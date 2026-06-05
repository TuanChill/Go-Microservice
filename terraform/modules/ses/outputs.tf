output "instance_profile_name" {
  description = "IAM instance profile name — attach to EC2 so it can call SES and S3 without access keys"
  value       = aws_iam_instance_profile.app.name
}

output "sender_email" {
  value = var.sender_email
}
