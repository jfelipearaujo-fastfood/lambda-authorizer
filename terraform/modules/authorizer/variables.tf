variable "lambda_name" {
  type        = string
  description = "The name of the lambda function"
}

variable "sign_key" {
  type        = string
  sensitive   = true
  description = "The sign key for the lambda function"
}

variable "db_port" {
  type        = number
  description = "The port of the database"
}

variable "db_name" {
  type        = string
  description = "The name of the database"
  default     = "fastfood"
}

variable "db_username" {
  type        = string
  description = "The username for the database"
}

variable "security_group_id" {
  type        = string
  description = "The ID of the security group"
}

variable "private_subnets" {
  type        = list(string)
  description = "The IDs of the private subnets"
}
