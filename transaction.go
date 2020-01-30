package upvest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Transaction represents a wallet transaction
// For more details, see https://doc.upvest.co/reference#kms_transaction_create
type Transaction struct {
	ID        string `json:"id"`
	TxHash    string `json:"txhash"`
	WalletID  string `json:"wallet_id"`
	AssetID   string `json:"asset_id"`
	AssetName string `json:"asset_name"`
	Exponent  string `json:"exponent"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Quantity  string `json:"quantity"`
	Fee       string `json:"fee"`
	Status    string `json:"status"`
}

// TransactionParams is the set of parameters that can be used when creating a transaction
// For more details see https://doc.upvest.co/reference#kms_transaction_create
type TransactionParams struct {
	Password  string `json:"password"`
	AssetID   string `json:"asset_id"`
	Quantity  int64  `json:"quantity"`
	Fee       int64  `json:"fee"`
	Recipient string `json:"recipient"`
}

// TransactionService handles operations related to the transaction
// For more details see https://doc.upvest.co/reference#kms_transaction_create
type TransactionService struct {
	service
}

// TransactionList is a list object for transactions
type TransactionList struct {
	Meta   ListMeta
	Values []Transaction `json:"results"`
}

// Create creates a new transaction
// For more details https://doc.upvest.co/reference#kms_transaction_create
func (s *TransactionService) Create(walletID string, tp *TransactionParams) (*Transaction, error) {
	u := fmt.Sprintf("/kms/wallets/%s/transactions/", walletID)
	transaction := &Transaction{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, tp, transaction, p)
	return transaction, err
}

// Get returns the details of a transaction.
// For more details see https://doc.upvest.co/reference#kms_transactions_read
func (s *TransactionService) Get(walletID, txnID string) (*Transaction, error) {
	u := fmt.Sprintf("/kms/wallets/%s/transactions/%s", walletID, txnID)
	txn := &Transaction{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodGet, u, nil, txn, p)
	return txn, err
}

// List returns list of all transactions.
// For more details see https://doc.upvest.co/reference#kms_transaction_list
func (s *TransactionService) List(walletID string) (*TransactionList, error) {
	path := fmt.Sprintf("/kms/wallets/%s/transactions/", walletID)
	u, _ := url.Parse(path)
	p := &Params{}
	p.SetAuthProvider(s.auth)

	var results []Transaction
	transactions := &TransactionList{}

	for {
		err := s.client.Call(http.MethodGet, u.String(), nil, transactions, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of transactions")
		}
		results = append(results, transactions.Values...)

		// append page_size param to the returned params
		u1, err := url.Parse(transactions.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if transactions.Meta.Next == "" {
			break
		}
	}

	return &TransactionList{Values: results}, nil
}

// CreateComplex creates a complex transaction
// For more details https://doc.upvest.co/docs/complex-transactions
func (s *TransactionService) CreateComplex(walletID string, password string, tx string, fund bool) (*Transaction, error) {
	u := fmt.Sprintf("/kms/wallets/%s/transactions/complex", walletID)
	txn := &Transaction{}
	data := map[string]interface{}{"password": password, "tx": tx, "fund": fund}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, data, txn, p)
	return txn, err
}

// CreateRaw creates a raw transaction
// For more details https://doc.upvest.co/docs/complex-transactions
func (s *TransactionService) CreateRaw(walletID string, password string,
	rawTx string, fund bool, inputFormat string) (*Transaction, error) {
	u := fmt.Sprintf("/kms/wallets/%s/transactions/raw", walletID)
	txn := &Transaction{}
	data := map[string]interface{}{
		"password": password,
		"raw_tx": rawTx,
		"fund": fund,
		"input_format": inputFormat,
	}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodPost, u, data, txn, p)
	return txn, err
}
