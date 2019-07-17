package upvest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

var errorTypes = map[int]string{
	400: "invalid_request_error",
	401: "authentication_error",
	409: "duplicate_user",
	500: "server_error",
}

// APIError represents an error response from the Upvest API server
type APIError struct {
	Type       string                 `json:"type,omitempty"`
	Message    string                 `json:"message,omitempty"`
	StatusCode int                    `json:"code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	URL        *url.URL               `json:"url,omitempty"`
	Header     http.Header            `json:"header,omitempty"`
}

// APIError supports the error interface
func (aerr *APIError) Error() string {
	ret, _ := json.Marshal(aerr)
	return string(ret)
}

func newAPIError(resp *http.Response) *APIError {
	p, _ := ioutil.ReadAll(resp.Body)

	var upvestErrorResp map[string]interface{}
	_ = json.Unmarshal(p, &upvestErrorResp)

	errorType := "server_error"
	if err, ok := errorTypes[resp.StatusCode]; ok {
		errorType = err
	}

	return &APIError{
		Type:       errorType,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Details:    upvestErrorResp,
		URL:        resp.Request.URL,
	}
}
