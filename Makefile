GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet


all: lint build test

lint:
	golangci-lint run ./...
 
build:
	go build -v -race ./...

test:
	go test -v -race ./...

cover:
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -html=profile.cov
 
clean:
	go clean
	rm --force ./profile.cov




