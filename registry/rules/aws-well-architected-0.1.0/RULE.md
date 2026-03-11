---
name: aws-well-architected
description: AWS Well-Architected Framework guidelines covering all five pillars for production infrastructure.
version: 0.1.0
---

# AWS Well-Architected Framework

Apply these principles when designing, reviewing, or modifying any AWS infrastructure. Every architectural decision should be evaluated against all five pillars.

## Operational Excellence

### Infrastructure as Code

- Define ALL infrastructure in CloudFormation or Terraform. No manual console changes in production.
- Use `aws cloudformation drift-detection` to detect manual drift.
- Store IaC in version control with PR review required before apply.
- Tag every resource with `Environment`, `Team`, `Service`, and `CostCenter`.

### Monitoring and Observability

- Enable CloudWatch alarms for every production service:
  - API Gateway: 5xx error rate > 1%, latency p99 > 2s
  - Lambda: error rate > 1%, duration > 80% of timeout
  - RDS: CPU > 80%, free storage < 20%, replica lag > 10s
  - ECS/EKS: CPU/memory utilization > 80%
- Use CloudWatch Logs Insights or OpenSearch for centralized logging.
- Enable AWS X-Ray or OpenTelemetry for distributed tracing.
- Create operational dashboards per service in CloudWatch or Grafana.

### Runbooks and Automation

- Document runbooks for every alarm: what it means, how to investigate, how to remediate.
- Use Systems Manager Automation for common operational tasks.
- Implement auto-remediation for known failure patterns (e.g., restart unhealthy tasks).
- Practice incident response with game days at least quarterly.

## Security

### Identity and Access Management

- Use IAM roles, never long-lived access keys. Rotate any remaining keys every 90 days.
- Apply least privilege: start with zero permissions, add only what is needed.
- Use permission boundaries to cap maximum permissions for delegated admin roles.
- Enable MFA for all human users. Require MFA for sensitive operations.
- Use AWS Organizations SCPs to enforce guardrails across accounts.
- Separate accounts by environment: dev, staging, production.

### Encryption

- Encrypt all data at rest: S3 (SSE-S3 or SSE-KMS), EBS (default encryption), RDS (encrypted instances).
- Encrypt all data in transit: enforce TLS 1.2+ on all endpoints.
- Use AWS KMS with customer-managed keys for sensitive workloads.
- Enable S3 bucket policies that deny `s3:PutObject` without encryption headers.

### Network Security

- Use VPC for all resources. No public subnets for databases or application servers.
- Security groups: allow only required ports, reference other security groups instead of CIDR ranges where possible.
- Use VPC endpoints (Gateway and Interface) for AWS service access to avoid internet traversal.
- Enable VPC Flow Logs for network auditing.
- Use AWS WAF on public-facing ALBs and API Gateways.

### Detection and Response

- Enable AWS CloudTrail in all regions and all accounts. Send to a centralized logging account.
- Enable AWS GuardDuty for threat detection.
- Enable AWS Config with conformance packs for compliance monitoring.
- Use AWS Security Hub to aggregate findings from GuardDuty, Inspector, and Config.

## Reliability

### Multi-AZ and Multi-Region

- Deploy all production workloads across at least 2 AZs.
- Use ALB/NLB for cross-AZ load balancing.
- RDS: use Multi-AZ deployments. For critical databases, use read replicas in another region.
- S3: already cross-AZ. Enable cross-region replication for disaster recovery.
- Define RTO and RPO for every service and validate with DR testing.

### Health Checks and Auto-Recovery

- Configure ALB health checks on application-specific endpoints, not just TCP.
- ECS: use container health checks in task definitions.
- EC2: enable auto-recovery actions for system status check failures.
- Use Route 53 health checks for DNS failover to secondary regions.

### Backups

- Enable AWS Backup for all stateful resources (RDS, EBS, DynamoDB, EFS).
- Test backup restoration quarterly. Document restoration procedures.
- Use DynamoDB point-in-time recovery (PITR).
- S3 versioning enabled on all production buckets.

### Throttling and Circuit Breakers

- Implement client-side retries with exponential backoff and jitter.
- Use API Gateway throttling to protect backend services.
- Set SQS/Lambda concurrency limits to prevent downstream overload.

## Performance Efficiency

### Right-Sizing

- Use AWS Compute Optimizer for EC2, Lambda, and EBS recommendations.
- Start small and scale up based on observed metrics, not guesses.
- Review instance types quarterly; switch to Graviton (ARM) for cost/performance gains.
- Use appropriate storage: gp3 over gp2, io2 only when IOPS matter.

### Caching

- Use ElastiCache (Redis/Memcached) for frequently-read, infrequently-changed data.
- Use DynamoDB DAX for DynamoDB-heavy read workloads.
- Implement application-level caching with TTLs matched to data freshness requirements.

### Content Delivery

- Use CloudFront for all static assets and public APIs.
- Enable compression (gzip/brotli) on CloudFront distributions.
- Use Lambda@Edge or CloudFront Functions for lightweight request transformations.
- Set appropriate Cache-Control headers at the origin.

### Database Performance

- Choose the right database engine for the access pattern: relational (RDS/Aurora), key-value (DynamoDB), document (DocumentDB), graph (Neptune).
- Use RDS Proxy for Lambda-to-RDS connections to manage connection pooling.
- Enable Performance Insights for RDS query analysis.

## Cost Optimization

### Reserved Capacity

- Use Savings Plans (Compute or EC2) for steady-state workloads. Target 1-year commitment first.
- Use Reserved Instances for RDS, ElastiCache, and OpenSearch steady-state.
- Use Spot Instances for fault-tolerant workloads (batch processing, CI/CD runners).

### Resource Lifecycle

- Implement tagging policies and enforce with AWS Config rules.
- Use AWS Cost Explorer and create monthly budgets with alerts at 80% and 100%.
- Delete unused resources: unattached EBS volumes, old snapshots, idle NAT gateways, unused Elastic IPs.
- Use S3 lifecycle policies to transition to Glacier/Deep Archive.
- Right-size or delete non-production environments on nights and weekends.

### Architecture for Cost

- Use serverless (Lambda, Fargate, DynamoDB on-demand) for variable/unpredictable workloads.
- Use S3 Intelligent-Tiering for objects with unknown access patterns.
- Consolidate AWS accounts under Organizations for volume discounts.
- Review Cost and Usage Reports weekly in early stages, monthly once stable.
