package types

import "encoding/json"

// Block represents a block structure containing the full transaction list or the transaction hashes.
// It contains one of the two depending on the second bool parameter of eth_getBlockByNumber
type Block struct {
	Difficulty        ComplexNumber   `json:"difficulty"`
	ExtraData         string          `json:"extraData"`
	GasLimit          ComplexNumber   `json:"gasLimit"`
	GasUsed           ComplexNumber   `json:"gasUsed"`
	Hash              string          `json:"hash"`
	LogsBloom         string          `json:"logsBloom"`
	Miner             string          `json:"miner"`
	MixHash           string          `json:"mixHash"`
	Nonce             ComplexNumber   `json:"nonce"`
	Number            ComplexNumber   `json:"number"`
	ParentHash        string          `json:"parentHash"`
	ReceiptsRoot      string          `json:"receiptsRoot"`
	Sha3Uncles        string          `json:"sha3Uncles"`
	Size              ComplexNumber   `json:"size"`
	StateRoot         string          `json:"stateRoot"`
	Timestamp         string          `json:"timestamp"`
	TotalDifficulty   ComplexNumber   `json:"totalDifficulty"`
	TransactionsRoot  string          `json:"transactionsRoot"`
	Uncles            []string        `json:"uncles"`
	RawTransactions   json.RawMessage `json:"transactions"`
	Transactions      []Transaction   `json:"transactionData"`
	TransactionHashes []string        `json:"transactionHashes"`
}

// Transaction represents a transaction structure
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

// Receipt represents a transaction receipt
type Receipt struct {
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Logs              []struct {
		Address             string   `json:"address"`
		BlockHash           string   `json:"blockHash"`
		BlockNumber         string   `json:"blockNumber"`
		Data                string   `json:"data"`
		LogIndex            string   `json:"logIndex"`
		Topics              []string `json:"topics"`
		TransactionHash     string   `json:"transactionHash"`
		TransactionIndex    string   `json:"transactionIndex"`
		TransactionLogIndex string   `json:"transactionLogIndex"`
		Type                string   `json:"type"`
	} `json:"logs"`
	LogsBloom        string `json:"logsBloom"`
	Root             string `json:"root"`
	Status           string `json:"status"`
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
}
