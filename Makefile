BINARY=wgetgo
BIN_DIR=bin
TEST_SCRIPT=./integration/test_e2e.sh

.PHONY: build test integration lint clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o ${BIN_DIR}/${BINARY} ./cmd/wget

lint:
	go vet ./...
	golangci-lint run ./...

clean:
	rm -rf ${BIN_DIR}
	rm -rf sites
	rm -f wgetgo