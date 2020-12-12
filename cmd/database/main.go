package main

import (
	"os"
	"fmt"
	"net"
	"flag"
	"time"
	"strconv"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	go startHTTP(&d)
	go startGRPC(&d)

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

func startHTTP(d *database.Database) {
	server := database.Server{
		Port: port,
		DB:   *d,
	}

	if err := server.Run(); err != nil {
		log.Fatal(err, 200)
	}
}

func startGRPC(d *database.Database) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
	  log.Fatal(err, 1)
	}
	var opts []grpc.ServerOption
	/*if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}*/

	grpcServer := grpc.NewServer(opts...)

	// allow grpcurl https://github.com/fullstorydev/grpcurl
	if os.Getenv("GO_ENV") == "development" {
		log.Trace("GO_ENV=development detected, enabling GRPC reflection")
		reflection.Register(grpcServer)
	}
	server := database.Server{
		DB: *d,
	}
	database.RegisterApiServer(grpcServer, server)
	log.Info("starting grpc api on port:", strconv.Itoa(grpcPort))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err, 1)
	}
}
