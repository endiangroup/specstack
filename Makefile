.PHONY: install
install: vendor $(GOPATH)/bin/spec $(GOPATH)/bin/specfmt

$(GOPATH)/bin/spec:
	@go install github.com/endiangroup/specstack/cmd/spec
  
$(GOPATH)/bin/specfmt:
	@go install github.com/endiangroup/specstack/cmd/specfmt

.PHONY: clean
clean:
	@rm -rf vendor $(GOPATH)/bin/spec $(GOPATH)/bin/specfmt

.PHONY: vendor
vendor: dep
	@dep ensure --vendor-only

.PHONY: test
test: dep godog
	go test ./...
	(cd cmd/ && godog ../features)

.PHONY: lint
lint: golangci-lint $(GOPATH)/bin/specfmt
	golangci-lint run ./...
	specfmt -l features/*.feature

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
endif

.PHONY: golangci-lint
GOLANGCI_BIN := $(shell command -v golangci-lint 2> /dev/null)
golangci-lint:
ifndef GOLANGCI_BIN
	-@go get -u github.com/golangci/golangci-lint
	@cd $(GOPATH)/src/github.com/golangci/golangci-lint/cmd/golangci-lint
	@go install -ldflags "-X 'main.version=$(git describe --tags)' -X 'main.commit=$(git rev-parse --short HEAD)' -X 'main.date=$(date)'"
endif
