---
name: k8s-best-practices
description: Kubernetes best practices for production workloads covering resource management, autoscaling, health checks, security, and network policies.
version: 0.1.0
---

# Kubernetes Best Practices

Apply these rules to every Kubernetes manifest, Helm chart, and cluster configuration. These are non-negotiable for production workloads.

## Resource Limits and Requests

Always set both `requests` and `limits` for every container. Pods without requests receive BestEffort QoS and are evicted first under memory pressure.

### Rules

- Set `requests.memory` to the observed steady-state usage of the application.
- Set `limits.memory` to 1.5x-2x the request to accommodate spikes without OOMKill.
- Set `requests.cpu` to the average CPU usage. This determines scheduling.
- Set `limits.cpu` to 2x-5x the request to allow bursting. Omit CPU limits if the cluster uses CPU throttling instead of hard caps (evaluate per cluster).
- Never deploy without resource requests. Enforce via `LimitRange` and `ResourceQuota` on namespaces.

### Example

```yaml
resources:
  requests:
    cpu: 100m
    memory: 256Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

### Namespace-Level Guardrails

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
spec:
  limits:
    - default:
        cpu: 500m
        memory: 512Mi
      defaultRequest:
        cpu: 100m
        memory: 128Mi
      type: Container
```

## HPA Configuration

### Rules

- Set `minReplicas >= 2` for production workloads to maintain availability during scaling events and node failures.
- Target 70% CPU utilization as a starting point. Adjust based on observed scaling behavior.
- Use the `autoscaling/v2` API for access to multiple metrics and scaling behavior controls.
- Configure scale-down stabilization to at least 300 seconds to prevent flapping.
- Limit scale-up rate to prevent sudden spikes from overwhelming downstream dependencies.
- For request-driven services, prefer custom metrics (requests-per-second, queue depth) over CPU-based scaling.

### Example

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app
  minReplicas: 2
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
```

### VPA Considerations

- Do not run VPA in `Auto` mode alongside HPA on the same metric. They will conflict.
- Use VPA in `Off` mode to get recommendations without automatic changes.
- Apply VPA recommendations to update resource requests during maintenance windows.

## Liveness, Readiness, and Startup Probes

### Rules

- Every production pod must have both a readinessProbe and a livenessProbe.
- Use startupProbe for applications that take more than 10 seconds to initialize.
- Readiness probes gate traffic. They should check that the application can serve requests (e.g., database connection established, cache warmed).
- Liveness probes detect deadlocks and stuck processes. They should check basic health, not downstream dependencies.
- Never point liveness and readiness probes at the same endpoint if that endpoint checks downstream dependencies. A failing database should remove traffic (readiness) but not restart the pod (liveness).
- Set liveness probe intervals conservatively (20s+) with high failure thresholds (3+) to avoid unnecessary restarts.

### Example

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
  periodSeconds: 10
  failureThreshold: 3
  successThreshold: 1
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  periodSeconds: 20
  failureThreshold: 3
  initialDelaySeconds: 0
```

### Endpoint Guidelines

- `/healthz` (liveness): Return 200 if the process is running and not deadlocked. Do not check external dependencies.
- `/ready` (readiness): Return 200 if the app can handle requests. Check database connectivity, cache availability, required config loaded.
- Use `tcpSocket` probes only if the application does not support HTTP health endpoints.
- Use `exec` probes as a last resort; they spawn a process in the container on every check.

## Security Contexts

### Rules

- Always set `runAsNonRoot: true`. Never run containers as root.
- Set `readOnlyRootFilesystem: true`. Mount writable `emptyDir` volumes for tmp/cache/log paths.
- Set `allowPrivilegeEscalation: false` on every container.
- Drop ALL capabilities, then add back only what is needed (rare in most applications).
- Set `seccompProfile: RuntimeDefault` at the pod level.
- Disable `automountServiceAccountToken` unless the pod needs Kubernetes API access.
- Create a dedicated ServiceAccount for each application. Never use the `default` ServiceAccount.

### Example

```yaml
spec:
  serviceAccountName: my-app-sa
  automountServiceAccountToken: false
  securityContext:
    runAsNonRoot: true
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: app
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
        readOnlyRootFilesystem: true
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - ALL
      volumeMounts:
        - name: tmp
          mountPath: /tmp
  volumes:
    - name: tmp
      emptyDir: {}
```

### Enforcement

- Use Pod Security Standards (PSS) at the namespace level with `restricted` profile for production.
- Use OPA Gatekeeper or Kyverno to enforce custom policies (e.g., no `latest` tags, required labels).

## NetworkPolicy

### Rules

- Apply a default-deny policy to every namespace first. Then create explicit allow rules.
- Every namespace must have both ingress and egress default-deny.
- Always allow DNS egress (UDP/TCP port 53) to kube-dns or the cluster DNS service.
- Scope policies as tightly as possible: specify both `podSelector` and `ports`.
- Use `namespaceSelector` with labels to control cross-namespace traffic.
- Document the intended traffic flow for each service before writing policies.

### Default Deny Template

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: my-namespace
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
  namespace: my-namespace
spec:
  podSelector: {}
  policyTypes:
    - Egress
  egress:
    - ports:
        - protocol: UDP
          port: 53
        - protocol: TCP
          port: 53
```

### Service-to-Service Allow Template

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-api
spec:
  podSelector:
    matchLabels:
      app: api
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

### Testing

- Use `kubectl run` with `--rm` to test connectivity after applying policies.
- Verify both allowed and denied paths. A misconfigured policy that allows everything is worse than no policy.
- Use a CNI that supports NetworkPolicy (Calico, Cilium, Weave Net). The default kubenet does not enforce policies.
