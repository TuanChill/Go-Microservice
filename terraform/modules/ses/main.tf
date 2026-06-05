# Verify sender identity — AWS sends a confirmation email; click the link to activate
resource "aws_sesv2_email_identity" "sender" {
  email_identity = var.sender_email
  tags           = { Name = "${var.env}-ses-sender" }
}

# Generic IAM role for the app EC2 instance — holds all AWS service permissions
resource "aws_iam_role" "app" {
  name = "${var.env}-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = { Name = "${var.env}-app-role" }
}

# SES send permission scoped to the verified sender address
resource "aws_iam_role_policy" "ses_send" {
  name = "${var.env}-ses-send-policy"
  role = aws_iam_role.app.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "ses:SendEmail",
        "ses:SendRawEmail",
        "sesv2:SendEmail",
      ]
      Resource = "*"
      Condition = {
        StringEquals = {
          "ses:FromAddress" = var.sender_email
        }
      }
    }]
  })
}

# S3 access is managed via S3 bucket policies (see modules/s3/main.tf)
# The instance profile grants the EC2 access to SES.
# For S3, either use bucket policies or pass s3_bucket_arn to grant IAM access below.

# IAM instance profile — attached to the app EC2 so AWS SDK picks up credentials via IMDSv2
resource "aws_iam_instance_profile" "app" {
  name = "${var.env}-app-instance-profile"
  role = aws_iam_role.app.name
}

# Output the IAM role ARN so S3 module can grant access via bucket policy
output "iam_role_arn" {
  value = aws_iam_role.app.arn
}
