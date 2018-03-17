include .env

.PHONY: vendor example

$(GOPATH)/bin/dep:
	@go get -u github.com/golang/dep/cmd/dep

vendor: $(GOPATH)/bin/dep
	@dep ensure

example:
	@go run example/main.go

test/fixtures.yaml:
	@LINODE_API_KEY=$(LINODE_API_KEY) go run test/main.go
	@sed -i "s/$(LINODE_API_KEY)/awesometokenawesometokenawesometoken/" $@

test: vendor test/fixtures.yaml
	@go test
