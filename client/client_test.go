package client_test

// import (
// 	"fmt"
// 	"math/big"
// 	"reflect"
// 	"testing"

// 	client "github.com/cleanunicorn/ethereum/client"
// 	"github.com/ethereum/go-ethereum/common"
// 	gethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"gitlab.com/cleanunicorn/eth-tipper/core"
// 	"gitlab.com/cleanunicorn/eth-tipper/core/types"
// )

// var testServerHTTP = "http://127.0.0.1:8545"

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

// 	if len(tHash) != 66 {
// 		t.Error("Expecting hash, got:", tHash)
// 	}
// }

// func TestJSONRPCEthereumServer_Net_version(t *testing.T) {
// 	type fields struct {
// 		url string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		want    int64
// 		wantErr bool
// 	}{
// 		{
// 			name: "Network ID should be 99",
// 			fields: fields{
// 				url: testServerHTTP,
// 			},
// 			want: 99,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := client.Client{
// 				HTTP: tt.fields.url,
// 			}
// 			got, err := s.Net_version()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("JSONRPCEthereumServer.Net_version() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("JSONRPCEthereumServer.Net_version() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_Eth_getBalance(t *testing.T) {
// 	type fields struct {
// 		HTTP string
// 	}
// 	type args struct {
// 		account string
// 		block   string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name: "Should be untouched in the latest block",
// 			fields: fields{
// 				HTTP: testServerHTTP,
// 			},
// 			args: args{
// 				account: "0x5a70a6b58d20c95a3fd29a0ee046412cddc98bde",
// 				block:   "latest",
// 			},
// 			want:    "100000000000000000000",
// 			wantErr: false,
// 		},
// 		{
// 			name: "Should be untouched in block 0x0",
// 			fields: fields{
// 				HTTP: testServerHTTP,
// 			},
// 			args: args{
// 				account: "0x5a70a6b58d20c95a3fd29a0ee046412cddc98bde",
// 				block:   "0x0",
// 			},
// 			want:    "100000000000000000000",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &client.Client{
// 				HTTP: tt.fields.HTTP,
// 			}
// 			got, err := s.Eth_getBalance(tt.args.account, tt.args.block)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.Eth_getBalance() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			want, _ := big.NewInt(0).SetString(tt.want, 10)
// 			if !reflect.DeepEqual(got, want) {
// 				t.Errorf("Client.Eth_getBalance() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
