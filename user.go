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

// List returns list of all users.
// For more details see https://doc.upvest.co/reference#tenancy_user_list
func (s *UserService) List() (*UserList, error) {
	path := "/tenancy/users/"
	u, _ := url.Parse(path)
	//q := u.Query()
	//q.Set("page_size", fmt.Sprintf("%d", maxPageSize))
	//u.RawQuery = q.Encode()

	var results []User
	users := &UserList{}

	for {
		err := s.client.Call("GET", u.String(), nil, users, s.auth)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of users")
		}
		results = append(results, users.Values...)

		// append page_size param to the returned params
		u1, err := url.Parse(users.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(maxPageSize))
		u.RawQuery = q.Encode()
		if users.Meta.Next == "" {
			break
		}
	}

	return &UserList{Values: results}, nil
}

// ListN returns a specific number of users
// For more details see https://doc.upvest.co/reference#tenancy_user_list
func (s *UserService) ListN(count int) (*UserList, error) {
	path := "/tenancy/users/"
	u, _ := url.Parse(path)
	q := u.Query()
	q.Set("page_size", fmt.Sprintf("%d", maxPageSize))
	u.RawQuery = q.Encode()

	var results []User
	users := &UserList{}

	total := 0

	for total <= count {
		err := s.client.Call("GET", u.String(), nil, users, s.auth)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of users")
		}
		results = append(results, users.Values...)
		total += len(users.Values)

		// append page_size param to the returned params
		u1, err := url.Parse(users.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(maxPageSize))
		u.RawQuery = q.Encode()
		if users.Meta.Next == "" {
			break
		}
	}

	return &UserList{Values: results}, nil
}
