output "app_public_ip" {
  value = aws_eip.app.public_ip
}

output "rds_endpoint" {
  value = module.rds.endpoint
}

output "redis_endpoint" {
  value = "${module.elasticache.endpoint}:${module.elasticache.port}"
}

output "rabbitmq_amqps_endpoint" {
  value = module.amazon_mq.endpoint
}

output "rabbitmq_console_url" {
  value = module.amazon_mq.console_url
}

output "s3_bucket_name" {
  value = module.s3.bucket_name
}

output "s3_bucket_region" {
  value = module.s3.bucket_region
}
