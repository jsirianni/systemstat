package main

import (
	"os"
	"flag"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database"
)

func main() {
	flag.Parse()

	d, err := database.NewPostgres()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	server := database.Server{
		Port: 9000,
		DB: d,
	}

	if err := server.Run(); err != nil {
		log.Fatal(err, 1)
		os.Exit(1)
	}
}
