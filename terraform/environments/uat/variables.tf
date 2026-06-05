variable "aws_region" {
  type    = string
  default = "ap-southeast-1"
}

variable "instance_type" {
  type    = string
  default = "t3.small"
}

variable "ami_id" {
  type        = string
  default     = ""
  description = "AMI ID for the app EC2. If empty, uses SSM parameter for Ubuntu 22.04."
}

variable "key_name" {
  type        = string
  description = "EC2 key pair name for SSH access"
}

variable "allowed_ssh_cidr" {
  type        = string
  description = "CIDR allowed to SSH — restrict to your IP"
  default     = "0.0.0.0/0"
}

variable "app_port" {
  type    = number
  default = 8080
}

variable "db_name" {
  type    = string
  default = "go_template"
}

variable "db_user" {
  type    = string
  default = "postgres"
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "redis_password" {
  type      = string
  sensitive = true
  default   = ""
}

variable "rabbit_user" {
  type    = string
  default = "admin"
}

variable "rabbit_password" {
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
