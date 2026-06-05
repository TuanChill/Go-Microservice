# Terraform Infrastructure

Infrastructure-as-Code for deploying the Go backend template on AWS.

## Structure

```
terraform/
├── versions.tf                        # Terraform & AWS provider constraints
├── modules/
│   ├── networking/                    # VPC, subnets, IGW, NAT Gateway
│   ├── security_groups/              # App, RDS, ElastiCache, Amazon MQ SGs
│   ├── rds/                          # RDS PostgreSQL 16
│   ├── elasticache/                  # ElastiCache Redis 7.1
│   ├── amazon_mq/                    # Amazon MQ RabbitMQ 3.13
│   ├── ses/                          # SES email identity + IAM role (EC2 metadata auth)
│   └── s3/                           # S3 bucket for file uploads
└── environments/
    ├── uat/                          # Single EC2 + Docker Compose (PostgreSQL, Redis, RabbitMQ)
    └── prod/                         # EC2 app + RDS, ElastiCache, Amazon MQ, S3
```

## Prerequisites

1. **Terraform** >= 1.6
2. **AWS CLI** configured with credentials (`aws configure`)
3. **Key pair** created in the target region:
   ```bash
   aws ec2 create-key-pair --key-name my-key-pair \
     --region ap-southeast-1 \
     --query 'KeyMaterial' \
     --output text > ~/.ssh/my-key-pair.pem
   chmod 400 ~/.ssh/my-key-pair.pem
   ```
4. **State backend** (optional but recommended):
   ```bash
   aws s3 mb s3://YOUR_BUCKET --region ap-southeast-1
   aws dynamodb create-table \
     --table-name terraform-locks \
     --attribute-definitions AttributeName=LockID,AttributeType=S \
     --key-schema AttributeName=LockID,KeyType=HASH \
     --billing-mode PAY_PER_REQUEST \
     --region ap-southeast-1
   ```

---

## Deploy UAT

```bash
cd terraform/environments/uat
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars — fill in key_name, passwords, and ses_sender_email
terraform init
terraform plan
terraform apply
```

**What gets created:**
- VPC `10.0.0.0/16` with 1 public subnet
- EC2 t3.small (Ubuntu 22.04) with Docker Compose pre-installed via user_data
- Docker Compose on EC2: PostgreSQL 16 + Redis 7 + RabbitMQ 3.13
- SES email identity (verification email sent — click the link)
- IAM instance profile: SES + S3 permissions
- S3 bucket (versioning off, force_destroy=true — safe for UAT)
- EIP attached to EC2

**After apply — configure the app:**

```bash
# Get the EIP
terraform output

# Point YAML config at the UAT server
# configs/yaml/config.dev.yml:
ses:
  region: "ap-southeast-1"
  sender: "noreply@yourdomain.com"

s3:
  region: "ap-southeast-1"
  bucket: "my-app-uat-uploads"  # from terraform output

database:
  host: "<terraform output app_public_ip>"
  port: "5432"
```

## Deploy Prod

```bash
cd terraform/environments/prod
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars — fill in key_name, passwords, ses_sender_email, s3_bucket_name
terraform init
terraform plan
terraform apply
```

**What gets created:**
- VPC `10.1.0.0/16` with 2 public + 2 private subnets across 2 AZs
- NAT Gateway + route tables for private subnet routing
- RDS PostgreSQL 16 (db.t3.small, Single-AZ, encrypted, deletion_protection=true)
- ElastiCache Redis 7.1 (cache.t3.micro, single node)
- Amazon MQ RabbitMQ 3.13 (mq.m5.large, single instance, not publicly accessible)
- SES email identity + IAM role
- S3 bucket (versioning on, lifecycle → IA 30d → Glacier 90d)
- EC2 t3.small (app server, no Docker Compose — deploy binary directly)
- EIP + IAM instance profile

**After apply — configure the app:**

```bash
terraform output
# Use RDS endpoint, ElastiCache endpoint, Amazon MQ AMQPS endpoint
# in configs/yaml/config.prod.yml
```

