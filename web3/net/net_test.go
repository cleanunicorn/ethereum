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
const testMainnetHTTPEndpoint = "https://mainnet.infura.io:8545"
const ganachePort = "58545"
const testGanacheHTTPEndpoint = "http://localhost:" + ganachePort
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
	args := fmt.Sprintf("run --net=host --name=ganache --rm trufflesuite/ganache-cli -s 99 -d 0 -i %d -p %s", testGanacheNetworkID, ganachePort)
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
