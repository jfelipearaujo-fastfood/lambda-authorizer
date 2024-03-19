module "database" {
  source = "./modules/database"

  db_name = "fastfood"
}

module "secret" {
  source = "./modules/secret"
}

module "authorizer" {
  source = "./modules/authorizer"

  lambda_name = "authorizer"

  sign_key = module.secret.sign_key

  db_port     = module.database.db_port
  db_name     = module.database.db_name
  db_username = module.database.db_username

  private_subnets   = var.private_subnets
  security_group_id = module.database.security_group_id

  depends_on = [
    module.secret
  ]
}
