# GO_BIN = /usr/local/go/bin
GO_BIN = $(shell which go)

.PHONY: namespace
namespace:
	sudo su -c "$(GO_BIN) run ./cmd/namespace/namespace.go"

.PHONY: cgroup
cgroup:
	sudo su -c "$(GO_BIN) run ./cmd/cgroup/cgroup.go"

.PHONY: docker
docker:
	sudo su -c "$(GO_BIN) run ./cmd/mydocker/ run -ti /bin/ls"


