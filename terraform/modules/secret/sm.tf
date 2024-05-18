data "aws_secretsmanager_secret" "sign_key" {
  name = "lambda_sign_key"
}

data "aws_secretsmanager_secret_version" "sign_key_val" {
  secret_id = data.aws_secretsmanager_secret.sign_key.id
}