## Clean Up

```bash
# UAT (force_destroy=true, so direct delete works)
cd terraform/environments/uat
terraform destroy

# Prod (deletion_protection on RDS — must disable first)
# Set deletion_protection=false in prod/main.tf, then:
terraform destroy
```

---

## Cost Estimation

### UAT (~$21/month)

| Resource | Spec | Hourly | Monthly |
|---|---|---|---|
| EC2 t3.small | 2 vCPU, 2 GiB | $0.027 | ~$20 |
| EBS gp3 storage | 20 GiB | ~$0.003 | ~$2 |
| **Total** | | | **~$22** |

- No NAT Gateway (public subnet, EIP only)
- Docker Compose: PostgreSQL 16 + Redis 7 + RabbitMQ 3.13 (all on EC2, no extra cost)
- SES: free tier eligible (<62K emails/month), beyond that $0.10/1K
- S3: negligible (<1 GB, no versioning)

### Prod (~$385/month)

| Resource | Spec | Hourly | Monthly |
|---|---|---|---|
| EC2 t3.small | 2 vCPU, 2 GiB | $0.027 | ~$20 |
| RDS db.t3.small PostgreSQL | Single-AZ (no Multi-AZ) | $0.056 | ~$40 |
| ElastiCache cache.t3.micro | Redis 7.1, single node | $0.025 | ~$18 |
| Amazon MQ mq.m5.large | RabbitMQ 3.13, single | $0.36 | ~$259 |
| S3 | <1 GB, versioning + lifecycle | ~$0.03 | ~$1 |
| Data transfer | Estimate | ~$0.05 | ~$4 |
| **Total** | | | **~$342** |

> **Amazon MQ is the biggest cost driver.** At $259/month it represents ~65% of the total. If budget is tight, consider running RabbitMQ on a larger EC2 instance (m5.xlarge ~$0.19/hr → ~$137/month) instead of Amazon MQ — you'd need to manage failover yourself, but the savings are significant.

### Cost Optimization Tips

1. **UAT** — stop EC2 when not in use: `aws ec2 stop-instances --instance-id <id>` (EBS data persists, ~$5/month while stopped)
2. **Prod** — consider `t3.micro` for even cheaper ElastiCache if Redis is not critical (~$9/month)
3. **SES sandbox** — stay in sandbox mode for dev/testing (no production request needed yet)
4. **RDS Multi-AZ** — only enable if you need <1 min RTO; adds ~$40/month
5. **Reserved Instances** — for prod, 1-year reserved saves ~40%: EC2 ~$14/mo, RDS ~$23/mo, ElastiCache ~$6/mo

---

## IAM Auth (No Access Keys)

The app uses **EC2 Instance Metadata Service (IMDSv2)** — no AWS access keys stored anywhere:

```
EC2 instance → IAM instance profile → STS assume role → SES / S3
```

Credentials are picked up automatically by the AWS SDK from `$AWS_*` env vars or EC2 metadata.

## SES Verification

After `terraform apply`, AWS sends a verification email to the `ses_sender_email` address. **Click the link** in that email to activate — SES will not send to unverified addresses in sandbox mode.

To request production access: [AWS SES Console → Account → Request production access](https://console.aws.amazon.com/ses/home?#/account)

## State Backend (Optional)

Uncomment the `backend "s3"` block in each `main.tf` after bootstrapping the S3 + DynamoDB table (instructions in the comments).

---

## Module Summary

| Module | UAT | Prod | Notes |
|---|---|---|---|
| VPC | 1 AZ, no NAT | 2 AZs, 1 shared NAT GW | NAT = cost + complexity |
| PostgreSQL | Docker on EC2 | RDS db.t3.small Single-AZ | RDS has automated backups |
| Redis | Docker on EC2 | ElastiCache cache.t3.micro | |
| RabbitMQ | Docker on EC2 | Amazon MQ mq.m5.large | |
| SES | SESv2 identity + IAM | Same | Production mode needed |
| S3 | No versioning, force_destroy | Versioning + lifecycle | |