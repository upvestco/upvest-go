package upvest

import (
	"fmt"
	"log"
	"net/url"

	"github.com/pkg/errors"
)

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
	err := s.client.Call("POST", u, data, usr, s.auth)

	return usr, err
}

// Get returns the details of a user.
// For more details see
func (s *UserService) Get(username string) (*User, error) {
	u := fmt.Sprintf("/tenancy/users/%s", username)
	user := &User{}
	log.Printf("==> %s = %+v", u, s)
	err := s.client.Call("GET", u, nil, user, s.auth)
	return user, err
}

// List returns a list of users.
// For more details see
func (s *UserService) List() (*UserList, error) {
	//u := fmt.Sprintf("/tenancy/users/?page_size=%d", maxPageSize)
	u := "/tenancy/users/"
	var allUsers []User
	users := &UserList{}
	for {
		err := s.client.Call("GET", u, nil, users, s.auth)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of users")
		}
		allUsers = append(allUsers, users.Values...)

		// build new url with params from the Next url
		url, err := url.Parse(users.Meta.Next)
		q := url.Query()
		u = fmt.Sprintf("/tenancy/users/?page_size=%d&%s", maxPageSize, q.Encode())
		if users.Meta.Next == "" {
			break
		}
	}
	return &UserList{Values: allUsers}, nil
}

// ListN returns a list of users
// For more details see
// func (s *UserService) ListN(count int) (*UserList, error) {
// 	// go until there's count size is met
// 	// stop if no next page
// 	return users, err
// }
