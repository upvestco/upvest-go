package upvest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	// URLEncodeHeader is the content-type header for OuAth2
	URLEncodeHeader = "application/x-www-form-urlencoded"
	// clientele
	oauthPath = "/clientele/oauth2/token"
	grantType = "password"
	scope     = "read write echo transaction"
)

// Headers represent the HTTP headers sent to Upvest API
type Headers map[string]string

// AuthProvider interface for authentication mechanisms supported by Upvest API
type AuthProvider interface {
	// GetHeaders returns authorization headers (or other info) to be attached to requests.
	GetHeaders(method, path string, body interface{}, c *Client) (Headers, error)
}

// OAuthResponse represents succesful OAuth response
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// OAuth (The OAuth2 Key Authentication) is used to authenticate requests on behalf of a user
type OAuth struct {
	clientID     string
	clientSecret string
	username     string
	password     string
}

// KeyAuth (The API Key Authentication) is used to authenticate requests as a tenant.
type KeyAuth struct {
	apiKey        string
	apiSecret     string
	apiPassphrase string
}

// GetHeaders returns authorization headers for requests as a tenant.
func (auth KeyAuth) GetHeaders(method, path string, body interface{}, c *Client) (Headers, error) {
	path1, _ := joinURLs(APIVersion, path)
	versionedPath := path1.String()

	var headers Headers
	timestamp := fmt.Sprintf("%d", makeTimestamp())
	// Compose the message as a concatenation of all info we are sending along with the request
	message := timestamp + method + versionedPath

	if body != nil {
		buf, err := jsonEncode(body)
		if err != nil {
			return nil, err
		}
		body1, _ := ioutil.ReadAll(buf)
		message = message + string(body1)
	}

	// Generate signature, in order to prevent manipulation of payload in flight
	h := hmac.New(sha512.New, []byte(auth.apiSecret))
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	// Generate message headers
	headers = Headers{
		"Content-Type":         "application/json",
		"X-UP-API-Key":         auth.apiKey,
		"X-UP-API-Signature":   signature,
		"X-UP-API-Timestamp":   timestamp,
		"X-UP-API-Passphrase":  auth.apiPassphrase,
		"X-UP-API-Signed-Path": versionedPath,
	}

	return headers, nil
}

// GetHeaders returns authorization headers for requests as a clientele
func (oauth OAuth) GetHeaders(method, path string, body interface{}, c *Client) (Headers, error) {
	resp, err := oauth.preFlight(c)
	if err != nil {
		return nil, errors.Wrap(err, "OAuth2 preflight request failed")
	}
	// Retrieve and return OAuth token
	headers := Headers{
		"Authorization": fmt.Sprintf("Bearer %s", resp.AccessToken),
		"Content-Type":  "application/json",
	}
	return headers, nil
}

func (oauth OAuth) preFlight(c *Client) (*OAuthResponse, error) {
	data := url.Values{}
	data.Add("grant_type", grantType)
	data.Add("scope", scope)
	data.Add("client_id", oauth.clientID)
	data.Add("client_secret", oauth.clientSecret)
	data.Add("username", oauth.username)
	data.Add("password", oauth.password)

	payload := bytes.NewBufferString(data.Encode())

	p := &Params{}
	// TODO: refactor this to pass content type to Call/CallRaw
	p.AddHeader("Content-Type", URLEncodeHeader)
	p.AddHeader("Cache-Control", "no-cache")

	resp := &OAuthResponse{}
	// 	err = c.CallRaw(http.MethodPost, oauthPath, buf, resp, contentType)
	err := c.Call(http.MethodPost, oauthPath, payload, resp, p)
	return resp, err
}
