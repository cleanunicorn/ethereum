package net

import (
	"encoding/json"

	"github.com/cleanunicorn/ethereum/provider"
)

type Net struct {
	provider provider.Provider
}

func NewNet(p provider.Provider) Net {
	return Net{
		provider: p,
	}
}

// Version returns the network id.
//
// See https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
func (c Net) Version() (int64, error) {
	reply, err := c.provider.Call("net_version", []interface{}{})
	if err != nil {
		return 0, err
	}

	type responseNetVersion struct {
		ID      int    `json:"id"`
		Jsonrpc string `json:"jsonrpc"`
		Result  int64  `json:"result,string"`
	}
	var networkIDReply responseNetVersion

	err = json.Unmarshal(reply, &networkIDReply)
	if err != nil {
		return 0, err
	}

	return networkIDReply.Result, nil
}
