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

// Tests an API call to get list of specific number of users
func TestListNUsers(t *testing.T) {
	expected := 10
	users, err := tenancyTestClient.User.ListN(expected)
	if err != nil {
		t.Errorf("List Users returned error: %v", err)
	}

	if len(users.Values) != expected {
		t.Errorf("Expected greater than %d users, got %d", expected, len(users.Values))
	}
}

// Tests an API call to update a user's password
func TestChangePassword(t *testing.T) {
	user, pw := createTestUser()
	newPassword := randomString(12)
	username := user.Username

	params := make(Params)
	params["old_password"] = pw
	params["new_password"] = newPassword

	user, _ = tenancyTestClient.User.Update(username, params)

	if user.Username != username {
		t.Errorf("Expected username %s, got %+v", username, user)
	}
}

// Tests an API call to update a user's password
func TestDeleteUser(t *testing.T) {
	user, _ := createTestUser()
	_ = tenancyTestClient.User.Delete(user.Username)
	usr, err := tenancyTestClient.User.Get(user.Username)
	aerr := err.(*APIError)

	if aerr.StatusCode != 404 {
		t.Errorf("Expected user not found, got %s", usr.Username)
	}
}

// Test to retrive all assets
func TestListAssets(t *testing.T) {
	assets, err := tenancyTestClient.Asset.List()

	if err != nil {
		t.Errorf("List assets returned error: %v", err)
	}

	asset1 := assets.Values[0]
	assertions := []bool{
		asset1.ID == "51bfa4b5-6499-5fe2-998b-5fb3c9403ac7",
		asset1.Name == "Arweave (internal testnet)",
		asset1.Symbol == "AR",
		asset1.Exponent == 12,
		asset1.Protocol == "arweave_testnet",
	}

	for _, isValid := range assertions {
		if !isValid {
			t.Errorf("Asset structure does not match expected")
		}
	}
}
