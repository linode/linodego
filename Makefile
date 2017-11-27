
.PHONY: vendor example

$(GOPATH)/bin/govendor:
	@go get -u github.com/kardianos/govendor

vendor: $(GOPATH)/bin/govendor
	@govendor sync

example:
	@go run example/main.go 