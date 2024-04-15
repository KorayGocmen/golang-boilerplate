.PHONY: all build build_linux clean run test

all: clean build run 

clean:
	rm -f ./bin/*

build:
	go build -ldflags "-X main.SHASUM=$(shell eval git rev-parse --short HEAD)" -o bin/api cmd/api/main.go

run:
	ENV_FILES=.env bin/api serve

test:
	go test -v ./...

db_reset:
	ENV_FILES=.env bin/api db reset

db_list:
	ENV_FILES=.env bin/api db list

db_up:
	ENV_FILES=.env bin/api db up

db_down:
	ENV_FILES=.env bin/api db down

db_create:
	ENV_FILES=.env bin/api db create '$(name)'

db_seed:
	ENV_FILES=.env bin/api db seed
