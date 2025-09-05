BINARY=wgetgo
BIN_DIR=bin

.PHONY: build lint clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o ${BIN_DIR}/${BINARY} ./cmd/wget

lint:
	go vet ./...
	golangci-lint run ./...

format:
	goimports -local github.com/aliskhannn/wget-go -w .

clean:
	rm -rf ${BIN_DIR}
	rm -rf sites
	rm -f wgetgo