package main

import (
	"flag"
	"time"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database"
)

var grpcPort int

func main() {
	flag.IntVar(&grpcPort, "grpc-port", 9100, "port to use for grpc server")
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Parse()

	d, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err, 100)
	}

	server := database.Server{}
	server.Port.GRPC = grpcPort
	server.DB = d

	go startGRPC(&server)

	for {
		if err := d.TestConnection(); err != nil {
			log.Error(err)
		} else {
			log.Trace("database test connection passed")
		}
		time.Sleep(time.Second * 10)
	}
}

func startGRPC(server *database.Server) {
	if err := server.RunGRPC(); err != nil {
		log.Fatal(err, 300)
	}
}
