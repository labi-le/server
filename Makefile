include .env
export $(shell sed 's/=.*//' .env)

PROJ_NAME = labile-me-serv

MAIN_PATH = cmd/main.go
BUILD_PATH = build/package/

DB_NAME = db

run-build:
	./$(BUILD_PATH)$(PROJ_NAME)

run:
	go run -v $(MAIN_PATH) -debug
	#go run -v $(MAIN_PATH)

# https://github.com/maxcnunes/gaper
watch:
	gaper \
    	--build-path cmd \
    	--watch .

resolve-deps:
	go get -d ./...

build-release: build-default

# пытается собрать статический бинарник
build-default: clean resolve-deps
	go build --ldflags '-extldflags "-static"' -v -o $(BUILD_PATH)$(PROJ_NAME) $(MAIN_PATH)

# меньше размером (возможны ошибки)
build-static: clean
	go build -ldflags "-w -linkmode external -extldflags "-static" -s" -v -o $(BUILD_PATH)$(PROJ_NAME) $(MAIN_PATH)

migrate-up:install-migrate-tool
	#migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" up
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-create:install-migrate-tool
	@migrate create -ext sql -dir ./migrations $(filter-out $@,$(MAKECMDGOALS))

install-migrate-tool:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

clean:
	rm -rf $(BUILD_PATH)*

tests:
	go test -cover -v ./...

lint:
	golangci-lint run