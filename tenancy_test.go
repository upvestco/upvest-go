package upvest

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
)

var tenancyTestClient *TenancyAPI

// createTenancyClient creates an Upvest tenant client for testing purposes
func createTestTenancyClient() {
	c := NewClient("", nil)

	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	apiPassphrase := os.Getenv("API_PASSPHRASE")
	log.Printf("\n\n===> key: %s secret: %s passphrase: %s \n\n", apiKey, apiSecret, apiPassphrase)
	tenancyTestClient = c.NewTenant(apiKey, apiSecret, apiPassphrase)
}

// Tests an API call to create a user"""
func TestRegisterUser(t *testing.T) {
	uid, _ := uuid.NewUUID()
	username := fmt.Sprintf("upvest_test_%s", uid.String())
	user, err := tenancyTestClient.User.Create(username, randomString(12))
	if err != nil {
		t.Errorf("CREATE User returned error: %v", err)
	}
	if user.Username != username {
		t.Errorf("Expected User username %v, got %v", username, user.Username)
	}
	if user.RecoveryKit == "" {
		t.Errorf("Expected User recovery kit to be set, got nil")
	}
}

func TestMain(m *testing.M) {
	createTestTenancyClient()
	code := m.Run()
	os.Exit(code)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}
