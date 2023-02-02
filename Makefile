# Force Go Modules
GO111MODULE = on

GOCC ?= go
GOFLAGS ?=

# make reproducable
GOFLAGS += -trimpath

# If set, override the install location for plugins
IPFS_PATH ?= $(HOME)/.ipfs

# If set, override the IPFS version to build against. This _modifies_ the local
# go.mod/go.sum files and permanently sets this version.
IPFS_VERSION ?= $(lastword $(shell $(GOCC) list -m github.com/ipfs/kubo))

.PHONY: install build

# go.mod: FORCE
# 	./set-target.sh $(IPFS_VERSION)

# FORCE:

azureplugin.so: plugin/main/main.go go.mod
	$(GOCC) build $(GOFLAGS) -buildmode=plugin -o "$@" "$<"
	chmod +x "$@"

build: azureplugin.so
	@echo "Built against" $(IPFS_VERSION)

install: build
	install -Dm700 azureplugin.so "$(IPFS_PATH)/plugins/go-ds-azure.so"