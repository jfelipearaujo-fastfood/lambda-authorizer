output "sign_key" {
  description = "The generated value for the sign key"
  sensitive   = true
  value       = data.aws_secretsmanager_secret_version.sign_key_val.secret_string
}
