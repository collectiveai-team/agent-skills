---
name: terraform-conventions
description: Terraform coding conventions covering file organization, state management, naming, variables, and module structure.
version: 0.1.0
---

# Terraform Conventions

Apply these conventions to all Terraform code. Consistency enables team velocity, reduces review friction, and prevents state corruption.

## File Organization

### Standard File Layout

Every Terraform root module and child module must follow this file structure:

```
module/
  main.tf          # Primary resource definitions
  variables.tf     # All input variable declarations
  outputs.tf       # All output declarations
  versions.tf      # Required providers and terraform version constraints
  locals.tf        # Local values (optional, when locals are needed)
  data.tf          # Data sources (optional, when data sources are needed)
  backend.tf       # Backend configuration (root modules only)
```

### Rules

- `main.tf` contains the core resource definitions. Split into multiple files by logical grouping (e.g., `networking.tf`, `compute.tf`, `iam.tf`) when `main.tf` exceeds 200 lines.
- `variables.tf` contains ALL variable declarations. Never define variables inline in other files.
- `outputs.tf` contains ALL output declarations. Never define outputs inline in other files.
- `versions.tf` contains the `terraform` block with `required_version` and `required_providers`. Pin provider versions to minor version ranges: `~> 5.0`.
- Never put resource definitions in `variables.tf` or `outputs.tf`.
- Group related resources together within a file. Separate groups with a blank line and a comment.

### Example `versions.tf`

```hcl
terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

## State Management

### Remote Backend

- Always use a remote backend for shared projects. Never commit `terraform.tfstate` to git.
- Use S3 + DynamoDB for AWS. Use GCS for GCP. Use Azure Blob for Azure.
- Enable state locking to prevent concurrent modifications.
- Enable versioning on the state bucket for recovery from corruption.
- Encrypt state at rest (S3 SSE, GCS CMEK).

### Example S3 Backend

```hcl
terraform {
  backend "s3" {
    bucket         = "myorg-terraform-state"
    key            = "prod/networking/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}
```

### State Organization

- Use one state file per environment per service. Never share state across environments.
- Use the key path pattern: `{environment}/{service}/terraform.tfstate`.
- Keep state files small. If `terraform plan` takes more than 2 minutes, split the state.
- Use `terraform_remote_state` data sources or SSM parameters to share outputs between states.

### State Operations

- Never use `terraform state rm` or `terraform state mv` without peer review.
- Use `terraform import` to bring existing resources under management.
- Before `terraform destroy`, always run `terraform plan -destroy` and review.
- Use `terraform workspace` only for identical environments. Prefer separate state files for environments with different configurations.

## Naming Conventions

### Resources and Data Sources

- Use `snake_case` for all resource names, variable names, output names, and local values.
- Use descriptive, specific names: `aws_security_group.api_ingress` not `aws_security_group.sg1`.
- Prefix resources with their purpose: `aws_iam_role.lambda_execution`, `aws_s3_bucket.log_archive`.
- Use `this` as the resource name only when a module creates a single primary resource of that type.

### Variables

- Prefix boolean variables with `enable_` or `is_`: `enable_monitoring`, `is_public`.
- Use plural names for lists and maps: `subnet_ids`, `tags`.
- Use `_id` suffix for resource identifiers: `vpc_id`, `subnet_id`.
- Use `_arn` suffix for ARNs: `role_arn`, `bucket_arn`.

### Tags

- Apply consistent tags to every resource. Use a `local.common_tags` pattern:

```hcl
locals {
  common_tags = {
    Environment = var.environment
    Service     = var.service_name
    Team        = var.team
    ManagedBy   = "terraform"
  }
}
```

- Merge common tags with resource-specific tags: `tags = merge(local.common_tags, { Name = "specific-name" })`.

## Variable Validation

### Rules

- Every variable must have a `description`. No exceptions.
- Every variable must have a `type` constraint.
- Use `validation` blocks for variables with specific constraints.
- Set `default` values only when a sensible default exists. Required variables should not have defaults.
- Use `sensitive = true` for secrets, tokens, and passwords.
- Never pass secrets as variable defaults. Use environment variables or secret managers.

### Examples

```hcl
variable "environment" {
  description = "Deployment environment (dev, staging, prod)"
  type        = string

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "instance_count" {
  description = "Number of EC2 instances to create"
  type        = number
  default     = 1

  validation {
    condition     = var.instance_count > 0 && var.instance_count <= 20
    error_message = "Instance count must be between 1 and 20."
  }
}

variable "cidr_blocks" {
  description = "List of CIDR blocks for ingress rules"
  type        = list(string)

  validation {
    condition     = alltrue([for cidr in var.cidr_blocks : can(cidrhost(cidr, 0))])
    error_message = "All values must be valid CIDR blocks."
  }
}

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}
```

## Module Structure

### When to Extract a Module

Extract a module when:
- The same resource pattern is used 3+ times across configurations.
- A group of resources represents a logical unit (e.g., VPC with subnets, route tables, NAT gateways).
- You want to enforce a standard pattern with constrained inputs.

Do not extract a module when:
- The pattern is used only once. Inline it.
- The module would have as many variables as the resources have attributes (passthrough module).
- The abstraction hides complexity that operators need to understand.

### Module Layout

```
modules/
  vpc/
    main.tf
    variables.tf
    outputs.tf
    versions.tf
```

### Module Versioning

- Use versioned module sources for shared modules: `source = "git::https://github.com/org/modules.git//vpc?ref=v1.2.0"`.
- Pin to exact tags, not branches. Never use `ref=main`.
- For local modules within the same repo, use relative paths: `source = "../../modules/vpc"`.
- Publish reusable modules to a private Terraform registry when adoption grows beyond 2-3 teams.

### Module Design Rules

- Expose the minimum necessary variables. Use opinionated defaults for internal complexity.
- Output all values that downstream consumers might need (IDs, ARNs, names, endpoints).
- Include a `README.md` with usage examples for shared modules.
- Never hardcode account IDs, regions, or environment names inside modules. Pass them as variables.
- Use `for_each` over `count` for resources that might be conditionally created or iterated. `for_each` produces stable resource addresses.

### Terraform Plan Review

- Always run `terraform plan` before `terraform apply`. Review every change.
- Use `terraform plan -out=plan.tfplan` and `terraform apply plan.tfplan` in CI/CD for deterministic applies.
- Flag any plan that shows `destroy` or `replace` for review by a second engineer.
- Use `-target` sparingly and only for debugging. Never use `-target` in production applies.
