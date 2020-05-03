PROJECT_NAME := "dofus-bot"

.PHONY: all generate dep build clean

all: build

generate: ## Execute `go generate` command to execute commands inside Go code
	@go generate ./...

dep: ## Get the dependencies
	@go mod tidy
	@go mod download

build: dep generate ## Build the binary file
	@go build -v -o $(PROJECT_NAME) .

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)
