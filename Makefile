include .env

.PHONY: vendor example

$(GOPATH)/bin/dep:
	@go get -u github.com/golang/dep/cmd/dep

vendor: $(GOPATH)/bin/dep
	@dep ensure

example:
	@go run example/main.go

test/fixtures.yaml:
	@go run test/main.go

test: vendor test/fixtures.yaml
	@go test
