package upvest

// https://doc.upvest.co/reference#tenancy
type TenancyAPI struct {
	User   *UserService
	client *Client
	auth   KeyAuth
}

// NewTenant creates a new tenant for interacting with your Upvest tenant
func (c *Client) NewTenant(apiKey, apiSecret, apiPassphrase string) *TenancyAPI {
	auth := KeyAuth{apiKey: apiKey, apiSecret: apiSecret, apiPassphrase: apiPassphrase}
	tenant := &TenancyAPI{
		User: (*UserService)(&c.common),
		//Asset:  (*AssetService)(&c.common),
		client: c,
		auth:   auth,
	}
	return tenant
}
