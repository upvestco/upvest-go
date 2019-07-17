package upvest

import "fmt"

// UserService handles operations related to the user
// For more details see https://doc.upvest.co/reference#tenancy_user_create
type UserService struct {
	service
}

// User is the resource representing your Upvest Tenant user.
// For more details see https://doc.upvest.co/reference#tenancy_user_create
type User struct {
	Username    string         `json:"username,omitempty"`
	RecoveryKit string         `json:"recoverykit,omitempty"`
	WalletIDs   map[int]string `json:"wallet_ids,omitempty"`
}

// UserList is a list object for users.
type UserList struct {
	Meta   ListMeta
	Values []User `json:"results"`
}

// Create creates a new user
// For more details https://doc.upvest.co/reference#tenancy_user_create
func (s *UserService) Create(username, password string) (*User, error) {
	u := fmt.Sprintf("/tenancy/users/")
	usr := &User{}
	data := map[string]string{"username": username, "password": password}
	err := s.Post(u, data, usr, s.auth)

	return usr, err
}
