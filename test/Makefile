.PHONY: testint

testint:
	@LINODE_FIXTURE_MODE="play" \
	LINODE_TOKEN="awesometokenawesometokenawesometoken" \
	LINODE_API_VERSION="v4beta" \
	GO111MODULE="on" \
	go test -v ./integration $(ARGS)

.PHONY: testunit

testunit:
	go test -v ./unit/...

.PHONY: smoketest

smoketest:
	@LINODE_FIXTURE_MODE="record" \
	LINODE_TOKEN=$(LINODE_TOKEN) \
	LINODE_API_VERSION="v4beta" \
	LINODE_URL="$(LINODE_URL)" \
	GO111MODULE="on" \
	go test -v -run smoke ./integration/...


.PHONY: unit-test
unit-test:
	go test -v ./unit $(ARGS)