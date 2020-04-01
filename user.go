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

// ChangePasswordParams is the set of parameters that can be used when changing user password
// For more details see https://doc.upvest.co/reference#tenancy_user_password_update
type ChangePasswordParams struct {
	//Params   `json:"-"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Create creates a new user
// For more details https://doc.upvest.co/reference#tenancy_user_create
func (s *UserService) Create(username, password string, assetIDs []string) (*User, error) {
	u := "/tenancy/users/"
	usr := &User{}
	data := map[string]interface{}{
		"username":  username,
		"password":  password,
		"asset_ids": assetIDs,
	}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, data, usr, p)
	return usr, err
}

// ChangePassword changes user password with the provided password
// For more details https://doc.upvest.co/reference#tenancy_user_password_update
func (s *UserService) ChangePassword(username string, params *ChangePasswordParams) (*User, error) {
	u := fmt.Sprintf("/tenancy/users/%s", username)
	usr := &User{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPatch, u, params, usr, p)
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
