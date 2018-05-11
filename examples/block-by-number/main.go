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

	b, err := c.Eth.GetBlockByNumber("0x2", false)
	if err != nil {
		fmt.Printf("Error getting latest block, err: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Block: %#v\n", b)
	fmt.Printf("Difficulty: %#v\n", b.Difficulty.Int64())
	fmt.Printf("Number: %#v\n", b.Number.Int64())
}
