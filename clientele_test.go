package upvest

import (
	"testing"

	"github.com/google/uuid"
)

// Tests an API call to create a wallet
func TestWalletCRUD(t *testing.T) {
	uid, _ := uuid.NewUUID()
	assets, err := tenancyTestClient.Asset.List()
	asset1 := assets.Values[0]
	wp := &WalletParams{
		Password: uid.String(),
		AssetID:  asset1.ID,
		Type:     "encrypted",
		Index:    0,
	}
	_, err = clienteleTestClient.Wallet.Create(wp)
	// create the wallet
	if err != nil {
		t.Errorf("CREATE Wallet returned error: %v", err)
	}
	// retrieve the wallet
	// retrieve the wallet list
}

func TestWalletList(t *testing.T) {
	_, err := clienteleTestClient.Wallet.List()
	if err != nil {
		t.Errorf("Wallet list returned error: %v", err)
	}
}
