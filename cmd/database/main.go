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
	a, err := d.AccountByEmail("")
	if err != nil {
		panic(err)
	}
	accountString, err := a.String()
	if err != nil {
		panic(err)
	}
	alertConfigString, err := a.AlertConfig.String()
	if err != nil {
		panic(err)
	}

	log.Info(accountString)
	log.Info(alertConfigString)
}
