package main

import (
	"flag"
	"time"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database"
)

const defaultPort = 9000
const defaultGRPCPort = 9100

var port int
var grpcPort int

func main() {
	flag.IntVar(&port, "port", defaultPort, "port to use for http server")
	flag.IntVar(&grpcPort, "grpc-port", defaultGRPCPort, "port to use for grpc server")
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Parse()

	d, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err, 100)
	}

	server := database.Server{}
	server.Port.HTTP = port
	server.Port.GRPC = grpcPort
	server.DB = d


	go startHTTP(&server)
	go startGRPC(&server)

	// TODO: servers should return errors over a channel, and handle exiting
	// in main ??
	for {
		if err := d.TestConnection(); err != nil {
			log.Error(err)
		} else {
			log.Trace("database test connection passed")
		}
		time.Sleep(time.Second * 10)
	}
}

func startHTTP(server *database.Server) {
	if err := server.RunHTTP(); err != nil {
		log.Fatal(err, 200)
	}
}

func startGRPC(server *database.Server) {
	if err := server.RunGRPC(); err != nil {
		log.Fatal(err, 300)
	}
}
