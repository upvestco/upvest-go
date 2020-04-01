package upvest

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

var (
	tenancyTestClient   *TenancyAPI
	clienteleTestClient *ClienteleAPI
	staticUser          = &User{
		Username: "upvest_test_72eaa6e8-b59a-11e9-9737-8c8590323ff7",
	}
	ethARWeaveAsset = &Asset{
		ID:       "51bfa4b5-6499-5fe2-998b-5fb3c9403ac7",
		Name:     "Arweave (internal testnet)",
		Symbol:   "AR",
		Exponent: 12,
		Protocol: "arweave_testnet",
	}

	ethWallet = &Wallet{
		ID: "8fc19cd0-8f50-4626-becb-c9e284d2315b",
		Balances: []Balance{
			Balance{
				Amount:   0,
				AssetID:  "cfc59efb-3b21-5340-ae96-8cadb4ce31a8",
				Name:     "Example coin",
				Symbol:   "COIN",
				Exponent: 12,
			}},
		Protocol: "ethereum_ropsten",
		Address:  "0xc4a284e55ab2f1c2feb23a0bfc56fca31b0c94a3",
		Status:   "ACTIVE",
		Index:    0,
	}
)

const (
	staticUserPW      = "GMWJGSAPGATL"
	ethRopstenAssetID = "deaaa6bf-d944-57fa-8ec4-2dd45d1f5d3f"
)

func init() {
	initTestClients()
}

func createTestUser() (*User, string) {
	uid, _ := uuid.NewUUID()
	username := fmt.Sprintf("upvest_test_%s", uid.String())
	password := randomString(12)
	user, _ := tenancyTestClient.User.Create(username, password, []string{})
	return user, password
}

// initTestClients creates an Upvest tenant client for testing purposes
func initTestClients() {
	c := NewClient("", nil)

	// use env var to enable debugging during development
	if _, ok := os.LookupEnv("DEBUG"); ok {
		log.Println("==> DEBUG mode")
		c.LoggingEnabled = true
	}

	// tenant API
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	apiPassphrase := os.Getenv("API_PASSPHRASE")
	tenancyTestClient = c.NewTenant(apiKey, apiSecret, apiPassphrase)

	// clientele API
	clientID := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	clienteleTestClient = c.NewClientele(clientID, clientSecret, staticUser.Username, staticUserPW)
}
