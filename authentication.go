package upvest

// AuthProvider interface for authentication mechanisms supported by Upvest API
type AuthProvider interface {
	// GetHeaders returns authorization headers (or other info) to be attached to requests.
	GetHeaders() map[string]string
}

// KeyAuth (The API Key Authentication) is used to authenticate requests as a tenant.
type KeyAuth struct {
	apiKey        string
	apiSecret     string
	apiPassphrase string
}

// GetHeaders returns authorization headers for requests as a tenant.
func (auth KeyAuth) GetHeaders() {

}
