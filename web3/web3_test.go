package web3_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3"

	"github.com/cleanunicorn/ethereum/core"
	"github.com/cleanunicorn/ethereum/web3/account"
	"github.com/cleanunicorn/ethereum/web3/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

var update = flag.Bool("update", false, "update golden files")

func TestHTTPClient_Eth_getBalance(t *testing.T) {
	defer startGanache(t)()

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
		{
			name: "Account with Eth in block 0x0",
			fields: fields{
				HTTP: testGanacheHTTPEndpoint,
			},
			args: args{
				account: ganacheAccount9,
				block:   "0x0",
			},
			want:    "100000000000000000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := web3.NewClient(provider.DialHTTP(tt.fields.HTTP))
			got, err := c.Eth.GetBalance(tt.args.account, tt.args.block)
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
	type args struct {
		blockNumberHex      string
		includeTransactions bool
	}
	tests := []struct {
		name     string
		endpoint string
		args     args
		wantErr  bool
	}{
		{
			name:     "Get Mainnet block 1",
			endpoint: testMainnetHTTPEndpoint,
			args: args{
				blockNumberHex: "0x1",
			},
			wantErr: false,
		},
		{
			name:     "Get Mainnet block 5000000 with transaction data",
			endpoint: testMainnetHTTPEndpoint,
			args: args{
				blockNumberHex:      "0x4C4B40",
				includeTransactions: true,
			},
			wantErr: false,
		},
		{
			name:     "Get Mainnet block 5000000 without transaction data",
			endpoint: testMainnetHTTPEndpoint,
			args: args{
				blockNumberHex:      "0x4C4B40",
				includeTransactions: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := web3.NewClient(provider.DialHTTP(tt.endpoint))
			got, err := c.Eth.GetBlockByNumber(tt.args.blockNumberHex, tt.args.includeTransactions)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.Eth_getBlockByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			golden := filepath.Join("test-fixtures/", tt.name+".golden")
			if *update {
				gotJSON, err := json.MarshalIndent(got, "", "    ")
				if err != nil {
					t.Error("Error marshalling golden file, err: ", err)
				}
				ioutil.WriteFile(golden, gotJSON, 0644)
			}
			wantJSON, _ := ioutil.ReadFile(golden)

			var want types.Block
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
			c := web3.NewClient(provider.DialHTTP(tt.endpoint))
			got, err := c.Eth.BlockNumber()
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

func TestHTTPClient_Eth_getTransactionCount(t *testing.T) {
	defer startGanache(t)()

	type args struct {
		account string
		block   string
	}
	tests := []struct {
		name     string
		endpoint string
		args     args
		want     uint64
		wantErr  bool
	}{
		{
			name:     "Account " + zeroAccount + " should have 0 transactions",
			endpoint: testMainnetHTTPEndpoint,
			args: args{
				account: zeroAccount,
				block:   "latest",
			},
			want: 0,
		},
		{
			name:     "Account " + ganacheAccount0 + " should have 0 transactions",
			endpoint: testGanacheHTTPEndpoint,
			args: args{
				account: ganacheAccount0,
				block:   "latest",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := web3.NewClient(provider.DialHTTP(tt.endpoint))
			got, err := c.Eth.GetTransactionCount(tt.args.account, tt.args.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.Eth_getTransactionCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HTTPClient.Eth_getTransactionCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPClient_Eth_sendRawTransaction(t *testing.T) {
	defer startGanache(t)()

	type args struct {
		signedTransaction string
	}
	tests := []struct {
		name               string
		endpoint           string
		args               args
		account1PrivateKey string
		account2PrivateKey string
		wantErr            bool
	}{
		{
			name:               "Test account creates transaction to another account",
			endpoint:           testGanacheHTTPEndpoint,
			account1PrivateKey: "09b2e5a4cec476e891c8b2aae556399953c271f769e22d17554030c7a58b8d88",
			account2PrivateKey: "8fb1d9dcc5812a63339fa6fef45f204338a9a136be4afd522df5648790fe9cb5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := web3.NewClient(provider.DialHTTP(tt.endpoint))
			networkID, _ := c.Net.Version()
			signer := core.CreateSigner(networkID)
			a, _ := account.FromHexKey(tt.account1PrivateKey)
			b, _ := account.FromHexKey(tt.account2PrivateKey)

			nonce, err := c.Eth.GetTransactionCount(a.Address(), "latest")
			if err != nil {
				t.Errorf("Could not get nonce for account: %s err: %s", a.Address(), err)
			}

			tx, err := core.SignTx(
				signer,
				a,
				nonce,
				common.HexToAddress(b.Address()),
				big.NewInt(0),
				21000,
				big.NewInt(1),
				[]byte{},
			)
			txs := gethtypes.Transactions{tx}
			transactionHash := fmt.Sprintf("0x%x", txs.GetRlp(0))

			got, err := c.Eth.SendRawTransaction(transactionHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.Eth_sendRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 66 {
				t.Errorf("HTTPClient.Eth_sendRawTransaction() = %v, err: %s", got, err)
			}
		})
	}
}

func TestHTTPClient_Eth_getTransactionReceipt(t *testing.T) {
	defer startGanache(t)()

	type fields struct {
		client   *http.Client
		endpoint string
	}
	type args struct {
		transactionHash string
	}
	tests := []struct {
		name     string
		endpoint string
		args     args
		wantErr  bool
	}{
		{
			name:     "Get receipt from transaction hash contract call",
			endpoint: testMainnetHTTPEndpoint,
			args: args{
				transactionHash: "0x4c65570f9ceab8a0a575af2f500b83c7d8077d595e42dff4c1f90e53b05c9ae8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := web3.NewClient(provider.DialHTTP(tt.endpoint))
			got, err := c.Eth.GetTransactionReceipt(tt.args.transactionHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPClient.Eth_getTransactionReceipt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			golden := filepath.Join("test-fixtures/", tt.name+".golden")
			if *update {
				gotJSON, _ := json.MarshalIndent(got, "", "    ")
				ioutil.WriteFile(golden, gotJSON, 0644)
			}
			wantJSON, _ := ioutil.ReadFile(golden)
			var want types.Receipt
			if err := json.Unmarshal(wantJSON, &want); err != nil {
				t.Error("Could not unmarshal expected response, err: ", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("HTTPClient.Eth_getTransactionReceipt() = %v, want %v", got, want)
			}
		})
	}
}

// Ethereum node variables
const testMainnetHTTPEndpoint = "https://mainnet.infura.io"
const emptyAccount = "0x00000000000000000000000000000000000000ff"
const zeroAccount = "0x0000000000000000000000000000000000000000"

const ganachePort = "58545"
const testGanacheHTTPEndpoint = "http://localhost:" + ganachePort
const testGanacheNetworkID = 99
const ganacheAccount0 = "0xbd1e71ca74e8665718be94189a9e9f8ea07087d1"
const ganacheAccount1 = "0xc5bda996ac2d16ee24e0c9f69e44e12f79ddeff3"
const ganacheAccount2 = "0xc9e654a52918ead651a19491e1e5d6e2b7c54805"
const ganacheAccount3 = "0x91f2e144578057c0d27aaea35f81e932397f7162"
const ganacheAccount4 = "0x942921a92195aeb1a72798ee19edb11b05456d70"
const ganacheAccount5 = "0x11775df7a5541e7affa3c07c415a9a88e8f878bf"
const ganacheAccount6 = "0xe07b84dbd005dd48fd3632ec910630d96462c22e"
const ganacheAccount7 = "0x417f0549f0c8afbf8e440bfca721905adc51c128"
const ganacheAccount8 = "0x95355382c7d4bcae94df7f0ff1178b325bd7512b"
const ganacheAccount9 = "0x0ed774f495f902952dca2ce019241433c0088686"

func startGanache(t *testing.T) func() {
	// Make sure ganache is stopped
	ganacheStop := exec.Command("docker", []string{"rm", "-f", "ganache"}...)
	ganacheStop.Start()
	ganacheStop.Wait()

	// Start Ganache
	command := "docker"
	args := fmt.Sprintf("run --net=host --name=ganache --rm trufflesuite/ganache-cli -s 99 -i %d -p %s", testGanacheNetworkID, ganachePort)
	ganache := exec.Command(command, strings.Split(args, " ")...)
	ganacheOut, err := ganache.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := ganache.Start(); err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 10000)
	var output string
	for {
		n, err := ganacheOut.Read(buff)
		if err != nil {
			t.Fatal("Error reading output, err = ", err)
		}
		output = output + string(buff[:n])

		if strings.Contains(output, "Listening on") {
			break
		}
	}

	return func() {
		ganache.Process.Kill()
		rmContainer := exec.Command("docker", "rm", "-f", "ganache")
		rmContainer.Run()
	}
}
