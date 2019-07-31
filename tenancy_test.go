package upvest

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func createTestUser() (*User, string) {
	uid, _ := uuid.NewUUID()
	username := fmt.Sprintf("upvest_test_%s", uid.String())
	password := randomString(12)
	user, _ := tenancyTestClient.User.Create(username, password)
	return user, password
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

// Tests an API call to get a specific user
func TestGetUser(t *testing.T) {
	user, _ := createTestUser()
	user1, err := tenancyTestClient.User.Get(user.Username)
	if err != nil {
		t.Errorf("GET User returned error: %v", err)
	}

	if user.Username != user1.Username {
		t.Errorf("Expected User username %v, got %v", user.Username, user1.Username)
	}
}

// Tests an API call to get list of users
func TestListUsers(t *testing.T) {
	expected := 10
	users, err := tenancyTestClient.User.List()
	if err != nil {
		t.Errorf("List Users returned error: %v", err)
	}
	if len(users.Values) < expected {
		t.Errorf("Expected greater than %d users, got %d", expected, len(users.Values))
	}
}

// func TestMain(m *testing.M) {
// 	createStaticUser()
// 	code := m.Run()
// 	os.Exit(code)
// }
