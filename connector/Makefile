NAME    := connector
VERSION := 0.0.1

BIN_DIR := ./bin

.PHONY: build

build:
	mkdir -p ${BIN_DIR}
	GOOS=linux go build -o ${BIN_DIR}/${NAME}-${VERSION} ./cmd/main.go

run: build
	${BIN_DIR}/${NAME}-${VERSION} ${flags}

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.37.1
	golangci-lint --version

lint:
	golangci-lint run ./...

test:
	go test -cover ./...  -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean
	rm -rf ${BIN_DIR}
