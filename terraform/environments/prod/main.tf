provider "aws" {
  region = var.aws_region
}

# Uncomment after bootstrapping state bucket (same bucket as UAT, different key):
# terraform {
#   backend "s3" {
#     bucket         = "YOUR_BUCKET"
#     key            = "prod/terraform.tfstate"
#     region         = "ap-southeast-1"
#     dynamodb_table = "terraform-locks"
#     encrypt        = true
#   }
# }

locals {
  azs = ["${var.aws_region}a", "${var.aws_region}b"]
}

module "networking" {
  source               = "../../modules/networking"
  env                  = "prod"
  vpc_cidr             = "10.1.0.0/16"
  public_subnet_cidrs  = ["10.1.1.0/24", "10.1.2.0/24"]
  private_subnet_cidrs = ["10.1.11.0/24", "10.1.12.0/24"]
  azs                  = local.azs
  enable_nat_gateway   = true
}

module "security_groups" {
  source           = "../../modules/security_groups"
  env              = "prod"
  vpc_id           = module.networking.vpc_id
  allowed_ssh_cidr = var.allowed_ssh_cidr
  app_port         = var.app_port
}

module "rds" {
  source              = "../../modules/rds"
  env                 = "prod"
  subnet_ids          = module.networking.private_subnet_ids
  security_group_id   = module.security_groups.rds_sg_id
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  instance_class      = "db.t3.small"
  allocated_storage   = 20
  multi_az            = false
  deletion_protection = true
}

module "elasticache" {
  source            = "../../modules/elasticache"
  env               = "prod"
  subnet_ids        = module.networking.private_subnet_ids
  security_group_id = module.security_groups.cache_sg_id
  node_type         = "cache.t3.micro"
}

module "amazon_mq" {
  source            = "../../modules/amazon_mq"
  env               = "prod"
  subnet_ids        = module.networking.private_subnet_ids
  security_group_id = module.security_groups.mq_sg_id
  mq_username       = var.mq_username
  mq_password       = var.mq_password
  instance_type     = "mq.m5.large"
  deployment_mode   = "SINGLE_INSTANCE"
}

module "s3" {
  source             = "../../modules/s3"
  env                = "prod"
  bucket_name        = var.s3_bucket_name
  versioning_enabled = true
  force_destroy      = false
}

module "ses" {
  source        = "../../modules/ses"
  env           = "prod"
  sender_email  = var.ses_sender_email
  s3_bucket_arn = module.s3.bucket_arn
}

data "aws_ami" "ubuntu_22" {
  most_recent = true
  owners      = ["099720109477"] # Canonical
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-22.04-amd64-server-*"]
  }
}

resource "aws_instance" "app" {
  ami                    = data.aws_ami.ubuntu_22.id
  instance_type          = var.app_instance_type
  key_name               = var.key_name
  subnet_id              = module.networking.public_subnet_ids[0]
  vpc_security_group_ids = [module.security_groups.app_sg_id]
  iam_instance_profile   = module.ses.instance_profile_name

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
    encrypted   = true
  }

  tags = { Name = "prod-app" }
}

resource "aws_eip" "app" {
  domain   = "vpc"
  instance = aws_instance.app.id
  tags     = { Name = "prod-app-eip" }
}
