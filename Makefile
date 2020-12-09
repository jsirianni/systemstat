ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

build: build.all

build.all: build.database
	go build ./cmd/agent
	go build ./cmd/cli
	go build ./cmd/api
	go build ./cmd/control
	go build ./cmd/alert

build.database:
	go build ./cmd/database
	docker build -t systemstat/database:latest . -f cmd/database/Dockerfile

test:
	go test ./...

test.integration: test deploy.local
	. ./test.env
	echo $DATABASE_HOST
	go test ./... -tags=integration

deploy.local:
	docker-compose build
	docker-compose up -d && sleep 5
	docker exec systemstat_db_1 /var/lib/postgresql/systemstat/initdb.sh

fmt:
	go fmt ./...

clean:
	docker-compose down --remove-orphans
