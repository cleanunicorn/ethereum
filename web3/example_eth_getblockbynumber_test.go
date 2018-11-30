package web3_test

import (
	"encoding/json"
	"fmt"

	"github.com/cleanunicorn/ethereum/web3/eth"

	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3"
)

func Example() {
	c := web3.NewClient(provider.DialHTTP("https://mainnet.infura.io"))

	// Get block number 16 without transaction data included
	b, err := c.Eth.GetBlockByNumber("0x10", false)
	if err != nil {
		fmt.Printf("Error getting latest block, err: %v", err)
	}

	// Print the block number
	fmt.Printf("Number: %#v\n", b.Number.Int64())

	// Output:
	// Number: 16
}

func Example_rawCall() {
	c := web3.NewClient(provider.DialHTTP("https://mainnet.infura.io"))

	res, err := c.Provider.Call(
		// Specify the call
		"eth_getBlockByNumber",
		// Encode the parameters as expected by the node
		[]interface{}{"0x10", true},
	)
	if err != nil {
		fmt.Printf("Error getting latest block, err: %v, res: %v", err, res)
	}

	// Decode the response
	var br eth.ResponseEthGetBlockByNumber
	err = json.Unmarshal(res, &br)
	if err != nil {
		fmt.Printf("Error unmarshalling response, err: %v", err)
	}

	// Print the block number
	fmt.Printf("Number: %#v\n", br.Result.Number.Int64())

	// Output:
	// Number: 16
}
