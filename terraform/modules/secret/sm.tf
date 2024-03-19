data "aws_secretsmanager_secret" "secret" {
  name = "lambda_sign_key"
}

data "aws_secretsmanager_secret_version" "secret_val" {
  secret_id = data.aws_secretsmanager_secret.secret.id
}
