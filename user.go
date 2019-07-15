package upvest

// UserService handles operations related to the user
// For more details see https://doc.upvest.co/reference#tenancy_user_create
type UserService service

// User is the resource representing your Upvest Tenant user.
// For more details see https://doc.upvest.co/reference#tenancy_user_create
type User struct {
	Username   string         `json:"username,omitempty"`
	ReoveryKit string         `json:"reoverykit,omitempty"`
	WalletIDs  map[int]string `json:"wallet_ids,omitempty"`
}

// UserList is a list object for users.
type UserList struct {
	Meta   ListMeta
	Values []User `json:"results"`
}
