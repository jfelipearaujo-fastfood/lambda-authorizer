# Lambda Authorizer

This project its responsible to validate the token and check if the user has the right permissions to access the API Gateway.

<!-- BEGIN_TF_DOCS -->

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | 1.7.4 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | 5.38.0 |
## Providers

No providers.
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_region"></a> [region](#input\_region) | The default region to use for AWS | `string` | `"us-east-1"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | The default tags to use for AWS resources | `map(string)` | <pre>{<br>  "App": "lambda-authorizer"<br>}</pre> | no |
| <a name="input_vpc_name"></a> [vpc\_name](#input\_vpc\_name) | The name of the VPC | `string` | `"vpc-fastfood"` | no |
## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_authorizer"></a> [authorizer](#module\_authorizer) | ./modules/authorizer | n/a |
| <a name="module_database"></a> [database](#module\_database) | ./modules/database | n/a |
| <a name="module_secret"></a> [secret](#module\_secret) | ./modules/secret | n/a |
## Resources

No resources.
## Outputs

No outputs.
<!-- END_TF_DOCS -->
