SHELL := /bin/bash

ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

build: build.all

build.all:
	go build -o bin/ ./cmd/agent
	go build -o bin/ ./cmd/cli
	go build -o bin/ ./cmd/frontend
	go build -o bin/ ./cmd/control
	go build -o bin/ ./cmd/alert
	go build -o bin/ ./cmd/database

test: shellcheck
	go test ./...

test.integration: clean test build.all deploy.local
	scripts/service/database/test/test.sh
	scripts/service/frontend/test/test.sh
	scripts/service/control/test/test.sh
	scripts/service/alert/test/test.sh
	source ./test.env && go test ./... -tags=integration

deploy.local:
	docker-compose build --parallel
	scripts/deploy/local/docker-compose.sh
	docker exec systemstat_postgres_1 /var/lib/postgresql/systemstat/initdb.sh
	docker exec systemstat_postgres_1 /var/lib/postgresql/systemstat/test_data.sh

fmt:
	go fmt ./...

clean:
	docker-compose down --remove-orphans
	rm -f bin/*

shellcheck:
	shellcheck scripts/postgres/initdb.sh
	shellcheck scripts/postgres/test_data.sh
	shellcheck scripts/service/database/test/test.sh
	shellcheck scripts/service/frontend/test/test.sh
	shellcheck scripts/service/control/test/test.sh
	shellcheck scripts/service/alert/test/test.sh
	shellcheck scripts/deploy/local/docker-compose.sh

protobuf.generate: protobuf.generate.database

protobuf.generate.database:
	cd internal/service/database && \
		protoc \
			--go_out=. \
			--go_opt=paths=source_relative \
    		--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
    		./api.proto
