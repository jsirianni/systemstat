SHELL := /bin/bash

ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

build: build.all
build.all: build.agent build.cli build.frontend build.control build.alert build.database
build.agent:
	go build -o bin/ ./cmd/agent
build.cli:
	go build -o bin/ ./cmd/cli
build.frontend:
	go build -o bin/ ./cmd/frontend
build.control:
	go build -o bin/ ./cmd/control
build.alert:
	go build -o bin/ ./cmd/alert
build.database:
	go build -o bin/ ./cmd/database

test: shellcheck
	go test ./...
test.scripts:
	scripts/service/database/test/test.sh
	scripts/service/frontend/test/test.sh
	scripts/service/control/test/test.sh
	scripts/service/alert/test/test.sh
test.integration: clean test deploy.local test.scripts
	source ./test.env && go test ./... -tags=integration

docker-compose.build:
	docker-compose build --parallel
docker-compose.deploy:
	scripts/deploy/local/docker-compose.sh

deploy.local: docker-compose.build docker-compose.deploy
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
