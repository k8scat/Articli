NAME := articli
CGO_ENABLED = 0
BUILD_GOOS = $(shell go env GOOS)
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
BIN_PATH := $(shell rm -f articli && which articli || echo "/usr/local/bin/articli")
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1) || echo "1.0.0")
BUILD_FLAGS = -ldflags "-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(shell date +'%Y-%m-%d')"
MAIN_SRC_FILE = cmd/articli/main.go

.PHONY: build

build: pre-build
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=$(BUILD_GOOS) GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/$(BUILD_GOOS)/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/$(BUILD_GOOS)/$(NAME)
	rm -rf $(NAME) && ln -s bin/$(BUILD_GOOS)/$(NAME) $(NAME)

darwin: pre-build
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)
	rm -rf $(NAME) && ln -s bin/darwin/$(NAME) $(NAME)

linux: pre-build
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)
	rm -rf $(NAME)
	ln -s bin/linux/$(NAME) $(NAME)

win: pre-build
	go get github.com/inconshreveable/mousetrap
	go get github.com/mattn/go-isatty
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 $(GO) $(BUILD_TARGET) $(BUILD_FLAGS) -o bin/windows/$(NAME).exe $(MAIN_SRC_FILE)

build-all: darwin linux win

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

clean: ## Clean the generated articlifacts
	rm -rf bin release
	rm -rf coverage.out
	rm -rf app/cmd/test-app.xml
	rm -rf app/test-app.xml
	rm -rf util/test-utils.xml

copy: build
	sudo cp bin/$(BUILD_GOOS)/$(NAME) $(BIN_PATH)

get-golint:
	go get -u golang.org/x/lint/golint

tools: get-golint

verify: dep tools lint

pre-build: fmt vet
	export GO111MODULE=on
	export GOPROXY=https://goproxy.io,direct
	go mod tidy

vet:
	go vet ./...

lint: vet
	golint -set_exit_status cmd/...
	golint -set_exit_status internal/...
	golint -set_exit_status pkg/...

fmt:
	go fmt ./cmd/...
	go fmt ./internal/...
	go fmt ./pkg/...
	gofmt -s -w .

test:
	mkdir -p bin
	go test ./pkg ./internal ./cmd/ -v -count=1 -coverprofile coverage.out
	go test ./app/cmd -v -count=1
#	go test ./util -v -count=1
#	go test ./client -v -count=1 -coverprofile coverage.out
#	go test ./app -v -count=1
#	go test ./app/health -v -count=1
#	go test ./app/helper -v -count=1
#	go test ./app/i18n -v -count=1
#	go test ./app/cmd -v -count=1

test-release:
	goreleaser release --rm-dist --snapshot --skip-publish

dep:
	go get github.com/AlecAivazis/survey/v2
	go get github.com/spf13/cobra
	go get github.com/spf13/viper
	go get gopkg.in/yaml.v2
	go get github.com/Pallinder/go-randomdata
	go install github.com/gosuri/uiprogress

image:
	docker build . -t k8scat/$(NAME)
