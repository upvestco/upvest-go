package upvest

// ClienteleAPI represents Upvest Clientele API
// For more details, please see https://doc.upvest.co/reference#clientele
type ClienteleAPI struct {
	Wallet *WalletService
	//Asset  *AssetService
}

// NewClintele creates a new clientele for interacting with your Upvest clients/users
func (c *Client) NewClintele(clientID, clientSecret, username, password string) *ClienteleAPI {
	auth := OAuth{clientID: clientID, clientSecret: clientSecret, username: username, password: password}
	svc := service{c, auth} // Reuse a single struct instead of allocating one for each service
	clientele := &ClienteleAPI{
		Wallet: &WalletService{svc},
	}
	return clientele
}
