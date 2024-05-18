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
| <a name="input_lambda_name"></a> [lambda\_name](#input\_lambda\_name) | The name of the lambda function | `string` | n/a | yes |
| <a name="input_private_subnets"></a> [private\_subnets](#input\_private\_subnets) | The IDs of the private subnets | `list(string)` | n/a | yes |
| <a name="input_security_group_id"></a> [security\_group\_id](#input\_security\_group\_id) | The ID of the security group | `string` | n/a | yes |
| <a name="input_sign_key"></a> [sign\_key](#input\_sign\_key) | The sign key for the lambda function | `string` | n/a | yes |
## Modules

No modules.
## Resources

| Name | Type |
|------|------|
| [aws_iam_policy_attachment.iam_role_policy_attachment_vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy_attachment) | resource |
| [aws_iam_policy_attachment.lambda_policy_attachment](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy_attachment) | resource |
| [aws_iam_role.lambda_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_lambda_function.lambda_function](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
## Outputs

No outputs.
<!-- END_TF_DOCS -->