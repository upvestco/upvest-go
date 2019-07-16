package upvest

// TenancyAPI represents Upvest tenancy API
// For more details, please see https://doc.upvest.co/reference#tenancy
type TenancyAPI struct {
	User *UserService
}

// NewTenant creates a new tenant for interacting with your Upvest tenant
func (c *Client) NewTenant(apiKey, apiSecret, apiPassphrase string) *TenancyAPI {
	auth := KeyAuth{apiKey: apiKey, apiSecret: apiSecret, apiPassphrase: apiPassphrase}
	svc := service{c, auth} // Reuse a single struct instead of allocating one for each service on the heap
	tenant := &TenancyAPI{
		User: &UserService{svc},
		// Asset: (*AssetService)(&svc),
	}
	return tenant
}
