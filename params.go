package upvest

import (
	"net/http"
)

// Params is the structure that contains the common properties
// of any *Params structure.
type Params struct {
	// Headers may be used to provide extra header lines on the HTTP request.
	Headers http.Header `json:"-"`

	// AuthProvider for authenticating the request
	AuthProvider AuthProvider `json:"-"`
}

// SetAuthProvider sets a value for the auth mechanism
func (p *Params) SetAuthProvider(auth AuthProvider) {
	p.AuthProvider = auth
}

// AddHeader adds a new arbitrary key-value pair to the request header
func (p *Params) AddHeader(key, value string) {
	if p.Headers == nil {
		p.Headers = make(http.Header)
	}

	p.Headers.Add(key, value)
}
