# type 'make help' to visualize the help information

##@ Utility
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ CI/CD
build: ## Build the application to the output folder (default: ./buil/main)
	@echo "Building..."	
	@go build -ldflags="-s -w" -o build/main cmd/main.go

clean: ## Clean the binary
	@echo "Cleaning..."
	@rm -f build/main

sec: ## Security checker
	@if command -v gosec > /dev/null; then \
		echo "Analyzing..."; \
		gosec ./...; \
	else \
		read -p "gosec is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/securego/gosec/v2/cmd/gosec@latest; \
			echo "Analyzing..."; \
			gosec ./...; \
		else \
			echo "You chose not to intall gosec. Exiting..."; \
			exit 1; \
		fi; \
	fi

lint: ## Go Linter
	@if command -v golangci-lint > /dev/null; then \
		echo "Analyzing..."; \
		golangci-lint run ./...; \
	else \
		read -p "golangci-lint is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2; \
			echo "Analyzing..."; \
			golangci-lint run ./...; \
		else \
			echo "You chose not to intall golangci-lint. Exiting..."; \
			exit 1; \
		fi; \
	fi

##@ Runner
run: build docker-up ## Run the application
	@if test ! -f .env; then \
		make env; \
	fi; \
	go run -race -ldflags="-s -w" cmd/main.go local

##@ Testing
test-queue: ## Test the queue
	@echo "Testing..."
	@./scripts/tests/send-message.sh

test: ## Test the application
	@if command -v gcc > /dev/null; then \
		echo "Testing..."; \
		go test -race -ldflags="-s -w" -count=1 ./internal/... -coverprofile=coverage.out; \
	else \
		read -p "gcc is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			sudo apt install build-essential; \
			echo "Testing..."; \
			go test -race -ldflags="-s -w" -count=1 ./internal/... -coverprofile=coverage.out; \
		else \
			echo "You chose not to intall gcc. Exiting..."; \
			exit 1; \
		fi; \
	fi

test-bdd: ## Run BDD tests
	@echo "Running BDD tests..."
	@go test -race -count=1 ./tests/... -test.v -test.run ^TestFeatures$

cover: ## View the coverage
	@echo "Analyzing coverage..."
	@go tool cover -func=coverage.out

test-cover: ## Run the tests and view the coverage
	make test && make cover

bench: ## Run benchmarks
	@echo "Benchmarking..."
	@go test -cpu=1,2,4,6,8,16 -benchmem -bench=. ./internal/features/transactions/handlers/create_transaction -run=^Benchmark

load: ## Run the load test using k6
	@if command -v k6 > /dev/null; then \
		echo "Generating..."; \
		k6 run tests/index.js; \
	else \
		read -p "k6 is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			brew install k6; \
			echo "Generating..."; \
			k6 run tests/index.js; \
		else \
			echo "You chose not to intall k6. Exiting..."; \
			exit 1; \
		fi; \
	fi

##@ Developing
env: ## Create the .env file based on example
	@echo "Generating..."
	@cp .env.example .env

docker-up: ## Run the containers
	@if command -v docker compose > /dev/null; then \
		docker compose up -d; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi

docker-down: ## Shutdown containers
	@if command -v docker compose > /dev/null; then \
		docker compose down; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

watch: env docker-up ## Live reload using air
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

TAG := $(shell git describe --tags --abbrev=0 2>/dev/null)
VERSION := $(shell echo $(TAG) | sed 's/v//')

tag: ## Create or bump the version tag
	@if [ -z "$(TAG)" ]; then \
        echo "No previous version found. Creating v1.0 tag..."; \
        git tag v1.0; \
    else \
        echo "Previous version found: $(VERSION)"; \
        read -p "Bump major version (M/m) or release version (R/r)? " choice; \
        if [ "$$choice" = "M" ] || [ "$$choice" = "m" ]; then \
            echo "Bumping major version..."; \
			major=$$(echo $(VERSION) | cut -d'.' -f1); \
            major=$$(expr $$major + 1); \
            new_version=$$major.0; \
        elif [ "$$choice" = "R" ] || [ "$$choice" = "r" ]; then \
            echo "Bumping release version..."; \
			release=$$(echo $(VERSION) | cut -d'.' -f2); \
            release=$$(expr $$release + 1); \
            new_version=$$(echo $(VERSION) | cut -d'.' -f1).$$release; \
        else \
            echo "Invalid choice. Aborting."; \
            exit 1; \
        fi; \
        echo "Creating tag for version v$$new_version..."; \
        git tag v$$new_version; \
    fi

