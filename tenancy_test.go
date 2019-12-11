package upvest

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// Tests API call to create, retrieve and delete user
func TestUserCRUD(t *testing.T) {
	uid, _ := uuid.NewUUID()
	username := fmt.Sprintf("upvest_test_%s", uid.String())
	user, err := tenancyTestClient.User.Create(username, randomString(12))
	if err != nil {
		t.Errorf("CREATE User returned error: %v", err)
	}
	if user.Username != username {
		t.Errorf("Expected User username %v, got %v", username, user.Username)
	}
	if user.RecoveryKit == "" {
		t.Errorf("Expected User recovery kit to be set, got nil")
	}

	// retrieve the user
	user1, err := tenancyTestClient.User.Get(user.Username)
	if err != nil {
		t.Errorf("GET User returned error: %v", err)
	}

	if user.Username != user1.Username {
		t.Errorf("Expected User username %v, got %v", user.Username, user1.Username)
	}

	// delete user
	_ = tenancyTestClient.User.Delete(user.Username)
	usr, err := tenancyTestClient.User.Get(user.Username)
	aerr := err.(*Error)

	if aerr.StatusCode != 404 {
		t.Errorf("Expected user not found, got %s", usr.Username)
	}
}

// Tests an API call to get list of users
func TestListUsers(t *testing.T) {
	expected := 10
	users, err := tenancyTestClient.User.List()
	if err != nil {
		t.Errorf("List Users returned error: %v", err)
	}
	if len(users.Values) < expected {
		t.Errorf("Expected greater than %d users, got %d", expected, len(users.Values))
	}
}

// Tests an API call to get list of specific number of users
func TestListNUsers(t *testing.T) {
	expected := 10
	users, err := tenancyTestClient.User.ListN(expected)
	if err != nil {
		t.Errorf("List Users returned error: %v", err)
	}

	if len(users.Values) != expected {
		t.Errorf("Expected greater than %d users, got %d", expected, len(users.Values))
	}
}

// Tests an API call to update a user's password
func TestChangePassword(t *testing.T) {
	user, pw := createTestUser()
	username := user.Username
	params := &ChangePasswordParams{
		OldPassword: pw,
		NewPassword: randomString(12),
	}
	user, _ = tenancyTestClient.User.ChangePassword(username, params)
	if user.Username != username {
		t.Errorf("Expected username %s, got %+v", username, user)
	}
}

// Test to retrive all assets
func TestListAssets(t *testing.T) {
	assets, err := tenancyTestClient.Asset.List()

	if err != nil {
		t.Errorf("List assets returned error: %v", err)
	}

	asset1 := assets.Values[0]
	assertions := []bool{
		asset1.ID == ethARWeaveAsset.ID,
		asset1.Name == ethARWeaveAsset.Name,
		asset1.Symbol == ethARWeaveAsset.Symbol,
		asset1.Exponent == ethARWeaveAsset.Exponent,
		asset1.Protocol == ethARWeaveAsset.Protocol,
	}

	for _, isValid := range assertions {
		if !isValid {
			t.Errorf("Asset structure does not match expected")
		}
	}
}

func TestWebhookVerify(t *testing.T) {
	verificationURL := os.Getenv("WEBHOOK_VERIFICATION_URL")
	isVerified := tenancyTestClient.Webhook.Verify(verificationURL)
	if !isVerified {
		t.Errorf("webhook verification failed: %s", verificationURL)
	}
}

// Tests API call to create, retrieve and delete webhook
func TestWebhookCRUD(t *testing.T) {
	url := os.Getenv("WEBHOOK_URL")
	uid, _ := uuid.NewUUID()
	webhook := &WebhookParams{
		URL:     url,
		Name:    fmt.Sprintf("test-webhook-%s", uid.String()),
		Headers: map[string]string{"X-Test": "Hello world!"},
		Version: "1.2",
		Status:  "ACTIVE",
		EventFilters: []EventFilterScope{
			"upvest.wallet.created",
			"ropsten.block.*",
			"upvest.echo.post",
		},
		HMACSecretKey: "abcdef",
	}

	wh, err := tenancyTestClient.Webhook.Create(webhook)
	if err != nil {
		t.Errorf("CREATE Webhook returned error: %v", err)
	}
	if wh.Name != webhook.Name {
		t.Errorf("Expected Webhook name %v, got %v", webhook.Name, wh.Name)
	}
	if wh.ID == "" {
		t.Errorf("Expected Webhook ID property to be set, got nil")
	}

	// retrieve the webhook
	wh1, err := tenancyTestClient.Webhook.Get(wh.ID)
	if err != nil {
		t.Errorf("GET Webhook returned error: %v", err)
	}

	if wh.Name != wh.Name {
		t.Errorf("Expected Webhook Name %v, got %v", wh.Name, wh1.Name)
	}

	// delete webhook
	_ = tenancyTestClient.Webhook.Delete(wh.ID)
	wh2, err := tenancyTestClient.Webhook.Get(wh.ID)
	aerr := err.(*Error)

	if aerr.StatusCode != 404 {
		t.Errorf("Expected Webhooknot found, got %s", wh2.Name)
	}
}

