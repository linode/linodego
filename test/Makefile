.PHONY: test-int

test-int:
	@LINODE_FIXTURE_MODE="play" \
	LINODE_TOKEN="awesometokenawesometokenawesometoken" \
	LINODE_API_VERSION="v4beta" \
	GO111MODULE="on" \
	go test -v ./integration $(TEST_ARGS) -timeout=$(TEST_TIMEOUT)

.PHONY: test-unit

test-unit:
	go test -v ./unit/... $(TEST_ARGS)

.PHONY: test-smoke

test-smoke:
	@LINODE_FIXTURE_MODE="record" \
	LINODE_TOKEN=$(LINODE_TOKEN) \
	LINODE_API_VERSION="v4beta" \
	LINODE_URL="$(LINODE_URL)" \
	GO111MODULE="on" \
	go test -v -run smoke ./integration/... -timeout=$(TEST_TIMEOUT)
