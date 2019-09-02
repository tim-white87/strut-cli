build: ./cmd/main.go
	go build -o strut ./cmd/main.go

.PHONY: run
run: ./cmd/main.go
	go run ./cmd/main.go

.PHONY: install
install: ./cmd/main.go
	go build -o $(GOPATH)/bin/strut ./cmd/main.go