build: clean
	go build ./cmd/main.go

.PHONY: run
run: clean
	go run ./cmd/main.go

install: ./cmd/main.go
	go build -o $(GOPATH)/bin/strut ./cmd/main.go

.PHONY: test
test: clean
	go test -v -coverprofile cover.out ./...

.PHONY: clean
clean:
	rm -f main
	rm -f ./test/testdata/**
	rm -f cover.*

coverage:
	go tool cover -html=cover.out -o cover.html
	open cover.html