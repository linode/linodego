.PHONY: testint

testint:
	@LINODE_FIXTURE_MODE="play" \
	LINODE_TOKEN="9f08d4534157252421c964c7b9bf53dd545a242c1f7ef4fe1ada1701366acf7f" \
	LINODE_API_VERSION="v4beta" \
	GO111MODULE="on" \
	go test -v ./integration $(ARGS)

.PHONY: testunit

testunit:
	go test -v ./unit/...

.PHONY: smoketest

smoketest:
	@LINODE_FIXTURE_MODE="play" \
	LINODE_TOKEN="awesometokenawesometokenawesometoken" \
	LINODE_API_VERSION="v4beta" \
	GO111MODULE="on" \
	go test -v -run smoke ./integration/...


.PHONY: unit-test
unit-test:
	go test -v ./unit $(ARGS)