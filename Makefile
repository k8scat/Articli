NAME := acli
CGO_ENABLED = 0
BUILD_GOOS = $(shell go env GOOS)
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
BIN_PATH := $(shell rm -f acli && which acli || echo "/usr/local/bin/acli")
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1) || echo "1.0.0")
BUILD_DATE := $(shell date +'%Y-%m-%d')
BUILD_FLAGS = -ldflags "-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(BUILD_DATE) -w -s" \
	-trimpath
MAIN_SRC_FILE = cmd/articli/main.go

.PHONY: build
build: pre-build
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=$(BUILD_GOOS) GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/$(BUILD_GOOS)/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/$(BUILD_GOOS)/$(NAME)
	rm -rf $(NAME) && ln -s bin/$(BUILD_GOOS)/$(NAME) $(NAME)

.PHONY: darwin
darwin: pre-build
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)
	rm -rf $(NAME) && ln -s bin/darwin/$(NAME) $(NAME)

.PHONY: linux
linux: pre-build
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)
	rm -rf $(NAME)
	ln -s bin/linux/$(NAME) $(NAME)

.PHONY: win
win: pre-build
	go get github.com/inconshreveable/mousetrap
	go get github.com/mattn/go-isatty
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/windows/$(NAME).exe $(MAIN_SRC_FILE)

.PHONY: build-all
build-all: darwin linux win

.PHONY: release
release: build-all
	mkdir -p release
	cd ./bin/darwin; upx $(NAME); \
		tar -zcvf ../../release/$(NAME)-darwin-amd64.tar.gz $(NAME); \
		cd ../../release/; \
		shasum -a 256 $(NAME)-darwin-amd64.tar.gz > $(NAME)-darwin-amd64.txt
	cd ./bin/linux; upx $(NAME); \
		tar -zcvf ../../release/$(NAME)-linux-amd64.tar.gz $(NAME); \
		cd ../../release/; \
		shasum -a 256 $(NAME)-linux-amd64.tar.gz > $(NAME)-linux-amd64.txt
	cd ./bin/windows; \
		upx $(NAME).exe; \
		tar -zcvf ../../release/$(NAME)-windows-386.tar.gz $(NAME).exe; \
		cd ../../release/; \
		shasum -a 256 $(NAME)-windows-386.tar.gz > $(NAME)-windows-386.txt

.PHONY: clean
clean: ## Clean the generated articlifacts
	rm -rf bin release
	rm -rf coverage.out

.PHONY: copy
copy: build
	sudo cp bin/$(BUILD_GOOS)/$(NAME) $(BIN_PATH)

.PHONY: get-golint
get-golint:
	go get -u golang.org/x/lint/golint

.PHONY: tools
tools: get-golint
	brew install goreleaser/tap/goreleaser
	brew install --build-from-source upx

.PHONY: verify
verify: dep tools lint

.PHONY: pre-build
pre-build: fmt vet
	export GO111MODULE=on
	# export GOPROXY=https://goproxy.io,direct
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: vet
	golint -set_exit_status cmd/...
	golint -set_exit_status internal/...
	golint -set_exit_status pkg/...

.PHONY: fmt
fmt:
	go fmt ./cmd/...
	go fmt ./internal/...
	go fmt ./pkg/...
	gofmt -s -w .

.PHONY: test
test:
	go test ./pkg/platform/juejin -v -count=1 -coverprofile coverage.out

.PHONY: test-release
test-release:
	goreleaser release --rm-dist --snapshot --skip-publish

.PHONY: dep
dep:
	go mod download

.PHONY: build-image
build-image:
	docker build -t k8scat/articli \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		.

TAG =
.PHONY: delete-tag
delete-tag:
ifneq ($(TAG),)
	git push --delete origin $(TAG)
	git tag --delete $(TAG)
else
	@echo "Usage: make delete-tag TAG=<tag>"
endif

.PHONY: npm-publish
npm-publish:
	npm publish --access public
