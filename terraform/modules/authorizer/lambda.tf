data "aws_secretsmanager_secret" "master_user_secret" {
  name = "db-${var.db_name}-secret"
}

data "aws_secretsmanager_secret_version" "master_user_secret_version" {
  secret_id = data.aws_secretsmanager_secret.master_user_secret.arn
}

resource "aws_lambda_function" "lambda_function" {
  function_name = "lambda_${var.lambda_name}"

  filename      = "./lambda.zip"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  memory_size   = 128
  timeout       = 30

  environment {
    variables = {
      SIGN_KEY = var.sign_key
      DB_HOST  = jsondecode(data.aws_secretsmanager_secret_version.master_user_secret_version.secret_string)["host"]
      DB_PORT  = var.db_port
      DB_NAME  = var.db_name
      DB_USER  = var.db_username
      DB_PASS  = jsondecode(data.aws_secretsmanager_secret_version.master_user_secret_version.secret_string)["password"]
    }
  }

  source_code_hash = filebase64sha256("./lambda.zip")

  vpc_config {
    ipv6_allowed_for_dual_stack = false
    subnet_ids                  = var.private_subnets
    security_group_ids          = [var.security_group_id]
  }
}
