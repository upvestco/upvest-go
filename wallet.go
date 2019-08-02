package upvest

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
