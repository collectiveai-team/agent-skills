---
name: k8s-specialist
description: Kubernetes deployment specialist for production workloads. Covers deployments, security hardening, HPA, NetworkPolicy, RBAC, and production-ready configurations.
---

# Kubernetes Specialist

Production-grade Kubernetes deployment guidance. Apply these patterns when creating or reviewing any K8s manifest, Helm chart, or cluster configuration.

## Deployment Best Practices

### Rolling Update Strategy

Always configure rolling updates to prevent downtime:

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0
```

Setting `maxUnavailable: 0` ensures zero-downtime deployments. Pair with `minReadySeconds: 10` to allow new pods to stabilize before old ones terminate.

### Resource Limits and Requests

Always set both requests and limits. Requests determine scheduling; limits prevent noisy-neighbor issues.

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

Guidelines:
- Set memory limit equal to or close to request to avoid OOMKill surprises
- Set CPU limit 2-5x the request to allow bursting
- Never omit requests; pods without requests get BestEffort QoS and are evicted first
- Use `LimitRange` and `ResourceQuota` at the namespace level as guardrails

### Probes

Configure all three probe types for production workloads:

```yaml
startupProbe:
  httpGet:
    path: /healthz
    port: 8080
  failureThreshold: 30
  periodSeconds: 10
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 20
  failureThreshold: 3
```

- **startupProbe**: Use for slow-starting apps. Prevents liveness from killing pods during init.
- **readinessProbe**: Gates traffic. Pod receives traffic only when ready.
- **livenessProbe**: Restarts stuck processes. Set conservatively to avoid restart loops.

### Pod Disruption Budgets

Always create a PDB for production deployments:

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: app-pdb
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: my-app
```

### Topology Spread and Anti-Affinity

Spread pods across nodes and zones:

```yaml
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        app: my-app
```

## Security Contexts

Apply the principle of least privilege to every container:

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
```

At the pod level, set:

```yaml
spec:
  securityContext:
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault
  automountServiceAccountToken: false
```

Key rules:
- Never run as root (`runAsNonRoot: true`)
- Always drop ALL capabilities, then add only what is needed
- Use read-only root filesystem; mount writable `emptyDir` volumes for tmp/cache
- Disable service account token auto-mount unless the pod needs K8s API access
- Set `seccompProfile: RuntimeDefault` to restrict syscalls

## HPA Configuration

### Standard CPU/Memory Based

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 2
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
```

Guidelines:
- Set `minReplicas >= 2` for HA
- Use scale-down stabilization (300s) to prevent flapping
- Target 70% CPU utilization as a starting point
- Use custom metrics (e.g., requests-per-second) for request-driven workloads via Prometheus adapter

## NetworkPolicy Patterns

### Default Deny All

Apply to every namespace first, then allowlist:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
```

### Allow Specific Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - protocol: TCP
          port: 8080
```

### Allow DNS Egress

Every namespace with default-deny needs DNS egress:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
spec:
  podSelector: {}
  policyTypes:
    - Egress
  egress:
    - to: []
      ports:
        - protocol: UDP
          port: 53
        - protocol: TCP
          port: 53
```

## RBAC Best Practices

- Use `Role` (namespaced) over `ClusterRole` whenever possible
- Bind to `ServiceAccount`, never to `User` or `Group` for workloads
- Create one ServiceAccount per application; never share the `default` SA
- Use `resourceNames` to scope access to specific objects when possible
- Audit with `kubectl auth can-i --list --as=system:serviceaccount:ns:sa`

Example least-privilege Role for a config-reader:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: config-reader
  namespace: app-ns
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list", "watch"]
    resourceNames: ["app-config"]
```

## Helm Chart Conventions

- Use `.Values` for all configurable parameters; never hardcode
- Set sensible defaults in `values.yaml` with comments explaining each value
- Use `_helpers.tpl` for reusable template functions (labels, names, selectors)
- Include standard labels: `app.kubernetes.io/name`, `app.kubernetes.io/version`, `app.kubernetes.io/managed-by`
- Template every resource with `{{ include "chart.fullname" . }}` for naming
- Always include `NOTES.txt` with post-install instructions
- Pin image tags to digests or semver; never use `latest`
- Use `{{- toYaml .Values.resources | nindent 12 }}` for passthrough values
