build: ./cmd/main.go
	go build -o strut ./cmd/main.go

.PHONY: run
run: ./cmd/main.go
	go run ./cmd/main.go

install: ./cmd/main.go
	go build -o $(GOPATH)/bin/strut ./cmd/main.go

.PHONY: test
test: ./cmd/main.go
	go test ./cmd/main.go