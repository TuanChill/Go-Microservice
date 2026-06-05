variable "aws_region" {
  type    = string
  default = "ap-southeast-1"
}

variable "key_name" {
  type        = string
  description = "EC2 key pair name for SSH access"
}

variable "allowed_ssh_cidr" {
  type        = string
  description = "CIDR allowed to SSH — restrict to your IP or bastion"
}

variable "app_port" {
  type    = number
  default = 8080
}

variable "app_instance_type" {
  type    = string
  default = "t3.small"
}

variable "db_name" {
  type    = string
  default = "go_template"
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "mq_username" {
  type = string
}

variable "mq_password" {
  type      = string
  sensitive = true
}

variable "ses_sender_email" {
  type        = string
  description = "Verified SES sender email address"
}

variable "s3_bucket_name" {
  type        = string
  description = "S3 bucket name — must be globally unique across all AWS accounts"
}
