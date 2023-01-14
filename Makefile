#!make

run:
	find . -name '*.go' 2>&1 | entr -r bash -c "go run ./cmd/api 2>&1"
