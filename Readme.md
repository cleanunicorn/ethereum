# ethereum : A set of Ethereum specific tools

## Overview [![GoDoc](https://godoc.org/github.com/cleanunicorn/ethereum?status.svg)](https://godoc.org/github.com/cleanunicorn/ethereum) [![Go Report Card](https://goreportcard.com/badge/github.com/cleanunicorn/ethereum)](https://goreportcard.com/report/github.com/cleanunicorn/ethereum) [![Sourcegraph](https://sourcegraph.com/github.com/cleanunicorn/ethereum/-/badge.svg)](https://sourcegraph.com/github.com/cleanunicorn/ethereum?badge)

## Install

```
go get github.com/cleanunicorn/ethereum
```

## Use

```
package main

import (
	"fmt"
	"os"

	"github.com/cleanunicorn/ethereum/client"
)

func main() {
	c, err := client.DialHTTP("https://mainnet.infura.io")
	if err != nil {
		fmt.Printf("Could not dial into HTTP, err: %s", err)
		os.Exit(1)
	}

	// You can make raw calls
	resp, err := c.RawCall("eth_getBlockByNumber", []interface{}{"0x1", true})
	if err != nil {
		fmt.Printf("Error making a raw call to eth_getBlockByNumber, err: %s", err)
		os.Exit(1)
	}
	fmt.Println("Block number 1 with included transactions")
	fmt.Println(string(resp))
	// Returns
	//
	// {
	// 	"jsonrpc":"2.0",
	// 	"id":1,
	// 	"result":{
	// 		"difficulty":"0x3ff800000",
	// 		"extraData":"0x476574682f76312e302e302f6c696e75782f676f312e342e32",
	// 		"gasLimit":"0x1388",
	// 		"gasUsed":"0x0",
	// 		"hash":"0x88e96d4537bea4d9c05d12549907b32561d3bf31f45aae734cdc119f13406cb6",
	// 		"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	// 		"miner":"0x05a56e2d52c817161883f50c441c3228cfe54d9f",
	// 		"mixHash":"0x969b900de27b6ac6a67742365dd65f55a0526c41fd18e1b16f1a1215c2e66f59",
	// 		"nonce":"0x539bd4979fef1ec4",
	// 		"number":"0x1",
	// 		"parentHash":"0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3",
	// 		"receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	// 		"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	// 		"size":"0x219",
	// 		"stateRoot":"0xd67e4d450343046425ae4271474353857ab860dbc0a1dde64b41b5cd3a532bf3",
	// 		"timestamp":"0x55ba4224",
	// 		"totalDifficulty":"0x7ff800000",
	// 		"transactions":[],
	// 		"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	// 		"uncles":[]
	// 	}
	// }

	// Or make specific requests
	// Like getting block number 1000000
	block, err := c.Eth_getBlockByNumber(fmt.Sprintf("0x%x", 1000000), true)
	if err != nil {
		fmt.Printf("Error calling eth_blockByNumber directly, err: %s", err)
		os.Exit(1)
	}
	fmt.Println("Block number 1000000 with include transactions")
	fmt.Printf("%+v", block)
}
```

Implemented requests:

- [x] [eth_getTransactionCount](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount)
- [x] [eth_getBalance](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount)
- [x] [eth_sendRawTransaction](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction)
- [x] [eth_getBlockByNumber](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber)
- [x] [eth_blockNumber](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber)
- [x] [net_version](https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version) 

## Author

Daniel Luca

## License

MIT.
