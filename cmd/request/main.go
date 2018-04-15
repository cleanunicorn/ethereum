package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/cleanunicorn/ethereum/client"
)

func main() {
	var (
		endpointHTTP = flag.String("http", "http://127.0.0.1:8545", "HTTP endpoint")
		method       = flag.String("method", "", "method")
	)
	flag.Parse()
	args := flag.Args()

	fmt.Println("endpointHTTP = ", *endpointHTTP)
	fmt.Println("Method = ", *method)
	fmt.Println("Args = ", args)

	c, err := client.DialHTTP(*endpointHTTP)
	if err != nil {
		log.Fatalf("Could not dial into HTTP endpoint: %s, err: %s", endpointHTTP, err)
	}
	log.Debug(c)

	response, err := c.MakeRawCall(*method, args)
	log.Println(response)
	log.Println(err)
}
