package upvest

import (
	"log"
	"os"
)

var tenancyTestClient *TenancyAPI

var staticUser *User

func init() {
	createTestTenancyClient()
	staticUser, _ = createTestUser()
}

// createTenancyClient creates an Upvest tenant client for testing purposes
func createTestTenancyClient() {
	c := NewClient("", nil)

	// use env var to enable debugging during development
	if _, ok := os.LookupEnv("DEBUG"); ok {
		log.Println("==> DEBUG mode")
		c.LoggingEnabled = true
	}

	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	apiPassphrase := os.Getenv("API_PASSPHRASE")
	tenancyTestClient = c.NewTenant(apiKey, apiSecret, apiPassphrase)
}