func TestListWebhooks(t *testing.T) {
	// list all webhooks
	webhooks, err := tenancyTestClient.Webhook.List()

	if err != nil {
		t.Errorf("List webhooks returned error: %v", err)
	}

	if reflect.ValueOf(webhooks).Kind() != reflect.Struct {
		t.Errorf("Retruned list is not webhook list: %s", err)
	}
}

const (
	EthProtocol       = "ethereum"
	EthRopstenNetwork = "ropsten"
)

func TestHistoricalAssetBalance(t *testing.T) {
	Address := "0x93b3d0b2894e99c2934bed8586ea4e2b94ce6bfd"
	balance, err := tenancyTestClient.Historical.GetAssetBalance(EthProtocol, EthRopstenNetwork, Address)
	if err != nil {
		t.Errorf("GET asset balance returned error: %v", err)
	}

	if balance.Address != Address {
		t.Errorf("Expected Asset Address %v, got %v", Address, balance.Address)
	}
}

func TestHistoricalContractBalance(t *testing.T) {
	ToAddr := "0x93b3d0b2894e99c2934bed8586ea4e2b94ce6bfd"
	ContractAddr := "0x1d7cf6ad190772cc6177beea2e3ae24cc89b2a10"
	balance, err := tenancyTestClient.Historical.GetContractBalance(EthProtocol, EthRopstenNetwork, ToAddr, ContractAddr)
	if err != nil {
		t.Errorf("GET contract balance returned error: %v", err)
	}

	if balance.Address != ToAddr {
		t.Errorf("Expected Contract Address %v, got %v", ToAddr, balance.Address)
	}
}

func TestHistoricalStatus(t *testing.T) {
	status, err := tenancyTestClient.Historical.GetStatus(EthProtocol, EthRopstenNetwork)
	if err != nil {
		t.Errorf("GET status returned error: %v", err)
	}

	if status.Lowest == "" {
		t.Errorf("Expected Status.Lowest to be set, got nil")
	}

	if status.Highest == "" {
		t.Errorf("Expected Status.Highest to be set, got nil")
	}

	if status.Latest == "" {
		t.Errorf("Expected Status.Latest to be set, got nil")
	}
}

func TestHistoricalTxByHash(t *testing.T) {
	TxHash := "0xa313aaad0b9b1fd356f7f42ccff1fa385a2f7c2585e0cf1e0fb6814d8bdb559a"
	tx, err := tenancyTestClient.Historical.GetTxByHash(EthProtocol, EthRopstenNetwork, TxHash)
	if err != nil {
		t.Errorf("GET transactions by txhash returned error: %v", err)
	}

	if tx.Hash != TxHash[1:len(TxHash)] {
		t.Errorf("Expected Transaction Hash %v, got %v", TxHash, tx.Hash)
	}
}

func TestHistoricalBlock(t *testing.T) {
	blockNumber := "6570890"
	block, err := tenancyTestClient.Historical.GetBlock(EthProtocol, EthRopstenNetwork, blockNumber)
	if err != nil {
		t.Errorf("GET block returned error: %v", err)
	}

	if block.Number != blockNumber {
		t.Errorf("Expected block number %v, got %v", blockNumber, block.Number)
	}
}

func TestHistoricalTransactions(t *testing.T) {
	address := "0x6590896988376a90326cb2f741cb4f8ace1882d5"
	confirmations := 1000
	opts := &TxFilters{Confirmations: confirmations}
	txns, err := tenancyTestClient.Historical.GetTransactions(EthProtocol, EthRopstenNetwork, address, opts)
	if err != nil {
		t.Errorf("GET transactions returned error: %v", err)
	}

	for _, v := range txns.Values {
		if !(v.Confirmations > confirmations) {
			t.Errorf("GET transactions confirmations mismatch. Expected %d, got %d", confirmations, v.Confirmations)
		}
	}
}
