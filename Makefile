BINARY_NAME=ably-test
MAIN_DIR=ablytest/cmd/ablytest

all: dep build

dep:
	go mod download

build:
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-amd64-linux ${MAIN_DIR}
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-amd64-darwin ${MAIN_DIR}
	GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}-arm64-darwin ${MAIN_DIR}

clean:
	go clean
	rm -f bin/${BINARY_NAME}*
