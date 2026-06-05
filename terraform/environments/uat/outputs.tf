output "app_public_ip" {
  value = aws_eip.app.public_ip
}

output "ssh_command" {
  value = "ssh ubuntu@${aws_eip.app.public_ip}"
}

output "rabbitmq_console" {
  value = "http://${aws_eip.app.public_ip}:15672"
}

output "s3_bucket_name" {
  value = module.s3.bucket_name
}

output "s3_bucket_region" {
  value = module.s3.bucket_region
}
