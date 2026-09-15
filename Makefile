GOPATH ?= $(shell go env GOPATH)
NPM_DIR = $(shell npm config get prefix)

.PHONY: all
all: install

.PHONY: lint-go
lint-go: 
	./scripts/install-lint.sh
	${GOPATH}/bin/golangci-lint run

.PHONY: lint-shell
lint-shell:
	shellcheck --source-path=.:scripts scripts/*.sh scripts/*/*.sh

.PHONY: lint-solidity
lint-solidity: 
	"$(NPM_DIR)/bin/solhint" $$(find ethereum -name '*.sol' -not -path '**/@openzeppelin/**')

.PHONY: lint
lint: lint-go lint-shell lint-solidity

.PHONY: format-go
format-go:
	test -x $(GOPATH)/bin/goimports || go install golang.org/x/tools/cmd/goimports@latest
	$(GOPATH)/bin/goimports -local github.com/athanorlabs/atomic-swap -w .

.PHONY: format-shell
format-shell:
	shfmt -w scripts/*.sh scripts/*/*.sh

.PHONY: format-solidity
format-solidity:
	"$(NPM_DIR)/bin/prettier" --print-width 100 --write \
		--plugin="$(NPM_DIR)/lib/node_modules/prettier-plugin-solidity/src/index.js" \
		$$(find ethereum -name '*.sol' -not -path '**/@openzeppelin/**')

.PHONY: format
format: format-go format-shell format-solidity

.PHONY: test
test: 
	./scripts/run-unit-tests.sh 2>&1 | tee test.log

.PHONY: test-integration
test-integration: 
	./scripts/run-integration-tests.sh 2>&1 | tee test-integration.log

# Instead of building from the local checked-out source, this will install
# the most recent commit with a release tag. Use the most recent tagged release
# for production swaps.
.PHONY: build-release
build-release:
	mkdir -p bin
	GOBIN=$(PWD)/bin go install -tags=prod github.com/athanorlabs/atomic-swap/cmd/...@latest

# If you don't have go installed but do have docker, you can build the most
# recent release using docker.
.PHONY: build-release-in-docker
build-release-in-docker:
	mkdir -p bin
	docker run --rm -v "$(PWD)/bin:/go/bin" -v $(PWD)/Makefile:/go/Makefile "golang:1.20" bash -c \
		"make build-release && chown $$(id -u):$$(id -g) bin/{swapd,swapcli,bootnode}"

# Install all the binaries into $HOME/go/bin (or alternative GOPATH bin directory)
.PHONY: install
install: 
	go install ./cmd/...

# Install swapd and swapcli into this directory (top of the project)
.PHONY: build
build: 
	./scripts/build.sh

# Test macos/arm build from linux. Use "make build" if compiling on macos.
.PHONY: build-darwin
build-darwin:
	mkdir -p bin/
	GOOS=darwin GOARCH=arm64 go build  -o ./bin ./cmd/...

# Same as build, but also includes some lesser used binaries
.PHONY: build-all
build-all: 
	ALL=true $(MAKE) build

.PHONY: docker-images
docker-images:
	scripts/docker-swapd/build-docker-image.sh
	scripts/docker-bootnode/build-docker-image.sh

# Go bindings for solidity contracts
.PHONY: bindings
bindings:
	./scripts/install-abigen.sh
	./scripts/generate-bindings.sh
	./ethereum/block/testdata/generate-bindings.sh

.PHONY: mock
mock:
	test -x $(GOPATH)/bin/mockgen || go install github.com/golang/mock/mockgen@v1.6.0
	go generate -run mockgen ./...
	$(MAKE) format-go

# Deletes all executables matching the directory names in cmd/
.PHONY: clean
clean:
	rm -r bin/

# monero-swap: the seller/buyer program against XmrSwap (needs Go 1.21; newer Go breaks quic-go 0.39)
.PHONY: build-swap
build-swap:
	@go version | grep -q 'go1.21' || { echo "need Go 1.21 (see README); found: $$(go version)"; exit 1; }
	mkdir -p bin && go build -o bin/monero-swap ./cmd/monero-swap

# monero-swap-v2: the seller against XmrSwapV2 (the live contract). Builds on modern Go —
# it does not pull in Athanor's p2p stack, so the Go 1.21 quic-go constraint does not apply.
.PHONY: build-swap-v2
build-swap-v2:
	mkdir -p bin && go build -o bin/monero-swap-v2 ./cmd/monero-swap-v2

# Regenerate the Go bindings for XmrSwapV2 from the Foundry build output.
.PHONY: bindings-v2
bindings-v2:
	cd contracts-v2 && forge build
	jq -r '.abi' contracts-v2/out/XmrSwapV2.sol/XmrSwapV2.json > /tmp/XmrSwapV2.abi
	jq -r '.bytecode.object' contracts-v2/out/XmrSwapV2.sol/XmrSwapV2.json | sed 's/^0x//' > /tmp/XmrSwapV2.bin
	$(GOPATH)/bin/abigen --abi /tmp/XmrSwapV2.abi --bin /tmp/XmrSwapV2.bin \
		--pkg xmrswap2 --type XmrSwapV2 --out xmrswap2/xmrswap2.go
	rm -f /tmp/XmrSwapV2.abi /tmp/XmrSwapV2.bin
