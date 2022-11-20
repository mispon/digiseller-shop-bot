LOCAL_BIN=$(CURDIR)/bin

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: build
build:
	go build -o $(LOCAL_BIN)/bot $(CURDIR)/cmd