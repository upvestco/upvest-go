package upvest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	// library version
	version = "0.1.0"

	// DefaultHTTPTimeout is the default timeout on the http client
	DefaultHTTPTimeout = 60 * time.Second

	// DefaultBaseURL for all requests. default to playground environment
	DefaultBaseURL = "https://api.playground.upvest.co/"

	// UserAgent used when communicating with the Upvest API.
	defaultUserAgent = "upvest-go/" + version

	// APIVersion is the currently supported API version
	APIVersion = "1.0"

	// Encoding is the text encoding to use
	Encoding = "utf-8"

	// MaxPageSize is the maximum page size when retrieving list
	MaxPageSize = 100
)

type service struct {
	client *Client
	auth   AuthProvider
}

// Client manages communication with the Upvest API
// Service specific actions are implemented on resource services mapped to the Upvest API.
// Miscellaneous actions are directly implemented on the Client object
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	baseURL   *url.URL
	useragent string

	LoggingEnabled bool
	Log            Logger
}

// Logger interface for custom loggers
type Logger interface {
	Printf(format string, v ...interface{})
}

// Response represents arbitrary response data
type Response map[string]interface{}

// ListMeta is pagination metadata for paginated responses from the Upvest API
type ListMeta struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

// NewClient creates a new Upvest API client with the given base URL
// and HTTP client, allowing overriding of the HTTP client to use.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultHTTPTimeout}
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	u, _ := url.Parse(baseURL)

	c := &Client{
		client:         httpClient,
		baseURL:        u,
		LoggingEnabled: false,
		Log:            log.New(os.Stderr, "", log.LstdFlags),
	}

	return c
}

func (c *Client) log(format string, a ...interface{}) {
	if c.LoggingEnabled {
		c.Log.Printf(format, a...)
	}
}

// SetUA sets the useragent to the provided value
func (c *Client) SetUA(userAgent string) {
	c.useragent = userAgent
}

// Call actually does the HTTP request to Upvest API
// TODO: refactor additional params into Param struct
func (c *Client) Call(method, path string, body, v interface{}, p *Params) error {
	req, err := c.NewRequest(method, path, body, p)
	if err != nil {
		return err
	}
	start := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	c.log("Completed in %v\n", time.Since(start))

	defer resp.Body.Close()
	return c.decodeResponse(resp, v)
}

// NewRequest is used by Call to generate an http.Request. It handles encoding
// parameters and attaching the appropriate headers.
func (c *Client) NewRequest(method, path string, body interface{}, params *Params) (*http.Request, error) {
	u, err := joinURLs(c.baseURL.String(), APIVersion, path)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request path")
	}

	var payload io.ReadWriter
	payload = new(bytes.Buffer)

	if params.Headers.Get("Content-Type") != URLEncodeHeader {
		if body == nil {
			body = map[string]string{}
		}
		payload, err = jsonEncode(body)
		if err != nil {
			return nil, errors.Wrap(err, "json encoding failed")
		}
	} else {
		payload = body.(io.ReadWriter)
	}

	req, err := http.NewRequest(method, u.String(), payload)
	c.log("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
	c.log("POST request data %v\n", payload)

	if err != nil {
		c.log("Cannot create Upvest request: %v\n", err)
		return nil, errors.Wrap(err, "could not create HTTP request object")
	}

	// set user agent
	if c.useragent != "" {
		req.Header.Set("User-Agent", c.useragent)
	} else {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	// add in extra headers
	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Set(k, v[0])
		}
	}

	// Get the headers from the auth provider
	if params.AuthProvider != nil {
		authHeaders, err := params.AuthProvider.GetHeaders(method, path, body, c)
		if err != nil {
			log.Println(err)
			return nil, errors.Wrap(err, "")
		}
		// Execute request with authenticated headers
		for k, v := range authHeaders {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

// decodeResponse decodes the JSON response from the Upvest API.
// The actual response will be written to the `v` parameter
func (c *Client) decodeResponse(httpResp *http.Response, v interface{}) error {
	if httpResp.StatusCode >= http.StatusBadRequest {
		err := NewError(httpResp)
		c.log("Upvest error: %+v", err)
		return err
	}

	var resp Response
	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(respBody, &resp)

	c.log("Upvest response: %v\n", resp)

	return mapstruct(resp, v)
}
