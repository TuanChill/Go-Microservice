variable "env" {
  type = string
}

variable "vpc_cidr" {
  type = string
}

variable "public_subnet_cidrs" {
  type = list(string)
}

variable "private_subnet_cidrs" {
  type    = list(string)
  default = []
}

variable "azs" {
  type = list(string)
}

variable "enable_nat_gateway" {
  type    = bool
  default = false
}
