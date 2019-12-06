package upvest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ErrorType is represents the allowed values for the error's type.
type ErrorType string

// List of values that ErrorType can take.
const (
	ErrInvalidRequest ErrorType = "invalid_request_error"
	ErrAuthorization            = "authorization_error"
	ErrAuthentication           = "authentication_error"
	ErrDuplicateUser            = "duplicate_user"
	ErrServer                   = "server_error"
)

var errorTypes = map[int]ErrorType{
	http.StatusBadRequest:          ErrInvalidRequest,
	http.StatusUnauthorized:        ErrAuthorization,
	http.StatusForbidden:           ErrAuthentication,
	http.StatusConflict:            ErrDuplicateUser,
	http.StatusInternalServerError: ErrServer,
}

// Error represents an error response from the Upvest API server
type Error struct {
	Type       ErrorType              `json:"type,omitempty"`
	Message    string                 `json:"message,omitempty"`
	StatusCode int                    `json:"code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	URL        *url.URL               `json:"url,omitempty"`
	Header     http.Header            `json:"header,omitempty"`
}

// httpError supports the error interface
func (aerr *Error) Error() string {
	ret, _ := json.Marshal(aerr)
	return string(ret)
}

// NewError parses http response and returns Upvest error type
func NewError(resp *http.Response) *Error {
	p, _ := ioutil.ReadAll(resp.Body)

	var upvestErrorResp map[string]interface{}
	_ = json.Unmarshal(p, &upvestErrorResp)
	var errorType ErrorType = "server_error"
	if err, ok := errorTypes[resp.StatusCode]; ok {
		errorType = err
	}

	// uperror := upvestErrorResp["error"].(map[string]interface{})

	return &Error{
		Type:       errorType,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Details:    upvestErrorResp,
		URL:        resp.Request.URL,
		//Message:    uperror["Message"].(string),
	}
}
