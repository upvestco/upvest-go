package upvest

import (
	"encoding/json"
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

	// defaultHTTPTimeout is the default timeout on the http client
	defaultHTTPTimeout = 60 * time.Second

	// base URL for all requests. default to playground environment
	defaultBaseURL = "https://api.playground.upvest.co/"

	// User agent used when communicating with the Upvest API.
	userAgent = "upvest-go/" + version

	apiVersion = "1.0"
	oauthPath  = "/clientele/oauth2/token"
	encoding   = "utf-8"
	grantType  = "password"
	scope      = "read write echo transaction"
)

type service struct {
	client *Client
}

// Client manages communication with the Upvest API
type Client struct {
	common service      // Reuse a single struct instead of allocating one for each service on the heap.
	client *http.Client // HTTP client used to communicate with the API.

	// the API Key used to authenticate all Upvest API requests
	key string

	baseURL *url.URL

	logger Logger
	// Services supported by the Upvest API.
	// Miscellaneous actions are directly implemented on the Client object
	User *UserService

	LoggingEnabled bool
	Log            Logger
}

// Logger interface for custom loggers
type Logger interface {
	Printf(format string, v ...interface{})
}

// Response represents arbitrary response data
type Response map[string]interface{}

// RequestValues aliased to url.Values as a workaround
type RequestValues url.Values

// MarshalJSON to handle custom JSON decoding for RequestValues
func (v RequestValues) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, 3)
	for k, val := range v {
		m[k] = val[0]
	}
	return json.Marshal(m)
}

// ListMeta is pagination metadata for paginated responses from the Upvest API
type ListMeta struct {
	Previous int `json:"previous"`
	Next     int `json:"next"`
}

// NewClient creates a new Upvest API client with the given base URL
// and HTTP client, allowing overriding of the HTTP client to use.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	u, _ := url.Parse(baseURL)

	c := &Client{
		client:         httpClient,
		baseURL:        u,
		LoggingEnabled: true,
		Log:            log.New(os.Stderr, "", log.LstdFlags),
	}

	c.common.client = c
	return c
}

// Get makes an HTTP GET request to Upvest API
func (c *Client) Get(method, path string, body, v interface{}, auth AuthProvider) error {
	return c.Call("GET", path, body, v, auth)
}

// Post makes an HTTP Post request to Upvest API
func (c *Client) Post(method, path string, body, v interface{}, auth AuthProvider) error {
	return c.Call("POST", path, body, v, auth)
}

// Patch makes an HTTP PATCH request to Upvest API
func (c *Client) Patch(method, path string, body, v interface{}, auth AuthProvider) error {
	return c.Call("Patch", path, body, v, auth)
}

// Delete makes an HTTP DELETE request to Upvest API
func (c *Client) Delete(method, path string, body, v interface{}, auth AuthProvider) error {
	return c.Call("Delete", path, body, v, auth)
}

// Call actually does the HTTP request to Upvest API
func (c *Client) Call(method, path string, body, v interface{}, auth AuthProvider) error {
	buf, err := jsonEncode(body)
	if err != nil {
		return errors.Wrap(err, "json encoding failed")
	}
	u, _ := c.baseURL.Parse(path)
	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		if c.LoggingEnabled {
			c.Log.Printf("Cannot create Upvest request: %v\n", err)
		}
		return errors.Wrap(err, "could not create HTTP request object")
	}

	// set headers
	req.Header.Set("User-Agent", userAgent)
	authHeaders, err := auth.GetHeaders(method, path, body)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range authHeaders {
		req.Header.Set(k, v)
	}

	if c.LoggingEnabled {
		c.Log.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
		c.Log.Printf("POST request data %v\n", buf)
	}

	start := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if c.LoggingEnabled {
		c.Log.Printf("Completed in %v\n", time.Since(start))
	}

	defer resp.Body.Close()
	return c.decodeResponse(resp, v)
}

// decodeResponse decodes the JSON response from the Upvest API.
// The actual response will be written to the `v` parameter
func (c *Client) decodeResponse(httpResp *http.Response, v interface{}) error {
	var resp Response
	respBody, err := ioutil.ReadAll(httpResp.Body)
	json.Unmarshal(respBody, &resp)

	if httpResp.StatusCode >= 300 {
		if c.LoggingEnabled {
			c.Log.Printf("Upvest error: %+v", err)
		}
		return newAPIError(httpResp)
	}

	if c.LoggingEnabled {
		c.Log.Printf("Upvest response: %v\n", resp)
	}

	return mapstruct(resp, v)
}
