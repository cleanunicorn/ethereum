# ethereum : A set of Ethereum specific tools

## Overview [![GoDoc](https://godoc.org/github.com/cleanunicorn/ethereum?status.svg)](https://godoc.org/github.com/cleanunicorn/ethereum) [![Go Report Card](https://goreportcard.com/badge/github.com/cleanunicorn/ethereum)](https://goreportcard.com/report/github.com/cleanunicorn/ethereum) [![Sourcegraph](https://sourcegraph.com/github.com/cleanunicorn/ethereum/-/badge.svg)](https://sourcegraph.com/github.com/cleanunicorn/ethereum?badge)

## Install

```
go get github.com/cleanunicorn/ethereum
```

## Use

Connect to an Ethereum endpoint
```go
c, err := client.DialHTTP("https://mainnet.infura.io")
if err != nil {
	fmt.Printf("Could not dial into HTTP, err: %s", err)
}
```

Make raw calls
```go
resp, err := c.RawCall("eth_getBlockByNumber", []interface{}{"0x1", true})
if err != nil {
	fmt.Printf("Error making a raw call to eth_getBlockByNumber, err: %s", err)
} else {
	fmt.Println(string(resp))
}
```

Or make specific requests like getting block number 1000000
```go
block, err := c.Eth_getBlockByNumber(fmt.Sprintf("0x%x", 1000000), true)
if err != nil {
	fmt.Printf("Error calling eth_blockByNumber directly, err: %s", err)
} else {
	fmt.Printf("%+v", block)
}
```

Check the [documentation](https://godoc.org/github.com/cleanunicorn/ethereum) 

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
