package upvest

// createTenancyClient creates an Upvest tenant client for testing purposes
func createTenancyClient() *TenancyAPI {
	c := NewClient("", nil)
	return c.NewTenant("", "", "")
}
