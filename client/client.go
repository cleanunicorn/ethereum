package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Default parameters
const defaultHTTP = "http://127.0.0.1:8545"

// Client is the client pointing to the Ethereum node
type HTTPClient struct {
	client   *http.Client
	endpoint string
}

func DialHTTP(endpointHTTP string) (*HTTPClient, error) {
	var httpTransport = &http.Transport{}
	var httpClient = &http.Client{Transport: httpTransport}

	return NewHTTPClient(httpClient, endpointHTTP), nil
}

func NewHTTPClient(hc *http.Client, endpoint string) *HTTPClient {
	return &HTTPClient{
		client:   hc,
		endpoint: endpoint,
	}
}

func (c *HTTPClient) MakeCall(method string, args interface{}) ([]byte, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  args,
		"id":      1,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest("POST", c.endpoint, strings.NewReader(string(dataJSON)))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	request.Header.Add("Content-Type", "application/json")

	log.Debug("Making JSON request: ", string(dataJSON))
	response, err := c.client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// Check if the result is an error
	type responseError struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	var respErr responseError
	if err := json.Unmarshal(body, &respErr); err != nil {
		return []byte{}, err
	}
	if respErr.Error.Code != 0 {
		return []byte{}, fmt.Errorf("code: %d, error: %s", respErr.Error.Code, respErr.Error.Message)
	}

	log.Debug("Got reply body: ", string(body))
	return body, nil
}

func (c *HTTPClient) MakeRawCall(method string, args []interface{}) ([]byte, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  args,
		"id":      1,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest("POST", c.endpoint, strings.NewReader(string(dataJSON)))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	request.Header.Add("Content-Type", "application/json")

	log.Debug("Making JSON request: ", string(dataJSON))
	response, err := c.client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// Check if the result is an error
	type responseError struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	var respErr responseError
	if err := json.Unmarshal(body, &respErr); err != nil {
		return []byte{}, err
	}
	if respErr.Error.Code != 0 {
		return []byte{}, fmt.Errorf("code: %d, error: %s", respErr.Error.Code, respErr.Error.Message)
	}

	log.Debug("Got reply body: ", string(body))
	return body, nil
}

// Eth_getTransactionCount returns how many transactions the account has
// Can be directly used as the account's nonce because the nonce is counted from 0.
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
func (c *HTTPClient) Eth_getTransactionCount(account string, block string) (uint64, error) {
	reply, err := c.MakeCall("eth_getTransactionCount", []interface{}{account, block})
	if err != nil {
		return 0, err
	}

	type eth_getTransactionCount struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
		ID      uint   `json:"id"`
	}
	var transactionCountReply eth_getTransactionCount

	err = json.Unmarshal(reply, &transactionCountReply)
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseUint(transactionCountReply.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Eth_getBalance returns the balance of the account at the specified block.
//
// block can one of
//
//    latest	// most recent account balance
//    0x1	// account's balance at block 1
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
func (c *HTTPClient) Eth_getBalance(account string, block string) (*big.Int, error) {
	reply, err := c.MakeCall("eth_getBalance", []interface{}{account, block})
	if err != nil {
		return big.NewInt(0), err
	}

	type eth_getBalance struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
		ID      uint   `json:"id"`
	}
	var balanceReply eth_getBalance

	err = json.Unmarshal(reply, &balanceReply)
	if err != nil {
		return big.NewInt(0), err
	}

	balance := hexStrToBigInt(balanceReply.Result)

	return balance, nil
}

func (c *HTTPClient) Eth_sendRawTransaction(signedTransaction string) (string, error) {
	reply, err := c.MakeCall("eth_sendRawTransaction", []interface{}{signedTransaction})
	if err != nil {
		return "", err
	}

	type eth_sendRawTransaction struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
		ID      uint   `json:"id"`
	}
	var transactionHashReply eth_sendRawTransaction

	err = json.Unmarshal(reply, &transactionHashReply)
	if err != nil {
		return "", err
	}

	if strings.Compare(transactionHashReply.Result, "") == 0 {
		type eth_sendRawTransactionError struct {
			Jsonrpc string `json:"jsonrpc"`
			Error   struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			ID int `json:"id"`
		}

		var transactionError eth_sendRawTransactionError
		err := json.Unmarshal(reply, &transactionError)
		if err != nil {
			return "", err
		}

		if transactionError.Error.Code == 0 {
			return "", fmt.Errorf("unknown error, got reply: %s", string(reply))
		}

		return "", fmt.Errorf("no transaction hash generated, got error code: %d message: %s", transactionError.Error.Code, transactionError.Error.Message)
	}

	return transactionHashReply.Result, nil
}

// Net_version returns the current network id as described here
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
func (c *HTTPClient) Net_version() (int64, error) {
	reply, err := c.MakeCall("net_version", []interface{}{})
	if err != nil {
		return 0, err
	}

	var networkIDReply response_netVersion

	err = json.Unmarshal(reply, &networkIDReply)
	if err != nil {
		return 0, err
	}

	return networkIDReply.Result, nil
}

// Eth_blockNumber returns the latest block number.
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber
func (c *HTTPClient) Eth_blockNumber() (*big.Int, error) {
	reply, err := c.MakeCall("eth_blockNumber", []interface{}{})
	if err != nil {
		return big.NewInt(0), err
	}

	var blockNumberReply response_ethBlockNumber

	err = json.Unmarshal(reply, &blockNumberReply)
	if err != nil {
		return big.NewInt(0), err
	}

	blockNumber := hexStrToBigInt(blockNumberReply.Result)

	return blockNumber, nil
}

// Eth_getBlockByNumber returns the block with or without the transaction list included
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func (c *HTTPClient) Eth_getBlockByNumber(blockNumberHex string, includeTransactions bool) (Block, error) {
	reply, err := c.MakeCall("eth_getBlockByNumber", []interface{}{
		blockNumberHex, includeTransactions,
	})
	if err != nil {
		return Block{}, err
	}

	var responseBlock response_ethGetBlockByNumber
	err = json.Unmarshal(reply, &responseBlock)
	if err != nil {
		return Block{}, err
	}

	return responseBlock.Result, nil
}
