module "database" {
  source = "./modules/database"

  db_name = "customers"
}

module "secret" {
  source = "./modules/secret"
}

module "authorizer" {
  source = "./modules/authorizer"

  lambda_name = "authorizer"

  sign_key = module.secret.sign_key

  private_subnets   = var.private_subnets
  security_group_id = module.database.security_group_id

  depends_on = [
    module.secret
  ]
}
