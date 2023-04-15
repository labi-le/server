CONTOUR = dev

include .env.$(CONTOUR)
export $(shell sed 's/=.*//' .env.$(CONTOUR))

PACKAGE = labile-serv

MAIN_PATH = cmd/main.go
BUILD_PATH = build/package/


VERSION=$(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2>/dev/null | sed 's/^.//')
COMMIT_HASH=$(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=-ldflags="-X '${PACKAGE}/internal.Version=${VERSION}' \
                   -X '${PACKAGE}/internal.CommitHash=${COMMIT_HASH}' \
                   -X '${PACKAGE}/internal.BuildTime=${BUILD_TIMESTAMP}' \
                   -extldflags '-static'"

.DEFAULT_GOAL := build
build: clean tests
	go build ${LDFLAGS} -v -o $(BUILD_PATH)$(PACKAGE) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

clean:
	@rm -rf $(BUILD_PATH)*

tests:
	@go test ./...

lint:
	golangci-lint run

migrate-up:install-migrate-tool
	./migrate -path ./migrations -database sqlite3://$(DB_NAME) up

migrate-create:install-migrate-tool
	./migrate create -ext sql -dir ./migrations $(filter-out $@,$(MAKECMDGOALS))

install-migrate-tool:
	@[ -f ./migrate ] \
		|| curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
		| tar zxvf - migrate

cloc:
	cloc . --exclude-ext=yml,mod,sum,xml