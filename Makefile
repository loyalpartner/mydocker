# GO_BIN = /usr/local/go/bin
GO_BIN = $(shell which go)

.PHONY: run
run:
	sudo su -c "$(GO_BIN) run cmd/mydocker.go"

