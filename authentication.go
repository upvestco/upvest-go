package upvest

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
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
	ExpiresIn    string `json:"exxpires_in"`
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
	body := map[string]string{
		"grant_type":    grantType,
		"scope":         scope,
		"client_id":     oauth.clientID,
		"client_secret": oauth.clientSecret,
		"username":      oauth.username,
		"password":      oauth.password,
	}
	buf, err := jsonEncode(body)
	if err != nil {
		return nil, err
	}

	p := &Params{}
	// TODO: refactor this to pass content type to Call/CallRaw
	p.Headers.Add("Content-Type", "application/x-www-form-urlencoded")
	resp := &OAuthResponse{}
	// 	err = c.CallRaw(http.MethodPost, oauthPath, buf, resp, contentType)
	err = c.Call(http.MethodPost, oauthPath, buf, resp, p)
	return resp, err
}
