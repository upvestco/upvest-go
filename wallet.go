package upvest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Balance has a quantity and an asset
type Balance struct {
	Amount   string `json:"amount"`
	AssetID  string `json:"asset_id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Exponent string `json:"exponent"`
}

// Wallet represents an Upvest wallet
type Wallet struct {
	ID       string    `json:"id"`
	Path     string    `json:"path"`
	Balances []Balance `json:"balances"`
	Protocol string    `json:"protocol"`
	Address  string    `json:"address"`
	Status   string    `json:"status"`
	Index    int64     `json:"index"`
}

// WalletService handles operations related to the wallet
// For more details see https://doc.upvest.co/reference#kms_wallet_create
type WalletService struct {
	service
}

// WalletList is a list object for wallets.
type WalletList struct {
	Meta   ListMeta
	Values []Wallet `json:"results"`
}

// WalletParams is the set of parameters that can be used when creating or updating a wallet
// For more details see https://doc.upvest.co/reference#kms_wallet_create
type WalletParams struct {
	//Params   `json:"-"`
	Password string `json:"password"`
	AssetID  string `json:"asset_id"`
	Type     string `json:"type,omitempty"`
	Index    int    `json:"index,omitempty"`
}

// Create creates a new wallet
// For more details https://doc.upvest.co/reference#kms_wallet_create
func (s *WalletService) Create(wp *WalletParams) (*Wallet, error) {
	u := "kms/wallets/"
	wallet := &Wallet{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, wp, wallet, p)
	return wallet, err
}

// Get returns the details of a wallet.
// For more details see https://doc.upvest.co/reference#kms_wallets_read
func (s *WalletService) Get(walletID string) (*Wallet, error) {
	u := fmt.Sprintf("/kms/wallets/%s", walletID)
	wallet := &Wallet{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodGet, u, nil, wallet, p)
	return wallet, err
}

// List returns list of all wallets.
// For more details see https://doc.upvest.co/reference#wallet
func (s *WalletService) List() (*WalletList, error) {
	path := "/kms/wallets/"
	u, _ := url.Parse(path)
	p := &Params{}
	p.SetAuthProvider(s.auth)

	var results []Wallet
	wallets := &WalletList{}

	for {
		err := s.client.Call(http.MethodGet, u.String(), nil, wallets, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of wallets")
		}
		results = append(results, wallets.Values...)

		// append page_size param to the returned params
		u1, err := url.Parse(wallets.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if wallets.Meta.Next == "" {
			break
		}
	}

	return &WalletList{Values: results}, nil
}

// ListN returns a specific number of wallets
// For more details see https://doc.upvest.co/reference#tenancy_wallet_list
func (s *WalletService) ListN(count int) (*WalletList, error) {
	path := "/kms/wallets/"
	u, _ := url.Parse(path)
	// q := u.Query()
	// q.Set("page_size", fmt.Sprintf("%d", maxPageSize))
	// u.RawQuery = q.Encode()

	p := &Params{}
	p.SetAuthProvider(s.auth)
	var results []Wallet
	wallets := &WalletList{}

	total := 0

	for total <= count {
		err := s.client.Call(http.MethodGet, u.String(), nil, wallets, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of wallets")
		}
		results = append(results, wallets.Values...)
		total += len(wallets.Values)

		// append page_size param to the returned params
		u1, err := url.Parse(wallets.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if wallets.Meta.Next == "" {
			break
		}
	}

	return &WalletList{Values: results}, nil
}
