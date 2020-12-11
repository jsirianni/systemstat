SHELL := /bin/bash

ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

build: build.all

build.all: build.database
	go build ./cmd/agent
	go build ./cmd/cli
	go build ./cmd/frontend
	go build ./cmd/control
	go build ./cmd/alert

build.database:
	go build ./cmd/database
	docker build -t systemstat/database:latest . -f cmd/database/Dockerfile

test: shellcheck
	go test ./...

test.integration: clean test deploy.local
	scripts/service/database/test/test.sh
	scripts/service/frontend/test/test.sh
	scripts/service/control/test/test.sh
	scripts/service/alert/test/test.sh
	source ./test.env && go test ./... -tags=integration

deploy.local:
	docker-compose build --parallel
	docker-compose up -d && sleep 5
	docker exec systemstat_db_1 /var/lib/postgresql/systemstat/initdb.sh
	docker exec systemstat_db_1 /var/lib/postgresql/systemstat/test_data.sh

fmt:
	go fmt ./...

clean:
	docker-compose down --remove-orphans
	rm -f agent alert frontend cli control database

shellcheck:
	shellcheck scripts/postgres/initdb.sh
	shellcheck scripts/postgres/test_data.sh
	shellcheck scripts/service/database/test/test.sh
	shellcheck scripts/service/frontend/test/test.sh
	shellcheck scripts/service/control/test/test.sh
	shellcheck scripts/service/alert/test/test.sh
