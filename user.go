package upvest

import (
	"fmt"
	"net/http"
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
	u := "/tenancy/users/"
	usr := &User{}
	data := map[string]string{"username": username, "password": password}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, data, usr, p)

	return usr, err
}

// Update changes user password with the provided password
// For more details https://doc.upvest.co/reference#tenancy_user_create
func (s *UserService) Update(username string, rp RequestParams) (*User, error) {
	u := fmt.Sprintf("/tenancy/users/%s", username)
	usr := &User{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPatch, u, rp, usr, p)

	return usr, err
}

// Delete permanently deletes a user
// For more details https://doc.upvest.co/reference#tenancy_user_create
func (s *UserService) Delete(username string) error {
	u := fmt.Sprintf("/tenancy/users/%s", username)
	resp := &Response{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodDelete, u, map[string]string{}, resp, p)

	return err
}

// Get returns the details of a user.
// For more details see
func (s *UserService) Get(username string) (*User, error) {
	u := fmt.Sprintf("/tenancy/users/%s", username)
	user := &User{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodGet, u, nil, user, p)
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

	p := &Params{}
	p.SetAuthProvider(s.auth)
	var results []User
	users := &UserList{}

	for {
		err := s.client.Call(http.MethodGet, u.String(), nil, users, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of users")
		}
		results = append(results, users.Values...)

		// append page_size param to the returned params
		u1, err := url.Parse(users.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
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
	// q := u.Query()
	// q.Set("page_size", fmt.Sprintf("%d", maxPageSize))
	// u.RawQuery = q.Encode()

	p := &Params{}
	p.SetAuthProvider(s.auth)
	var results []User
	users := &UserList{}

	total := 0

	for total <= count {
		err := s.client.Call(http.MethodGet, u.String(), nil, users, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of users")
		}
		results = append(results, users.Values...)
		total += len(users.Values)

		// append page_size param to the returned params
		u1, err := url.Parse(users.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if users.Meta.Next == "" {
			break
		}
	}

	return &UserList{Values: results}, nil
}
