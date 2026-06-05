output "app_sg_id" {
  value = aws_security_group.app.id
}

output "rds_sg_id" {
  value = aws_security_group.rds.id
}

output "cache_sg_id" {
  value = aws_security_group.cache.id
}

output "mq_sg_id" {
  value = aws_security_group.mq.id
}
