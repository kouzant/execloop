BIN_DIR="bin"

.PHONY: all
all: test build

.PHONY: test
test:
	@go test ./...

.PHONY: test-no-cache
test-no-cache:
	@go test -count=1 ./...
	
.PHONY: build
build:
	@echo "Building binary in ${BIN_DIR}"
	@go vet -c=3 ./...
	@GOOS=linux GOARCH=amd64 go build -v -o ${BIN_DIR}/execloop

.PHONY: clean
clean:
	@echo "Deleting ${BIN_DIR}"
	@rm -rf bin/
	@go clean -i -x -cache -testcache 