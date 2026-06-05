resource "aws_db_subnet_group" "main" {
  name       = "${var.env}-rds-subnet-group"
  subnet_ids = var.subnet_ids
  tags       = { Name = "${var.env}-rds-subnet-group" }
}

resource "aws_db_instance" "main" {
  identifier        = "${var.env}-postgres"
  engine            = "postgres"
  engine_version    = "16"
  instance_class    = var.instance_class
  allocated_storage = var.allocated_storage
  storage_type      = "gp3"
  storage_encrypted = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [var.security_group_id]

  multi_az                  = var.multi_az
  deletion_protection       = var.deletion_protection
  backup_retention_period   = 7
  skip_final_snapshot       = false
  final_snapshot_identifier = "${var.env}-postgres-final-snapshot"

  tags = { Name = "${var.env}-postgres" }
}
