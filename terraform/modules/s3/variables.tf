variable "env" {
  type = string
}

variable "bucket_name" {
  type        = string
  description = "S3 bucket name — must be globally unique"
}

variable "versioning_enabled" {
  type    = bool
  default = false
}

variable "force_destroy" {
  type        = bool
  default     = false
  description = "Allow Terraform to delete bucket even when non-empty (use true for UAT only)"
}
