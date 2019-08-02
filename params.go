package upvest

import (
	"net/http"
)

// Params is the structure that contains the common properties
// of any *Params structure.
type Params struct {
	// Headers may be used to provide extra header lines on the HTTP request.
	Headers http.Header

	// AuthProvider for authenticating the request
	AuthProvider AuthProvider
}

// SetAuthProvider sets a value for the auth mechanism
func (p *Params) SetAuthProvider(auth AuthProvider) {
	p.AuthProvider = auth
}
