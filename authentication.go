package upvest

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

func (auth KeyAuth) GetHeaders() {

}
