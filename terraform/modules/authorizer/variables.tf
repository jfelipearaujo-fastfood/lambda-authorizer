variable "lambda_name" {
  type        = string
  description = "The name of the lambda function"
}

variable "sign_key" {
  type        = string
  sensitive   = true
  description = "The sign key for the lambda function"
}

variable "security_group_id" {
  type        = string
  description = "The ID of the security group"
}

variable "private_subnets" {
  type        = list(string)
  description = "The IDs of the private subnets"
}
