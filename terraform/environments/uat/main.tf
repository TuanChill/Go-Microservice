provider "aws" {
  region = var.aws_region
}

# Uncomment after bootstrapping state bucket:
#   aws s3 mb s3://YOUR_BUCKET --region ap-southeast-1
#   aws dynamodb create-table --table-name terraform-locks \
#     --attribute-definitions AttributeName=LockID,AttributeType=S \
#     --key-schema AttributeName=LockID,KeyType=HASH \
#     --billing-mode PAY_PER_REQUEST --region ap-southeast-1
#
# terraform {
#   backend "s3" {
#     bucket         = "YOUR_BUCKET"
#     key            = "uat/terraform.tfstate"
#     region         = "ap-southeast-1"
#     dynamodb_table = "terraform-locks"
#     encrypt        = true
#   }
# }

module "networking" {
  source              = "../../modules/networking"
  env                 = "uat"
  vpc_cidr            = "10.0.0.0/16"
  public_subnet_cidrs = ["10.0.1.0/24"]
  azs                 = ["${var.aws_region}a"]
  enable_nat_gateway  = false
}

module "security_groups" {
  source           = "../../modules/security_groups"
  env              = "uat"
  vpc_id           = module.networking.vpc_id
  allowed_ssh_cidr = var.allowed_ssh_cidr
  app_port         = var.app_port
}

module "s3" {
  source             = "../../modules/s3"
  env                = "uat"
  bucket_name        = var.s3_bucket_name
  versioning_enabled = false
  force_destroy      = true
}

module "ses" {
  source        = "../../modules/ses"
  env           = "uat"
  sender_email  = var.ses_sender_email
  s3_bucket_arn = module.s3.bucket_arn
}

# Ubuntu 22.04 LTS via AWS SSM public parameter (Canonical's official path)
# Falls back to hardcoded AMI if SSM parameter is unavailable in this account.
data "aws_ssm_parameter" "ubuntu_22_ami" {
  count = var.ami_id == "" ? 1 : 0
  name  = "/aws/service/canonical/ubuntu/server/22.04/stable/current/amd64/hvm/ebs-gp3/ami-id"
}

# When ami_id is provided, use it directly. When empty, use SSM parameter.
resource "aws_instance" "app" {
  ami = length(var.ami_id) > 0 ? var.ami_id : data.aws_ssm_parameter.ubuntu_22_ami[0].value
  instance_type          = var.instance_type
  key_name               = var.key_name
  subnet_id              = module.networking.public_subnet_ids[0]
  vpc_security_group_ids = [module.security_groups.app_sg_id]
  iam_instance_profile   = module.ses.instance_profile_name

  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }

  user_data = templatefile("${path.module}/user_data.sh.tpl", {
    db_name         = var.db_name
    db_user         = var.db_user
    db_password     = var.db_password
    redis_password  = var.redis_password
    rabbit_user     = var.rabbit_user
    rabbit_password = var.rabbit_password
  })

  tags = { Name = "uat-app" }
}

resource "aws_eip" "app" {
  domain   = "vpc"
  instance = aws_instance.app.id
  tags     = { Name = "uat-app-eip" }
}
