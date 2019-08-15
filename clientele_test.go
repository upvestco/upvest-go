package upvest

import (
	"encoding/hex"
	"testing"
)

// Tests an API call to create a wallet
func TestWalletCRUD(t *testing.T) {
	wp := &WalletParams{
		Password: staticUserPW,
		AssetID:  ethWallet.Balances[0].AssetID,
		// Type:     "encrypted",
		// Index:    0,
	}

	// create the wallet
	wallet, err := clienteleTestClient.Wallet.Create(wp)
	if err != nil {
		t.Errorf("CREATE Wallet returned error: %v", err)
	}

	// retrieve the wallet
	wallet1, err := clienteleTestClient.Wallet.Get(wallet.ID)
	if err != nil {
		t.Errorf("GET Wallet returned error: %v", err)
	}

	if wallet.Address != wallet1.Address {
		t.Errorf("Expected Wallet address %v, got %v", wallet.Address, wallet1.Address)
	}
}

func TestWalletList(t *testing.T) {
	_, err := clienteleTestClient.Wallet.List()
	if err != nil {
		t.Errorf("Wallet list returned error: %v", err)
	}
}

func TestTransactionCRUD(t *testing.T) {
	tp := &TransactionParams{
		Password:  staticUserPW,
		AssetID:   ethRopstenAssetID,
		Quantity:  10000000000000000,
		Fee:       41180000000000,
		Recipient: "0xf9b44Ba370CAfc6a7AF77D0BDB0d50106823D91b",
	}

	// create the transaction
	txn, err := clienteleTestClient.Transaction.Create(ethWallet.ID, tp)
	if err != nil {
		t.Errorf("CREATE Transaction returned error: %v", err)
	}

	// retrieve the transaction
	txn1, err := clienteleTestClient.Transaction.Get(ethWallet.ID, txn.ID)
	if err != nil {
		t.Errorf("GET Transaction returned error: %v", err)
	}

	if txn1.ID != txn.ID {
		t.Errorf("Expected Transaction ID %v, got %v", txn.ID, txn1.ID)
	}

	if txn1.AssetID != ethRopstenAssetID {
		t.Errorf("Expected Transaction asset ID %v, got %v", ethRopstenAssetID, txn1.AssetID)
	}

	if txn1.TxHash != txn.TxHash {
		t.Errorf("Expected Transaction hash %v, got %v", txn.TxHash, txn1.TxHash)
	}
}

func TestTransactionList(t *testing.T) {
	wallets, err := clienteleTestClient.Wallet.List()
	wallet1 := wallets.Values[0]
	_, err = clienteleTestClient.Transaction.List(wallet1.ID)
	if err != nil {
		t.Errorf("Transaction list returned error: %v", err)
	}
}

func TestWalletSign(t *testing.T) {
	sp := &SignatureParams{
		Password:     staticUserPW,
		ToSign:       hex.EncodeToString([]byte(randomString(32))),
		InputFormat:  "hex",
		OutputFormat: "hex",
	}
	_, err := clienteleTestClient.Wallet.Sign(ethWallet.ID, sp)
	if err != nil {
		t.Errorf("Wallet sign returned error: %v", err)
	}
}
