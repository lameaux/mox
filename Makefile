SRC_DIR := .
BUILD_DIR := ./bin
GIT_HASH := $(shell git rev-parse --short HEAD)
DOCKER_REPO := ghcr.io
DOCKER_IMAGE := lameaux/mox

GO_FILES := $(shell find $(SRC_DIR) -name '*.go' ! -path '$(SRC_DIR)/protos/*go')

.PHONY: all
all: clean build lint test

.PHONY: build
build:
	go build -ldflags "-X main.GitHash=$(GIT_HASH)" -o $(BUILD_DIR)/mox $(SRC_DIR)/cmd/mox/*.go

.PHONY: fmt
fmt:
	gci write $(GO_FILES) --skip-generated -s standard -s default
	gofumpt -l -w $(GO_FILES)

.PHONY: fmt-install
fmt-install:
	go install github.com/daixiang0/gci@latest
	go install mvdan.cc/gofumpt@latest

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-install
lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: test
test:
	go test $(SRC_DIR)/... -coverprofile=coverage.out

.PHONY: coverage
coverage:
	go tool cover -func coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 80) != 1) }'


.PHONY: install
install: build
	mv $(BUILD_DIR)/mox $(GOPATH)/bin

.PHONY: run
run: build
	$(BUILD_DIR)/mox $(ARGS)

.PHONY: serve
serve: run

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: docker-build
docker-build:
	docker build --build-arg GIT_HASH=$(GIT_HASH) -t $(DOCKER_IMAGE):$(GIT_HASH) .

.PHONY: docker-push
docker-push:
	docker tag $(DOCKER_IMAGE):$(GIT_HASH) ghcr.io/$(DOCKER_IMAGE):latest
	docker push ghcr.io/$(DOCKER_IMAGE):latest

.PHONY: docker-release
docker-release: docker-build docker-push
