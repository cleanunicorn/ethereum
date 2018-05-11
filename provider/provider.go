package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Provider encapsulates the client that will be used for calls
type Provider struct {
	HTTPClient   *http.Client
	HTTPEndpoint string
}

// DialHTTP takes a HTTP endpoint and returns the HTTPClient structure.
func DialHTTP(endpointHTTP string) Provider {
	var httpTransport = &http.Transport{}
	var httpClient = &http.Client{Transport: httpTransport}

	p := Provider{
		HTTPClient:   httpClient,
		HTTPEndpoint: endpointHTTP,
	}

	return p
}

// Call makes a request with a specified method and parameters.
func (c Provider) Call(method string, params interface{}) ([]byte, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest("POST", c.HTTPEndpoint, strings.NewReader(string(dataJSON)))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	request.Header.Add("Content-Type", "application/json")

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// Check if the result is an error
	type responseError struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	var respErr responseError
	if err := json.Unmarshal(body, &respErr); err != nil {
		return []byte{}, err
	}
	if respErr.Error.Code != 0 {
		return []byte{}, fmt.Errorf("code: %d, error: %s", respErr.Error.Code, respErr.Error.Message)
	}

	return body, nil
}

// RawCall calls a method with a JSON encoded list of params.
//
// method should be a string
//    ex: "eth_getBlockByNumber"
//
// params is an interface
//    ex: {"0x1", true}
//
func (c Provider) RawCall(method string, args []interface{}) ([]byte, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  args,
		"id":      1,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest("POST", c.HTTPEndpoint, strings.NewReader(string(dataJSON)))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	request.Header.Add("Content-Type", "application/json")

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// Check if the result is an error
	type responseError struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	var respErr responseError
	if err := json.Unmarshal(body, &respErr); err != nil {
		return []byte{}, err
	}
	if respErr.Error.Code != 0 {
		return []byte{}, fmt.Errorf("code: %d, error: %s", respErr.Error.Code, respErr.Error.Message)
	}

	return body, nil
}