##@ Auto generated files
gen-mocks: ## Gen mock files using mockery
	@if command -v mockery > /dev/null; then \
		echo "Generating..."; \
		mockery; \
	else \
		read -p "Go 'mockery' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/vektra/mockery/v2@latest; \
			echo "Generating..."; \
			mockery; \
		else \
			echo "You chose not to intall mockery. Exiting..."; \
			exit 1; \
		fi; \
	fi

gen-docs: ## Gen Swagger docs using swag
	@if command -v swag > /dev/null; then \
		swag init --parseDependency -g main.go -d cmd/api,internal; \
	else \
		read -p "Go 'swag' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/swaggo/swag/cmd/swag@latest; \
			swag init --parseDependency -g main.go -d cmd/api,internal; \
		else \
			echo "You chose not to intall swag. Exiting..."; \
			exit 1; \
		fi; \
	fi

gen-pkg-docs: ## Gen Package docs using gomarkdoc
	@if command -v gomarkdoc > /dev/null; then \
		gomarkdoc --output '{{.Dir}}/README.md' ./...; \
	else \
		read -p "Go 'gomarkdoc' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest; \
			gomarkdoc --output '{{.Dir}}/README.md' ./...; \
		else \
			echo "You chose not to intall gomarkdoc. Exiting..."; \
			exit 1; \
		fi; \
	fi

fmt-docs: ## Format generated Swagger docs using swag
	@if command -v swag > /dev/null; then \
		swag fmt -d internal -g cmd/api/main.go; \
	else \
		read -p "Go 'swag' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/swaggo/swag/cmd/swag@latest; \
			swag fmt -d internal -g cmd/api/main.go; \
		else \
			echo "You chose not to intall swag. Exiting..."; \
			exit 1; \
		fi; \
	fi

gen-scaffold-bdd: ## Gen BDD scaffold using godog
	@if command -v godog > /dev/null; then \
		echo "Generating BDD scaffold..."; \
		godog ./tests/features; \
	else \
		read -p "Go 'godog' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cucumber/godog/cmd/godog@latest; \
			echo "Generating BDD scaffold..."; \
			godog ./tests/features; \
		else \
			echo "You chose not to intall godog. Exiting..."; \
			exit 1; \
		fi; \
	fi

gen-tf-docs: ## Generate Terraform Docs
	@echo "Generating Terraform Docs..."
	@terraform-docs markdown table terraform

gen-cloud-diagrams: ## Generate Cloud Diagrams
	@echo "Generating Cloud Diagrams..."
	@cd docs && python3 cloud_aws_database_migrations.py

##@ Terraform
init:
	@echo "Initializing..."
	@cd terraform \
		&& terraform init -reconfigure

check:
	@echo "Checking..."
	make fmt && make validate && make plan

fmt:
	@echo "Formatting..."
	@cd terraform \
		&& terraform fmt -check

validate:
	@echo "Validating..."
	@cd terraform \
		&& terraform validate

plan:
	@echo "Planning..."
	@cd terraform \
		&& terraform plan -var-file="local.tfvars" -out=plan \
		&& terraform show -json plan > plan.tfgraph

apply:
	@echo "Applying..."
	@cd terraform \
		&& terraform apply plan

destroy:
	@echo "Destroying..."
	@cd terraform \
		&& terraform destroy -auto-approve

build-binary:
	@echo "Building..."
	@env GOOS=linux GOARCH=arm64 go build -o terraform/bootstrap cmd/main.go

zip-binary:
	@echo "Zipping..."
	@zip terraform/lambda.zip terraform/bootstrap

.PHONY: build run test clean
