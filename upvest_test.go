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
	staticUser          = &User{Username: "upvest_test_72eaa6e8-b59a-11e9-9737-8c8590323ff7"}
)

const (
	staticUserPW = "GMWJGSAPGATL"
)

func init() {
	initTestClients()
}

func createTestUser() (*User, string) {
	uid, _ := uuid.NewUUID()
	username := fmt.Sprintf("upvest_test_%s", uid.String())
	password := randomString(12)
	user, _ := tenancyTestClient.User.Create(username, password)
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
