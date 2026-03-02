# Design: Implement Basic Functional App CRD

## Architectural Decisions

### CRD Fields Overview
The `App` CRD will expose the following fields under `.spec`:
- `image` (string): Docker image URL.
- `containerPort` (int32): Main port on which the container listens.
- `domains` ([]string): List of DNS names to expose the app on.
- `tls` (bool): If true, an Ingress with TLS/cert-manager support will be created.
- `env` ([]corev1.EnvVar): Environment variables for the container.
- `replicas` (int32): Number of pod replicas.
- `resources` (corev1.ResourceRequirements): CPU/Memory requests and limits.

### Kubernetes Resources Mapping
When an `App` is created:
1.  **Deployment:** Will use the `image`, `env`, `replicas`, and `resources` fields.
2.  **Service:** Will expose the `containerPort`.
3.  **Ingress:** Will map `domains` to the `Service`. If `tls` is true, it will include the `tls` section and potentially a cert-manager annotation.

### Ownership and Garbage Collection
- Every resource (Deployment, Service, Ingress) must have an `OwnerReference` pointing to the `App` CR. This ensures that when an `App` is deleted, all its Kubernetes resources are also removed.

### TLS Strategy
For now, if `tls: true`, we will add `cert-manager.io/cluster-issuer` or `cert-manager.io/issuer` annotations to the Ingress resource. We will assume a default issuer exists on the cluster.

### Error Handling
If the reconciliation fails (e.g., invalid image or domain), the `AppStatus` will be updated with a `Degraded` condition and a descriptive message.

## Trade-offs
- **Fixed Ingress annotations:** We might need more flexibility in the future, but for the first version, we'll hardcode or use a sane default for cert-manager integration.
- **Single container only:** We only support one container per `App` for now.
