package client_test

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/big"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cleanunicorn/ethereum/client"
)

var update = flag.Bool("update", false, "update golden files")

var testMainnetHTTPEndpoint = "https://mainnet.infura.io:8545"

// func TestGetTransactionCount(t *testing.T) {
// 	s := client.Client{
// 		HTTP: testServerHTTP,
// 	}

// 	// TODO: add an account with transaction number > 0
// 	transactionCount, err := s.Eth_getTransactionCount("0xfB8ab195c0134B6c809b176B5d829aC2e058e6b4", "latest")
// 	if err != nil {
// 		t.Error("Received error while getting transaction count", err)
// 	}

// 	if transactionCount != 0 {
// 		t.Error("Error getting transaction count \texpected:", 0, "\tgot:", transactionCount)
// 	}
// }

// func TestGetBalance(t *testing.T) {
// 	s := client.Client{
// 		HTTP: testServerHTTP,
// 	}

// 	// TODO: add an account with balance > 0
// 	balance, err := s.Eth_getBalance("0xfB8ab195c0134B6c809b176B5d829aC2e058e6b4", "latest")
// 	if err != nil {
// 		t.Error("Received error while getting balance", err)
// 	}

// 	if balance.Cmp(big.NewInt(0)) != 0 {
// 		t.Error("Error getting balance \texpected:", 0, "\tgot:", balance)
// 	}
// }

// func TestSendSignedTransaction(t *testing.T) {
// 	s := client.Client{
// 		HTTP: testServerHTTP,
// 	}

// 	a, _ := types.AccountFromHexKey("5905ed74bb339cf0f456020ecd63415d80588f234ffcffca4fe119b13b8ef32a")
// 	b, _ := types.AccountFromHexKey("d1ecb25acf8387b949e50809ceedc47abfeeca1e04a8ddfb083f3aebe6d5e680")

// 	signer := core.CreateSigner(99)

// 	nonce, err := s.Eth_getTransactionCount("0xd84cf7a5a3c7985398c591bc61662b8be438dab8", "latest")
// 	if err != nil {
// 		t.Error("Could not get account nonce err: ", err)
// 	}
// 	tx, err := core.SignTx(signer, a, nonce, common.HexToAddress(b.Address()), big.NewInt(0), 21000, big.NewInt(1), []byte{})

// 	txs := gethtypes.Transactions{tx}
// 	txH := fmt.Sprintf("0x%x", txs.GetRlp(0))

// 	tHash, err := s.Eth_sendRawTransaction(txH)
// 	if err != nil {
// 		t.Error("Error sending signed transaction", err)
// 	}
// Network
// 	if len(tHash) != 66 {
// 		t.Error("Expecting hash, got:", tHash)
// 	}
// }

func TestHTTPClient_Net_version(t *testing.T) {
	type fields struct {
		url string
	}
	tests := []struct {
		name     string
		endpoint string
		want     int64
		wantErr  bool
	}{
		{
			name:     "Mainnet ID should be 1",
			endpoint: testMainnetHTTPEndpoint,
			want:     1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := client.DialHTTP(tt.endpoint)
			got, err := c.Net_version()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONRPCEthereumServer.Net_version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JSONRPCEthereumServer.Net_version() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPClient_Eth_getBalance(t *testing.T) {
	const emptyAccount = "0x00000000000000000000000000000000000000ff"

	type fields struct {
		HTTP string
	}
	type args struct {
		account string
		block   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Should be untouched in the latest block",
			fields: fields{
				HTTP: testMainnetHTTPEndpoint,
			},
			args: args{
				account: emptyAccount,
				block:   "latest",
			},
			want:    "0",
			wantErr: false,
		},
		{
			name: "Should be untouched in block 0x0",
			fields: fields{
				HTTP: testMainnetHTTPEndpoint,
			},
			args: args{
				account: emptyAccount,
				block:   "0x0",
			},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := client.DialHTTP(testMainnetHTTPEndpoint)
			got, err := c.Eth_getBalance(tt.args.account, tt.args.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Eth_getBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want, _ := big.NewInt(0).SetString(tt.want, 10)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Client.Eth_getBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPClient_Eth_getBlockByNumber(t *testing.T) {
	type fields struct {
		client   *http.Client
		endpoint string
	}
	type args struct {
		blockNumberHex      string
		includeTransactions bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Get block 1",
			args: args{
				blockNumberHex: "0x1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := client.DialHTTP(testMainnetHTTPEndpoint)
			got, err := c.Eth_getBlockByNumber(tt.args.blockNumberHex, tt.args.includeTransactions)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.Eth_getBlockByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			golden := filepath.Join("test-fixtures/", tt.name+".golden")
			if *update {
				gotJSON, _ := json.MarshalIndent(got, "", "    ")
				ioutil.WriteFile(golden, gotJSON, 0644)
			}
			wantJSON, _ := ioutil.ReadFile(golden)
			var want client.Block
			if err := json.Unmarshal(wantJSON, &want); err != nil {
				t.Error("Could not unmarshal expected response, err: ", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("HTTPClient.Eth_getBlockByNumber() = %v, want %v", got, want)
			}
		})
	}
}

func TestHTTPClient_Eth_blockNumber(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		wantErr  bool
	}{
		{
			name:     "Block number should be greater than 0",
			endpoint: testMainnetHTTPEndpoint,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := client.DialHTTP(tt.endpoint)
			got, err := c.Eth_blockNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.Eth_blockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cmp := got.Cmp(big.NewInt(1)); cmp != 1 {
				t.Errorf("HTTPClient.Eth_blockNumber() = %v, cmp = %d", got, cmp)
			}
		})
	}
}
