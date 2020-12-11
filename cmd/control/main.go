package main

import (
	"flag"
	"os"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/control"
)

const defaultPort = 9080
const defaultEndpoint = "http://database:9000"

var port int
var endpoint string

func main() {
	flag.IntVar(&port, "port", defaultPort, "port to use for http server")
	flag.StringVar(&endpoint, "endpoint", defaultEndpoint, "backend database api endpoint")
	flag.Parse()

	server := control.Server{}
	server.Port = port
	server.Database.Endpoint = endpoint

	if err := server.Run(); err != nil {
		log.Fatal(err, 1)
		os.Exit(1)
	}
}
