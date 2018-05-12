package eth

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/cleanunicorn/ethereum/helper"

	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3/types"
)

// Eth module
type Eth struct {
	provider provider.Provider
}

// NewEth returns an instance of the eth module
func NewEth(p provider.Provider) Eth {
	return Eth{
		provider: p,
	}
}

// ResponseEthGetTransactionCount is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
type ResponseEthGetTransactionCount struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      uint   `json:"id"`
}

// GetTransactionCount returns how many transactions the account has.
// Can be directly used as the account's nonce because the nonce is counted from 0.
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
func (c Eth) GetTransactionCount(account string, block string) (uint64, error) {
	reply, err := c.provider.Call("eth_getTransactionCount", []interface{}{account, block})
	if err != nil {
		return 0, err
	}

	var transactionCountReply ResponseEthGetTransactionCount
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

// ResponseEthGetBalance is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance
type ResponseEthGetBalance struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      uint   `json:"id"`
}

// GetBalance returns the balance of the account at the specified block.
//
// block can be one of
//
//    latest	// most recent account balance
//    0x1	// account's balance at block 1
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
func (c Eth) GetBalance(account string, block string) (*big.Int, error) {
	reply, err := c.provider.Call("eth_getBalance", []interface{}{account, block})
	if err != nil {
		return big.NewInt(0), err
	}

	var balanceReply ResponseEthGetBalance
	err = json.Unmarshal(reply, &balanceReply)
	if err != nil {
		return big.NewInt(0), err
	}

	balance := helper.HexStrToBigInt(balanceReply.Result)

	return balance, nil
}

// ResponseEthSendRawTransaction is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
type ResponseEthSendRawTransaction struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      uint   `json:"id"`
}

// ResponseEthSendRawTransactionError is the structure returned when an invalid signed transaction was sent
type ResponseEthSendRawTransactionError struct {
	Jsonrpc string `json:"jsonrpc"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

// SendRawTransaction send a signed transaction to the endpoint and returns the transaction hash.
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func (c Eth) SendRawTransaction(signedTransaction string) (string, error) {
	reply, err := c.provider.Call("eth_sendRawTransaction", []interface{}{signedTransaction})
	if err != nil {
		return "", err
	}

	var transactionHashReply ResponseEthSendRawTransaction
	err = json.Unmarshal(reply, &transactionHashReply)
	if err != nil {
		return "", err
	}

	if strings.Compare(transactionHashReply.Result, "") == 0 {
		var transactionError ResponseEthSendRawTransactionError
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

// ResponseEthBlockNumber is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber
type ResponseEthBlockNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      uint   `json:"id"`
}

// BlockNumber returns the latest block number.
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber
func (c Eth) BlockNumber() (*big.Int, error) {
	reply, err := c.provider.Call("eth_blockNumber", []interface{}{})
	if err != nil {
		return big.NewInt(0), err
	}

	var blockNumberReply ResponseEthBlockNumber
	err = json.Unmarshal(reply, &blockNumberReply)
	if err != nil {
		return big.NewInt(0), err
	}

	blockNumber := helper.HexStrToBigInt(blockNumberReply.Result)

	return blockNumber, nil
}

// ResponseEthGetBlockByNumber is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
type ResponseEthGetBlockByNumber struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  types.Block `json:"result"`
}

// GetBlockByNumber returns the block with or without the transaction list included
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func (c Eth) GetBlockByNumber(blockNumberHex string, includeTransactions bool) (types.Block, error) {
	reply, err := c.provider.Call("eth_getBlockByNumber", []interface{}{
		blockNumberHex, includeTransactions,
	})
	if err != nil {
		return types.Block{}, err
	}

	var getBlockReply ResponseEthGetBlockByNumber
	err = json.Unmarshal(reply, &getBlockReply)
	if err != nil {
		return types.Block{}, err
	}

	b := getBlockReply.Result
	b.RawTransactions = json.RawMessage(`{}`)

	if includeTransactions {
		err := json.Unmarshal(getBlockReply.Result.RawTransactions, &b.Transactions)
		if err != nil {
			return types.Block{}, err
		}
	} else {
		err := json.Unmarshal(getBlockReply.Result.RawTransactions, &b.TransactionHashes)
		if err != nil {
			return types.Block{}, err
		}
	}

	return b, nil
}

// ResponseEthGetTransactionReceipt is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
type ResponseEthGetTransactionReceipt struct {
	Jsonrpc string        `json:"jsonrpc"`
	Result  types.Receipt `json:"result"`
	ID      int           `json:"id"`
}

// GetTransactionReceipt returns a transaction receipt
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
func (c Eth) GetTransactionReceipt(transactionHash string) (types.Receipt, error) {
	reply, err := c.provider.Call("eth_getTransactionReceipt", []interface{}{
		transactionHash,
	})
	if err != nil {
		return types.Receipt{}, err
	}

	var responseReceipt ResponseEthGetTransactionReceipt
	err = json.Unmarshal(reply, &responseReceipt)
	if err != nil {
		return types.Receipt{}, err
	}

	return responseReceipt.Result, nil
}
