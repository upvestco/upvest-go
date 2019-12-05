package upvest

import (
	"fmt"
	"net/http"
)

// HDBlock represents block object from historical data API
type HDBlock struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             string   `json:"size"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Transactions     []string `json:"transactions"`
	Timestamp        string   `json:"timestamp"`
	Uncles           []string `json:"uncles"`
}

// HDTransaction represents transaction object from historical data API
type HDTransaction struct {
	blockHash        string `json:"blockHash"`
	blockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	transactionIndex string `json:"transactionIndex"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	Input            string `json:"input"`
	Confirmations    string `json:"confirmations"`
}

// HDBalance reprents balance of an asset or contract
// if native asset balance,contract is set to address of the contract
type HDBalance struct {
	ID               string `json:"id"`
	Address          string `json:"address"`
	Contract         string `json:"contract"`
	Balance          string `json:"balance"`
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionindex"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	Timestamp        string `json:"timestamp"`
	IsMainChain      bool   `json:"isMainChain"`
}

// HDTransactionList is a list of HDTransaction objects
type HDTransactionList struct {
	Transactions []HDTransaction `json:"result"`
	NextCursor   string          `json:"next_cursor"`
}

// HDStatus represents historical data API status object
type HDStatus struct {
	Lowest  string `json:"lowest"`
	Highest string `json:"highest"`
	Latest  string "latest"
}

// HistoricalDataService handles operations related to the historical data
type HistoricalDataService struct {
	service
}

// GetTxByHash transaction (single) by txhash
func (s *HistoricalDataService) GetTxByHash(protocol, network, txhash string) (*HDTransaction, error) {
	u := fmt.Sprintf("/data/%s/%s/transaction/%s", protocol, network, txhash)
	p := NewParams(s.auth)
	txn := &HDTransaction{}
	r := &hdresult{}
	err := s.client.Call(http.MethodGet, u, nil, r, p)
	if err == nil {
		err = mapstruct(r.Result, txn)
	}
	return txn, err
}

// TxFilters is for filtering historical Data API queries
type TxFilters struct {
	Before        string `json:"before,omitempty"`
	After         string `json:"after,omitempty"`
	Confirmations string `json:"confirmations,omitempty"`
}

// GetTransactions returns transactions that have been sent to and received by an address
func (s *HistoricalDataService) GetTransactions(protocol, network, address string, txfilter TxFilters) (*[]HDTransaction, error) {
	u := fmt.Sprintf("/data/%s/%s/transactions/%s", protocol, network, address)
	p := NewParams(s.auth)
	txns := &[]HDTransaction{}
	r := &hdresult{}
	err := s.client.Call(http.MethodGet, u, txfilter, r, p)
	if err == nil {
		err = mapstruct(r.Result, txns)
	}
	return txns, err
}

// GetBlock returns block details by blockNumber
func (s *HistoricalDataService) GetBlock(protocol, network, blockNumber string) (*HDBlock, error) {
	u := fmt.Sprintf("/data/%s/%s/block/%s", protocol, network, blockNumber)
	p := NewParams(s.auth)
	block := &HDBlock{}
	r := &hdresult{}
	err := s.client.Call(http.MethodGet, u, nil, r, p)
	if err == nil {
		err = mapstruct(r.Result, block)
	}
	return block, err
}

// GetAssetBalance returns native asset balance by address
func (s *HistoricalDataService) GetAssetBalance(protocol, network, address string) (*HDBalance, error) {
	u := fmt.Sprintf("/data/%s/%s/balance/%s", protocol, network, address)
	p := NewParams(s.auth)
	hdbalance := &HDBalance{}
	r := &hdresult{}
	err := s.client.Call(http.MethodGet, u, nil, r, p)
	if err == nil {
		err = mapstruct(r.Result, hdbalance)

	}
	return hdbalance, err
}

// GetContractBalance returns contract balance by address
func (s *HistoricalDataService) GetContractBalance(protocol, network, address, contractAddr string) (*HDBalance, error) {
	u := fmt.Sprintf("/data/%s/%s/balance/%s/%s", protocol, network, address, contractAddr)
	p := NewParams(s.auth)
	hdbalance := &HDBalance{}
	r := &hdresult{}
	err := s.client.Call(http.MethodGet, u, nil, r, p)
	if err == nil {
		err = mapstruct(r.Result, hdbalance)
	}
	return hdbalance, err
}

// GetStatus return Historical Data API status
func (s *HistoricalDataService) GetStatus(protocol, network string) (*HDStatus, error) {
	u := fmt.Sprintf("/data/%s/%s/status", protocol, network)
	p := NewParams(s.auth)
	hdstatus := &HDStatus{}
	r := &hdresult{}
	err := s.client.Call(http.MethodGet, u, nil, r, p)
	if err == nil {
		err = mapstruct(r.Result, hdstatus)
	}
	return hdstatus, err
}

type hdresult struct {

	// get result and map to struct
	Result map[string]interface{} `json:"result"`
}
