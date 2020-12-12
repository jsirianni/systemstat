package database

import (
    "github.com/jsirianni/systemstat/api"
)

type Server struct {
	Port struct {
		HTTP int
		GRPC int
	}
	DB   Database

	api.UnimplementedApiServer
}
