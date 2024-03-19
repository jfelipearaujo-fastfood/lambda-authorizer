data "aws_security_group" "db_security_group" {
  name = "db-sg-${var.db_name}"
}
