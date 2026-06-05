variable "env" {
  type = string
}

variable "sender_email" {
  type        = string
  description = "Email address to verify in SES as the sender (must click verification link)"
}

variable "s3_bucket_arn" {
  type        = string
  default     = ""
  description = "S3 bucket ARN to grant app read/write access — leave empty to skip"
}
