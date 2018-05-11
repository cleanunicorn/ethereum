# ethereum : A set of Ethereum specific tools

## Overview [![GoDoc](https://godoc.org/github.com/cleanunicorn/ethereum?status.svg)](https://godoc.org/github.com/cleanunicorn/ethereum) [![Go Report Card](https://goreportcard.com/badge/github.com/cleanunicorn/ethereum)](https://goreportcard.com/report/github.com/cleanunicorn/ethereum) [![Sourcegraph](https://sourcegraph.com/github.com/cleanunicorn/ethereum/-/badge.svg)](https://sourcegraph.com/github.com/cleanunicorn/ethereum?badge)

## Install

```
go get github.com/cleanunicorn/ethereum/web3
```

## Use

Connect to an Ethereum endpoint and ask for a block number
```go
c := web3.NewClient(provider.DialHTTP("https://mainnet.infura.io/:8545"))

b, err := c.Eth.GetBlockByNumber("0x2", false)
if err != nil {
	fmt.Printf("Error getting latest block, err: %v", err)
	os.Exit(1)
}

fmt.Printf("Block: %#v\n", b)
```

Make raw calls
```go
c := web3.NewClient(provider.DialHTTP("https://mainnet.infura.io/:8545"))

res, err := c.Provider.Call("eth_blockNumber", nil)
if err != nil {
	fmt.Printf("Error getting latest block, err: %v", err)
	os.Exit(1)
}
fmt.Println(string(res))
```

Check [examples](https://github.com/cleanunicorn/ethereum/tree/restructure-namespaces/examples) for more sample code

Check the [documentation](https://godoc.org/github.com/cleanunicorn/ethereum) 

Implemented requests:

- [ ] web3_clientVersion                      
- [ ] web3_sha3                               
- [x] net_version                             
- [ ] net_peerCount                           
- [ ] net_listening                           
- [ ] eth_protocolVersion                     
- [ ] eth_syncing                             
- [ ] eth_coinbase                            
- [ ] eth_mining                              
- [ ] eth_hashrate                            
- [ ] eth_gasPrice                            
- [ ] eth_accounts                            
- [x] eth_blockNumber                         
- [x] eth_getBalance                          
- [ ] eth_getStorageAt (deprecated)
- [x] eth_getTransactionCount                 
- [ ] eth_getBlockTransactionCountByHash      
- [ ] eth_getBlockTransactionCountByNumber    
- [ ] eth_getUncleCountByBlockHash            
- [ ] eth_getUncleCountByBlockNumber          
- [ ] eth_getCode                             
- [ ] eth_sign                                
- [ ] eth_sendTransaction                     
- [x] eth_sendRawTransaction                  
- [ ] eth_call                                
- [ ] eth_estimateGas                         
- [ ] eth_getBlockByHash                      
- [x] eth_getBlockByNumber                    
- [ ] eth_getTransactionByHash                
- [ ] eth_getTransactionByBlockHashAndIndex   
- [ ] eth_getTransactionByBlockNumberAndIndex 
- [x] eth_getTransactionReceipt               
- [ ] eth_getUncleByBlockHashAndIndex         
- [ ] eth_getUncleByBlockNumberAndIndex       
- [ ] eth_getCompilers                        
- [ ] eth_compileLLL                          
- [ ] eth_compileSolidity (deprecated)                    
- [ ] eth_compileSerpent                      
- [ ] eth_newFilter                           
- [ ] eth_newBlockFilter                      
- [ ] eth_newPendingTransactionFilter         
- [ ] eth_uninstallFilter                     
- [ ] eth_getFilterChanges                    
- [ ] eth_getFilterLogs                       
- [ ] eth_getLogs                             
- [ ] eth_getWork                             
- [ ] eth_submitWork                          
- [ ] eth_submitHashrate                      
- [ ] db_putString                            
- [ ] db_getString                            
- [ ] db_putHex                               
- [ ] db_getHex                               
- [ ] shh_post                                
- [ ] shh_version                             
- [ ] shh_newIdentity                         
- [ ] shh_hasIdentity                         
- [ ] shh_newGroup                            
- [ ] shh_addToGroup                          
- [ ] shh_newFilter                           
- [ ] shh_uninstallFilter                     
- [ ] shh_getFilterChanges                    
- [ ] shh_getMessages                         
- [ ] personal_listAccounts                   
- [ ] personal_newAccount                     
- [ ] personal_sendTransaction                
- [ ] personal_unlockAccount                  


## Author

Daniel Luca

## License

MIT.
