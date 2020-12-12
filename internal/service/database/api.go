package database

import (
    "github.com/jsirianni/systemstat/internal/proto"
)

type Server struct {
	Port struct {
		HTTP int
		GRPC int
	}
	DB   Database

	proto.UnimplementedApiServer
}
