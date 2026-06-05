variable "env" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "security_group_id" {
  type = string
}

variable "mq_username" {
  type = string
}

variable "mq_password" {
  type      = string
  sensitive = true
}

variable "instance_type" {
  type    = string
  default = "mq.m5.large"
}

variable "deployment_mode" {
  type    = string
  default = "SINGLE_INSTANCE"
}
