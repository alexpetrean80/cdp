# lint golang using golangci-lint
#name of the executable
APP_NAME=cdp

deps:
	@echo "installing dependencies"
	@go mod tidy

# install golangci-lint
install_tools:
	@echo "installing golangci-lint"
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

lint:
	@echo "linting golang code"
	@golangci-lint run ./...

# build executable from cmd/main.go
build:
	@echo "building executable"
	@go build -o bin/$(APP_NAME) cmd/main.go

test:
	@echo "running tests"
	@go test -v ./...

cleanup:
	@echo "cleaning up"
	@rm -rf ./bin ./target

install:
	@make build
	@sudo cp ./bin/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	@make cleanup
