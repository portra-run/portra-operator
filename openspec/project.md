# Project Context

## Purpose
Portra is a self-hosted PaaS (Platform as a Service) designed specifically for Kubernetes. It aims to provide an experience similar to Coolify, Dokploy, or Caprover, but leveraging native Kubernetes resources. The core of the project is a Kubernetes Operator that manages the "App" Custom Resource Definition (CRD), automating the creation and lifecycle of deployments, persistent storage, and networking configuration.

## Tech Stack
- **Language:** Go (v1.25+)
- **Framework:** [Kubebuilder](https://book.kubebuilder.io/) / [Operator SDK](https://sdk.operatorframework.io/)
- **Core Libraries:** 
  - `sigs.k8s.io/controller-runtime`
  - `k8s.io/client-go`
- **Testing:** Ginkgo and Gomega
- **Infrastructure/Manifests:** Kustomize
- **Containerization:** Docker

## Project Conventions

### Code Style
- Follow standard Go idioms and formatting (`go fmt`).
- Use Kubebuilder markers for CRD validation and RBAC generation.
- Proper error handling using `fmt.Errorf` with context.

### Architecture Patterns
- **Operator Pattern:** Uses a controller to reconcile the desired state (App CRD) with the actual state of the cluster.
- **Declarative API:** Users define their application requirements in a YAML manifest (App spec).
- **Resource Ownership:** All generated resources (Deployments, PVCs, etc.) should be owned by the App CR to ensure automatic cleanup via garbage collection.

### Testing Strategy
- **Unit/Integration Tests:** Use EnvTest with Ginkgo/Gomega to verify controller logic without a full cluster.
- **E2E Tests:** Use Kind or a real cluster to verify the full lifecycle of an application.

### Git Workflow
- Feature branches for all changes.
- Descriptive commit messages following the Conventional Commits pattern where possible.

## Domain Context
- **App CRD:** The primary unit of management. It encapsulates the application's configuration, including container images, environment variables, storage requirements, and exposure (API Gateway).
- **Orchestration:** The operator is responsible for creating:
    - **Deployments:** For running the application containers.
    - **PersistentVolumeClaims (PVCs):** For stateful application data.
    - **Networking/Gateway:** Configuration for the API Gateway (e.g., Ingress, Gateway API, or specific controller-related CRDs).

## Important Constraints
- Must be compatible with modern Kubernetes clusters (v1.30+ recommended).
- Designed for self-hosting with minimal external dependencies.
- Security-first approach for multi-tenant application isolation.

## External Dependencies
- **Kubernetes API:** The primary interface for all operations.
- **API Gateway/Ingress Controller:** (e.g., Traefik, NGINX, or Envoy Gateway) for external access.
- **CSI Drivers:** For persistent storage provisioning.
