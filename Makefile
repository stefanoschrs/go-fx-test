#!make

run:
	find . -name '*.go' 2>&1 | entr -r bash -c "go run ./cmd/api 2>&1"

run-nofx:
	cd nofx && find . -name '*.go' 2>&1 | entr -r bash -c "go run ./cmd/api 2>&1"

build:
	go build -trimpath  -ldflags "-s -w" -o ./dist/api ./cmd/api

build-nofx:
	cd nofx && go build -trimpath  -ldflags "-s -w" -o ../dist/api-nofx ./cmd/api
