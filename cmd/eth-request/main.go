package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/cleanunicorn/ethereum/client"
)

func main() {
	var (
		endpointHTTP = flag.String("http", "http://127.0.0.1:8545", "HTTP endpoint")
		method       = flag.String("method", "", "method")
	)
	flag.Parse()
	args := flag.Args()

	c, err := client.DialHTTP(*endpointHTTP)
	if err != nil {
		fmt.Printf("Could not dial into HTTP endpoint: %s, err: %s", *endpointHTTP, err)
		os.Exit(1)
	}

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

	response, err := c.Call(*method, params)
	if err != nil {
		fmt.Println("Error making request, err: ", err)
	}
	fmt.Println(string(response))
}
