.PHONY: test
test: dep godog
	go test ./...
	(cd cmd/ && godog ../features)

.PHONY: lint
lint: gometalinter
	gometalinter.v2 --errors --vendor --exclude=vendor --config gometalinter.json ./...

.PHONY: mock
dir ?= .
mock: export filename=$(shell echo $(name) | tr A-Z a-z)_mock.go
mock: mockery
ifndef name
	@echo "Please specify an interface name: $ make mock name=MyInterface"
	exit 1
endif
	@echo "Generating mock for $(dir)/$(name)..."
	@cd $(dir) && mockery -inpkg -print -name $(name) >_$(filename)
	@cd $(dir) && mv _$(filename) $(filename)

.PHONY: mockery
MOCKERY_BIN := $(shell command -v mockery 2> /dev/null)
mockery:
ifndef MOCKERY_BIN
	@echo "Installing mockery..."
	@go get github.com/vektra/mockery/cmd/mockery
endif

.PHONY: godog
GODOG_BIN := $(shell command -v godog 2> /dev/null)
godog:
ifndef GODOG_BIN
	@echo "Installing godog..."
	@go get github.com/DATA-DOG/godog/cmd/godog
endif

.PHONY: dep
DEP_BIN := $(shell command -v dep 2> /dev/null)
dep:
ifndef DEP_BIN
	@echo "Installing dep..."
	@go get github.com/golang/dep/cmd/dep
	@dep ensure --vendor-only
endif

.PHONY: gomentalinter
GOMETALINTER_BIN := $(shell command -v gometalinter.v2 2> /dev/null)
gometalinter:
ifndef GOMETALINTER_BIN
	@echo "Installing gometalinter..."
	@go get gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install
endif
