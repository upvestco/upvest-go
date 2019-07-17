package upvest

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

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

// func TestMain(m *testing.M) {
// 	createTestTenancyClient()
// 	code := m.Run()
// 	os.Exit(code)
// }
