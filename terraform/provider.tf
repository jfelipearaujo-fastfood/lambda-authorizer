terraform {
  required_version = "1.7.4"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.38.0"
    }
  }

  backend "s3" {
    region = "us-east-1"
    key    = "terraform/lambda/authorizer/terraform.tfstate"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}
