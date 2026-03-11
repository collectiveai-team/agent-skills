---
name: security-scanning
description: Security scanning requirements for SAST, SCA, container scanning, and IaC scanning in CI/CD pipelines.
version: 0.1.0
---

# Security Scanning Requirements

Integrate these scanning stages into every CI/CD pipeline. No code reaches production without passing all applicable scans. Findings must be triaged, not ignored.

## SAST (Static Application Security Testing)

Static analysis examines source code for security vulnerabilities without executing it.

### When to Run

- On every pull request (blocking merge for high/critical findings).
- On every merge to main/default branch.
- Scheduled weekly full scans to catch newly-disclosed vulnerability patterns.

### Tools by Language

| Language | Primary Tool | Alternative |
|----------|-------------|-------------|
| Python | Bandit | Semgrep |
| JavaScript/TypeScript | ESLint security plugin + Semgrep | CodeQL |
| Go | gosec | Semgrep |
| Java | SpotBugs + Find Security Bugs | CodeQL |
| Multi-language | Semgrep | CodeQL (GitHub) |

### Configuration Rules

- Enable all security-relevant rule sets. Disable individual rules only with documented justification.
- Configure severity thresholds: block on `HIGH` and `CRITICAL`. Warn on `MEDIUM`.
- Maintain a suppression file (e.g., `.semgrepignore`, `bandit.yaml`) in version control with comments explaining each suppression.
- Review and prune suppressions quarterly. Suppressions older than 6 months must be re-justified.

### Semgrep Example CI Configuration

```yaml
- name: SAST Scan
  run: |
    pip install semgrep
    semgrep scan --config=auto --config=p/security-audit \
      --error --severity ERROR \
      --json --output=semgrep-results.json
```

### Common Findings to Watch For

- SQL injection (string concatenation in queries)
- Command injection (unsanitized input in `os.system`, `exec`, `subprocess`)
- Path traversal (user input in file paths without sanitization)
- Hardcoded secrets (API keys, passwords, tokens in source)
- Insecure deserialization (`pickle.loads`, `yaml.load` without safe loader)
- Server-side request forgery (SSRF) via user-controlled URLs

## SCA (Software Composition Analysis)

SCA scans dependencies for known vulnerabilities (CVEs) and license compliance issues.

### When to Run

- On every pull request that modifies dependency files (requirements.txt, package.json, go.mod, pom.xml).
- Scheduled daily on the default branch to catch newly-disclosed CVEs.
- Before every production deployment.

### Tools

| Ecosystem | Primary Tool | Alternative |
|-----------|-------------|-------------|
| Universal | Snyk | Grype |
| Universal | Trivy (fs mode) | OWASP Dependency-Check |
| Python | pip-audit | Safety |
| Node.js | npm audit | Snyk |
| Go | govulncheck | Trivy |

### Configuration Rules

- Block on `HIGH` and `CRITICAL` vulnerabilities with available fixes.
- For vulnerabilities without fixes, document the risk and set a review date (max 30 days).
- Pin dependencies to exact versions in lock files. Never use floating ranges in production.
- Enable license scanning: block `AGPL`, `GPL` (for proprietary projects), and `SSPL` licenses. Warn on unknown licenses.
- Generate and store SBOM (Software Bill of Materials) for every release using SPDX or CycloneDX format.

### Trivy Example CI Configuration

```yaml
- name: SCA Scan
  run: |
    trivy fs --scanners vuln --severity HIGH,CRITICAL \
      --exit-code 1 --format json --output sca-results.json .
```

### Remediation Process

1. Upgrade to the patched version if available.
2. If no patch exists, evaluate if the vulnerable code path is reachable in your application.
3. If not reachable, suppress with documentation and set a review date.
4. If reachable and no patch exists, implement a workaround or remove the dependency.

## Container Scanning

Scan container images for OS-level and application-level vulnerabilities before deployment.

### When to Run

- After every image build in CI, before pushing to the registry.
- Scheduled daily on all images in the production registry to catch newly-disclosed CVEs.
- As an admission control check (e.g., using Kyverno or OPA Gatekeeper to block unscanned images).

### Tools

