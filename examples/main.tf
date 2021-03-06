data "aws_caller_identity" "current" {
  provider = "aws"
}
data "aws_caller_identity" "ident" {
  provider = "aws.ident"
}
terraform {
  required_providers {
    aws = ">= 2.15.0"
  }
}