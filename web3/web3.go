package web3

import (
	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3/eth"
	"github.com/cleanunicorn/ethereum/web3/net"
)

// Default parameters
const defaultHTTP = "http://127.0.0.1:8545"

type Client struct {
	Provider provider.Provider
	Eth      eth.Eth
	Net      net.Net
}

func NewClient(p provider.Provider) Client {
	c := Client{
		Provider: p,
	}

	c.Eth = eth.NewEth(p)
	c.Net = net.NewNet(p)

	return c
}