| Tool | Strengths |
|------|-----------|
| Trivy | Fast, comprehensive, good CI integration |
| Grype | Fast, SBOM-aware, Anchore ecosystem |
| Snyk Container | Deep analysis, fix suggestions, IDE integration |
| AWS ECR Scanning | Native AWS integration, automatic on push |
| GCP Artifact Analysis | Native GCP integration, on-push scanning |

### Configuration Rules

- Block deployment of images with `CRITICAL` vulnerabilities.
- Set a maximum age for base images: rebuild if the base image is older than 30 days.
- Use minimal base images: `distroless`, `alpine`, or `scratch` where possible.
- Never use `latest` tag. Pin base images to specific digests.
- Scan both OS packages and application dependencies within the image.
- Sign images with Cosign or Notary and verify signatures at deployment.

### Trivy Image Scan Example

```yaml
- name: Container Scan
  run: |
    trivy image --severity HIGH,CRITICAL \
      --exit-code 1 --format json --output image-scan.json \
      ${{ env.IMAGE_NAME }}:${{ env.IMAGE_TAG }}
```

### Base Image Hygiene

- Maintain a list of approved base images. Only build from approved images.
- Automate base image updates with Dependabot or Renovate.
- Multi-stage builds: use a full image for building, a minimal image for runtime.
- Remove build tools, package managers, and shells from runtime images when possible.

## IaC Scanning (Infrastructure as Code)

Scan Terraform, CloudFormation, Kubernetes manifests, and Helm charts for misconfigurations.

### When to Run

- On every pull request that modifies IaC files (blocking merge for high-severity findings).
- Before every `terraform apply` or infrastructure deployment.
- Scheduled weekly on the default branch for drift detection.

### Tools

| Tool | Targets | Strengths |
|------|---------|-----------|
| Checkov | Terraform, CloudFormation, K8s, Helm, Docker | Comprehensive policy library, custom policies |
| tfsec / Trivy config | Terraform | Fast, focused Terraform scanning |
| KICS | Terraform, K8s, CloudFormation, Ansible, Docker | Multi-framework, extensive rules |
| Terrascan | Terraform, K8s, Helm | OPA-based custom policies |
| kube-linter | K8s manifests, Helm | K8s-focused, fast |

### Configuration Rules

- Block on `HIGH` and `CRITICAL` findings. Warn on `MEDIUM`.
- Maintain a skip list in version control with documented justifications for each exception.
- Review skip lists quarterly. Remove stale exceptions.
- Run `terraform plan` output through scanning when possible (catches runtime-evaluated values).

### Checkov Example CI Configuration

```yaml
- name: IaC Scan
  run: |
    pip install checkov
    checkov -d . --framework terraform \
      --soft-fail-on LOW \
      --hard-fail-on HIGH,CRITICAL \
      --output json --output-file iac-scan.json
```

### Common IaC Findings

**Terraform:**
- S3 buckets without encryption or public access blocks
- Security groups with `0.0.0.0/0` ingress on non-HTTP ports
- IAM policies with `*` resource or `*` action
- RDS instances without encryption or Multi-AZ
- Missing CloudTrail or logging configuration

**Kubernetes:**
- Containers running as root
- Missing resource limits or requests
- Missing network policies
- Privileged containers or host network/PID access
- Secrets in environment variables instead of mounted volumes
- Images using `latest` tag

**Docker:**
- Running as root user (missing `USER` instruction)
- Using `ADD` instead of `COPY` (potential remote URL risks)
- Missing health check instructions
- Storing secrets in build args or environment variables

## Triage and Reporting

### Severity-Based Response Times

| Severity | Response Time | Resolution Time |
|----------|--------------|-----------------|
| CRITICAL | Same day | 48 hours |
| HIGH | 3 business days | 2 weeks |
| MEDIUM | 2 weeks | 1 month |
| LOW | Next sprint | Best effort |

### Reporting

- Aggregate all scan results into a single dashboard (Snyk, Defect Dojo, or custom).
- Track mean time to remediation (MTTR) by severity.
- Report vulnerability trends monthly to engineering leadership.
- Maintain a risk register for accepted vulnerabilities with documented business justification and review dates.
