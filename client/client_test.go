package client_test

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

	"github.com/cleanunicorn/ethereum/client"
	"github.com/cleanunicorn/ethereum/core"
	"github.com/cleanunicorn/ethereum/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

var update = flag.Bool("update", false, "update golden files")

//
const testMainnetHTTPEndpoint = "https://mainnet.infura.io:8545"
const emptyAccount = "0x00000000000000000000000000000000000000ff"
const zeroAccount = "0x0000000000000000000000000000000000000000"

//
const testGanacheHTTPEndpoint = "http://localhost:8545"
const testGanacheNetworkID = 99
const ganacheAccount0 = "0xd84cf7a5a3c7985398c591bc61662b8be438dab8"
const ganacheAccount1 = "0xb2622b59630e294578852c3d591e87dcb6507037"
const ganacheAccount2 = "0x534695bacbdf0428ec1b7cec15eccc88680959bd"
const ganacheAccount3 = "0x32692e49169679212038d207ef38cc045af55244"
const ganacheAccount4 = "0xe87b5397e46e4960256191defe9ae40029e9875b"
const ganacheAccount5 = "0xc0b8d4d03aadb9c82b468620a948c75628719ddd"
const ganacheAccount6 = "0x273040d239bc002ccdebf6c5d69e4a751731402b"
const ganacheAccount7 = "0x9f594debb9d77a77e5a63477f0c5e5d3464c208e"
const ganacheAccount8 = "0xb2784b58c227c8968fc78ae4829c87996d497d70"
const ganacheAccount9 = "0x5a70a6b58d20c95a3fd29a0ee046412cddc98bde"

func startGanache(t *testing.T) func() {
	// Make sure ganache is stopped
	ganacheStop := exec.Command("docker", []string{"rm", "-f", "ganache"}...)
	ganacheStop.Start()
	ganacheStop.Wait()

	// Start Ganache
	command := "docker"
	args := "run --net=host --name=ganache --rm trufflesuite/ganache-cli -s 99 -d 0 -i " + fmt.Sprintf("%d", testGanacheNetworkID)
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
			t.Fatal("Error reading output")
		}
		output = output + string(buff[:n])

		if strings.Contains(output, "Listening on localhost:8545") {
			break
		}
	}

	return func() {
		ganache.Process.Kill()
		rmContainer := exec.Command("docker", "rm", "-f", "ganache")
		rmContainer.Run()
	}
}

func TestHTTPClient_Net_version(t *testing.T) {
	defer startGanache(t)()

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
		{
			name:     "Ganache ID should be 99",
			endpoint: testGanacheHTTPEndpoint,
			want:     testGanacheNetworkID,
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
			c, err := client.DialHTTP(tt.fields.HTTP)
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
			c, err := client.DialHTTP(tt.endpoint)
			got, err := c.Eth_getBlockByNumber(tt.args.blockNumberHex, tt.args.includeTransactions)
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
			c, err := client.DialHTTP(tt.endpoint)
			got, err := c.Eth_getTransactionCount(tt.args.account, tt.args.block)
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
			account1PrivateKey: "5905ed74bb339cf0f456020ecd63415d80588f234ffcffca4fe119b13b8ef32a",
			account2PrivateKey: "d1ecb25acf8387b949e50809ceedc47abfeeca1e04a8ddfb083f3aebe6d5e680",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := client.DialHTTP(tt.endpoint)

			networkID, _ := c.Net_version()
			signer := core.CreateSigner(networkID)
			a, _ := types.AccountFromHexKey(tt.account1PrivateKey)
			b, _ := types.AccountFromHexKey(tt.account2PrivateKey)

			nonce, err := c.Eth_getTransactionCount(a.Address(), "latest")
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

			got, err := c.Eth_sendRawTransaction(transactionHash)
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
			c, err := client.DialHTTP(
				tt.endpoint,
			)
			got, err := c.Eth_getTransactionReceipt(tt.args.transactionHash)
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
			var want client.Receipt
			if err := json.Unmarshal(wantJSON, &want); err != nil {
				t.Error("Could not unmarshal expected response, err: ", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("HTTPClient.Eth_getTransactionReceipt() = %v, want %v", got, want)
			}
		})
	}
}
