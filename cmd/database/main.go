package main

import (
	"os"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database"
)

func main() {
	d, err := database.NewPostgres()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	d.TestConnection()
	log.Info("working!")
}
