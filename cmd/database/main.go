package main

import (
	"os"
	"flag"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database"
)

const defaultPort = 9000
var port int

func main() {
	flag.IntVar(&port, "port", defaultPort, "port to use for http server")
	flag.Parse()

	d, err := database.NewPostgres()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	server := database.Server{
		Port: port,
		DB: d,
	}

	if err := server.Run(); err != nil {
		log.Fatal(err, 1)
		os.Exit(1)
	}
}
