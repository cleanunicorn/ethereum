package net

import (
	"encoding/json"

	"github.com/cleanunicorn/ethereum/provider"
)

// Net module
type Net struct {
	provider provider.Provider
}

// NewNet returns an instance of net module
func NewNet(p provider.Provider) Net {
	return Net{
		provider: p,
	}
}

// ResponseNetVersion is the structure returned by https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
type ResponseNetVersion struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  int64  `json:"result,string"`
}

// Version returns the network id
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
func (c Net) Version() (int64, error) {
	reply, err := c.provider.Call("net_version", []interface{}{})
	if err != nil {
		return 0, err
	}

	var networkIDReply ResponseNetVersion
	err = json.Unmarshal(reply, &networkIDReply)
	if err != nil {
		return 0, err
	}

	return networkIDReply.Result, nil
}
