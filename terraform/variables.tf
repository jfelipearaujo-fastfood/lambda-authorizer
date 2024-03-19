variable "tags" {
  type        = map(string)
  description = "The default tags to use for AWS resources"
  default = {
    App = "lambda-authorizer"
  }
}

variable "region" {
  type        = string
  description = "The default region to use for AWS"
  default     = "us-east-1"
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC"
}

variable "vpc_cidr_block" {
  type        = string
  description = "The CIDR block of the VPC"
  default     = "10.0.0.0/16"
}

variable "private_subnets" {
  type        = list(string)
  description = "The IDs of the private subnets"
}
