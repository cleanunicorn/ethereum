package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3"
)

func main() {
	var (
		endpointHTTP = flag.String("http", "https://mainnet.infura.io:8545", "HTTP endpoint")
	)
	flag.Parse()

	c := web3.NewClient(provider.DialHTTP(*endpointHTTP))

	res, err := c.Provider.Call("eth_blockNumber", nil)
	if err != nil {
		fmt.Printf("Error getting latest block, err: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(res))
}
