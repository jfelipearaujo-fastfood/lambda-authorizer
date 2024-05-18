output "security_group_id" {
  description = "The ID of the security group"
  value       = data.aws_security_group.db_security_group.id
}
