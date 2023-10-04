# GO_BIN = /usr/local/go/bin
GO_BIN = $(shell which go)

.PHONY: docker
docker:
	sudo su -c "$(GO_BIN) run cmd/mydocker/mydocker.go"

.PHONY: cgroup
cgroup:
	sudo su -c "$(GO_BIN) run cmd/cgroup/cgroup.go"


