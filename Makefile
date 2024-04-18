build:
	@echo "Building..."
	@env GOOS=linux GOARCH=arm64 go build -o terraform/bootstrap main.go

zip:
	@echo "Zipping..."
	@zip terraform/lambda.zip terraform/bootstrap

upload:
	@echo "Creating folder..."
	@aws s3api put-object --bucket jfelipearaujo-org --key "migrations/" > /dev/null
	@echo "Uploading..."
	@aws s3 cp scripts/ s3://jfelipearaujo-org/migrations --recursive

init:
	@echo "Initializing..."
	@cd terraform \
		&& terraform init -reconfigure

check:
	@echo "Checking..."
	make fmt && make validate && make plan

plan:
	@echo "Planning..."
	@cd terraform \
		&& terraform plan -var-file="local.tfvars" -out=plan \
		&& terraform show -json plan > plan.tfgraph

fmt:
	@echo "Formatting..."
	@cd terraform \
		&& terraform fmt -check -recursive

validate:
	@echo "Validating..."
	@cd terraform \
		&& terraform validate

apply:
	@echo "Applying..."
	@cd terraform \
		&& terraform apply plan

destroy:
	@echo "Destroying..."
	@cd terraform \
		&& terraform destroy -auto-approve

gen-tf-docs:
	@echo "Generating Terraform Docs..."
	@terraform-docs markdown table terraform

gen-scaffold-bdd:
	@echo "Generating BDD scaffold..."
	@godog ./tests/features

test:
	@echo "Running tests..."
	@go test -count=1 ./src/... -v

test-bdd:
	@echo "Running BDD tests..."
	@go test -count=1 ./tests/... -test.v -test.run ^TestFeatures$