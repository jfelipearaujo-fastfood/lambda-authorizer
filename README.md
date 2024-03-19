# Fast-Food Lambda Authorizer

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
| <a name="input_private_subnets"></a> [private\_subnets](#input\_private\_subnets) | The IDs of the private subnets | `list(string)` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The default region to use for AWS | `string` | `"us-east-1"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | The default tags to use for AWS resources | `map(string)` | <pre>{<br>  "App": "lambda-authorizer"<br>}</pre> | no |
| <a name="input_vpc_cidr_block"></a> [vpc\_cidr\_block](#input\_vpc\_cidr\_block) | The CIDR block of the VPC | `string` | `"10.0.0.0/16"` | no |
| <a name="input_vpc_id"></a> [vpc\_id](#input\_vpc\_id) | The ID of the VPC | `string` | n/a | yes |
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

## Modules

- [database](./docs/database.md)
- [authorizer](./docs/authorizer.md)
- [secret](./docs/secret.md)