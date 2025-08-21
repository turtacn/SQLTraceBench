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

plugins: build/starrocks build/clickhouse

build/starrocks:
	go build -buildmode=plugin -o plugins/starrocks.so ./plugins/starrocks

build/clickhouse:
	go build -buildmode=plugin -o plugins/clickhouse.so ./plugins/clickhouse
