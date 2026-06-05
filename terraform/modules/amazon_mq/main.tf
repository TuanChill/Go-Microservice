resource "aws_mq_broker" "main" {
  broker_name        = "${var.env}-rabbitmq"
  engine_type        = "RabbitMQ"
  engine_version     = "3.13"
  host_instance_type = var.instance_type
  deployment_mode    = var.deployment_mode

  # SINGLE_INSTANCE requires exactly one subnet
  subnet_ids      = [var.subnet_ids[0]]
  security_groups = [var.security_group_id]

  publicly_accessible = false

  user {
    username = var.mq_username
    password = var.mq_password
  }

  tags = { Name = "${var.env}-rabbitmq" }
}
