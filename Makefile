BIN := "./bin/antibruteforce"
DOCKER_IMG="antibruteforce:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/antibruteforce

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -run TestCreateBucket  -v -count=1 -race -timeout=1m internal/storage/storage_test.go
	go test -run TestHelloWorldHandler  -v -count=1 -race -timeout=1m internal/server/http/server_test.go

int-tests:
	docker rm -f test-postgres
	docker run --name test-postgres -p 5432:5432 -e POSTGRES_PASSWORD=dbpass -d postgres
	sleep 3
	goose up
	go test ./... -v -count=1 -race -timeout=1m .


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img test lint
