build: ./cmd/main.go
	go build ./cmd/main.go

.PHONY: run
run: ./cmd/main.go
	go run ./cmd/main.go

install: ./cmd/main.go
	go build -o $(GOPATH)/bin/strut ./cmd/main.go

.PHONY: test
test:
	go test ./...