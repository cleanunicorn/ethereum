package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3"
)

func main() {
	var (
		endpointHTTP = flag.String("http", "https://mainnet.infura.io:8545", "HTTP endpoint")
		method       = flag.String("method", "", "method")
	)
	flag.Parse()
	args := flag.Args()

	c := web3.NewClient(provider.DialHTTP(*endpointHTTP))

	var params []interface{}
	switch *method {
	case "eth_getBlockByNumber":
		blockNumber, err := strconv.ParseInt(args[0], 0, 64)
		if err != nil {
			fmt.Printf("Error transforming %s to number, err: %s", args[0], err)
			os.Exit(1)
		}

		includeTransactionData, err := strconv.ParseBool(args[1])
		if err != nil {
			fmt.Printf("Error transforming %s to bool, err: %s", args[1], err)
			os.Exit(1)
		}

		params = []interface{}{
			fmt.Sprintf("0x%x", blockNumber),
			includeTransactionData,
		}
	case "eth_getUncleCountByBlockNumber":
		blockNumber, err := strconv.ParseInt(args[0], 0, 64)
		if err != nil {
			fmt.Printf("Error transforming %s to number, err: %s", args[0], err)
			os.Exit(1)
		}

		params = []interface{}{
			fmt.Sprintf("0x%x", blockNumber),
		}
	default:
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	response, err := c.Provider.Call(*method, params)
	if err != nil {
		fmt.Println("Error making request, err: ", err)
	}

	// pretty print
	buf := new(bytes.Buffer)
	json.Indent(buf, response, "", "  ")
	fmt.Println(buf)
}
