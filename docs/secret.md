<!-- BEGIN_TF_DOCS -->

## Requirements

No requirements.
## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
## Inputs

No inputs.
## Modules

No modules.
## Resources

| Name | Type |
|------|------|
| [aws_secretsmanager_secret.sign_key](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/secretsmanager_secret) | data source |
| [aws_secretsmanager_secret_version.sign_key_val](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/secretsmanager_secret_version) | data source |
## Outputs

| Name | Description |
|------|-------------|
| <a name="output_sign_key"></a> [sign\_key](#output\_sign\_key) | The generated value for the sign key |
<!-- END_TF_DOCS -->