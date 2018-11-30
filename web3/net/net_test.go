package net_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/cleanunicorn/ethereum/provider"
	"github.com/cleanunicorn/ethereum/web3"
)

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
			c := web3.NewClient(provider.DialHTTP(tt.endpoint))
			got, err := c.Net.Version()
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

// Ethereum node variables
const testMainnetHTTPEndpoint = "https://mainnet.infura.io"
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
