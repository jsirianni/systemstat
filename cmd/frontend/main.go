package main

import (
	"flag"
	"os"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/frontend"
)

const defaultPort = 9090
const defaultGRPCPort = 9091
const defaultEndpoint = "database:9100"

var port int
var grpcPort int
var endpoint string

func main() {
	flag.IntVar(&port, "port", defaultPort, "port to use for http server")
	flag.IntVar(&grpcPort, "grpc-port", defaultGRPCPort, "port to use for grpc server")
	flag.StringVar(&endpoint, "database=grpc", defaultEndpoint, "database grpc endpoint")
	flag.Parse()

	server := frontend.Server{}
	server.Port.HTTP = port
	server.Port.GRPC = grpcPort
	server.Database.Endpoint = endpoint

	if err := server.Server(); err != nil {
		log.Fatal(err, 1)
		os.Exit(1)
	}
}
