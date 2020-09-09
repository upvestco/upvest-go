package upvest

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

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

// Balance has a quantity and an asset
type Balance struct {
	Amount   int64  `json:"amount"`
	AssetID  string `json:"asset_id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Exponent int    `json:"exponent"`
}

// SignatureParams is the set of parameters that can be used when signing a wallet
// For more details see https://doc.upvest.co/reference#kms_sign
type SignatureParams struct {
	//Params   `json:"-"`
	Password     string `json:"password"`
	ToSign       string `json:"to_sign"`
	InputFormat  string `json:"input_format,omitempty"`
	OutputFormat string `json:"output_format,omitempty"`
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

// Signature represents the signed wallet signature
// For more details, see https://doc.upvest.co/reference#kms_sign
type Signature struct {
	// Has the same value as the "output_format" parameter.
	// The name of the string format for the big numbers in the signature. (Some JSON implementations can not handle integers which need more than 64 bits to be represented.)
	BigNumberFormat string `json:"big_number_format"`

	// The encryption algorithm used. (Currently only ECDSA)
	Algorithm string `json:"algorithm"`

	// The name of the elliptic curve used.
	Curve string `json:"curve"`

	// The x coordinate of the public key of the wallet
	PublicKey map[string]interface{} `json:"public_key"`

	// The "r" signature component.
	// Represented in the format given in the "big_number_format" field.
	R string `json:"r"`

	// The "s" signature component.
	// Represented in the format given in the "big_number_format" field.
	S string `json:"s"`

	// The "recover" signature component, sometimes also called "v".
	// Since this is a small integer with less than 64 bits, this is an actual JSON integer, and NOT represented in the big integer format.
	Recover string `json:"recover"`
}

// Sign signs (the hash of) data with the private key corresponding to this wallet.
// For more details, see https://doc.upvest.co/reference#kms_sign
func (s *WalletService) Sign(walletID string, sp *SignatureParams) (*Signature, error) {
	u := fmt.Sprintf("/kms/wallets/%s/sign", walletID)
	sig := &Signature{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, sp, sig, p)
	return sig, err
}

// Create creates a new wallet
// For more details https://doc.upvest.co/reference#kms_wallet_create
func (s *WalletService) Create(wp *WalletParams) (*Wallet, error) {
	u := "/kms/wallets/"
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
		if err != nil {
			return nil, errors.Wrap(err, "Can not parse url")
		}
		q := u1.Query()
		q.Set("page_size", strconv.Itoa(MaxPageSize))
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
		if err != nil {
			return nil, errors.Wrap(err, "Can not parse url")
		}
		q := u1.Query()
		q.Set("page_size", strconv.Itoa(MaxPageSize))
		u.RawQuery = q.Encode()
		if wallets.Meta.Next == "" {
			break
		}
	}

	return &WalletList{Values: results}, nil
}
