# lint golang using golangci-lint
#name of the executable
APP_NAME=cdp

deps:
	@echo "installing dependencies"
	@go mod tidy

# install golangci-lint
install_tools:
	@echo "installing cobra-cli"
	@go install github.com/spf13/cobra-cli@latest

lint:
	@echo "linting golang code"
	@golangci-lint run ./...

# build executable from cmd/main.go
build:
	@echo "building executable"
	@go build -o bin/$(APP_NAME) main.go

run:
	@go run main.go

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
