BIN := ck

export GO111MODULE=on

all: clean build

build:
	go build -o $(BIN) ./cmd/$(BIN)/

install:
	go install ./...

clean:
	rm -rf $(BIN)
	go clean

.PHONY: test
test:
	go test -v ./...

