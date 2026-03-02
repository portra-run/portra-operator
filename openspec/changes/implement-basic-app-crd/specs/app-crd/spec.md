# Capability: Functional App CRD

## ADDED Requirements

### Requirement: Define App Specification
The `App` CRD must allow users to define the basic configuration for their application.

#### Scenario: Define a basic application
- **Given** an `App` resource manifest.
- **When** the user specifies `image`, `containerPort`, `domains`, `tls`, `env`, `replicas`, and `resources`.
- **Then** the Kubernetes API should accept and store this configuration.

### Requirement: Reconcile Deployment
The operator must create or update a Kubernetes `Deployment` based on the `App` spec.

#### Scenario: Create deployment for an app
- **Given** an `App` with `image: "my-app:v1"`, `replicas: 3`, and specific `resources`.
- **When** the operator reconciles the `App`.
- **Then** a `Deployment` should be created with the specified image, replicas, and resource limits.

### Requirement: Reconcile Service
The operator must create a Kubernetes `Service` to expose the application within the cluster.

#### Scenario: Expose app via service
- **Given** an `App` with `containerPort: 8080`.
- **When** the operator reconciles the `App`.
- **Then** a `Service` of type `ClusterIP` should be created targeting port `8080`.

### Requirement: Reconcile Ingress
The operator must create an `Ingress` resource to expose the application to the internet if domains are specified.

#### Scenario: Expose app via ingress with TLS
- **Given** an `App` with `domains: ["example.com"]` and `tls: true`.
- **When** the operator reconciles the `App`.
- **Then** an `Ingress` should be created with the domain host, TLS configuration, and cert-manager annotations.

### Requirement: Resource Ownership
All created Kubernetes resources must be owned by the `App` instance.

#### Scenario: Ensure garbage collection
- **Given** an `App` that has created a `Deployment` and `Service`.
- **When** the `App` is deleted.
- **Then** the `Deployment` and `Service` should be automatically removed by Kubernetes garbage collection.
