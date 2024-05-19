<!-- BEGIN_TF_DOCS -->

## Requirements

No requirements.
## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_db_name"></a> [db\_name](#input\_db\_name) | The name of the database | `string` | n/a | yes |
## Modules

No modules.
## Resources

| Name | Type |
|------|------|
| [aws_db_instance.db](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/db_instance) | data source |
| [aws_security_group.db_security_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/security_group) | data source |
## Outputs

| Name | Description |
|------|-------------|
| <a name="output_security_group_id"></a> [security\_group\_id](#output\_security\_group\_id) | The ID of the security group |
<!-- END_TF_DOCS -->