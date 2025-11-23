.PHONY: build test clean lint

APP := sqltracebench
LDFLAGS := -w -s -X github.com/turtacn/SQLTraceBench/cmd.Version=$(shell cat VERSION)

build: clean
	go build -ldflags '$(LDFLAGS)' -o bin/$(APP) ./cmd/sql_trace_bench

clean:
	rm -rf bin/

test:
	go test ./...

lint:
	golangci-lint run

plugins: build-plugin-starrocks

build-plugin-starrocks:
	go build -o bin/sqltracebench-starrocks-plugin ./cmd/plugins/starrocks/main.go
