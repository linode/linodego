.PHONY: test

test:
	@LINODE_FIXTURE_MODE="play" \
	LINODE_TOKEN="awesometokenawesometokenawesometoken" \
	LINODE_API_VERSION="v4beta" \
	GO111MODULE="on" \
	go test -v ./integration $(ARGS)
