# Proposal: Implement Basic Functional App CRD

## Summary
Define the initial functional `App` Custom Resource Definition (CRD) to support basic application management features. This includes container configuration, networking (domains/TLS), scaling (replicas), and resource constraints.

## Problem Statement
The current `App` CRD is a scaffold with no functional fields. To make Portra useful as a PaaS, it must at least be able to define what image to run, how to expose it, and how to scale it.

## Proposed Changes
1.  **CRD Schema Update:** Add fields for `image`, `containerPort`, `domains`, `tls`, `env`, `replicas`, and `resources` to the `AppSpec`.
2.  **Controller Implementation:** Update the `AppReconciler` to create/update Kubernetes `Deployment`, `Service`, and `Ingress` (or Gateway API) based on the `App` spec.
3.  **Status Reporting:** Update `AppStatus` to reflect the health of the underlying resources.
4.  **Testing:** Add Ginkgo/Gomega tests to verify the reconciliation logic.

## Dependencies
- Kubernetes cluster with an Ingress Controller (e.g., NGINX or Traefik) if we use Ingress resources.
- Cert-manager if we want to automate TLS (to be decided in design).

## Alternatives Considered
- Using `Gateway API` instead of `Ingress`: While more modern, `Ingress` is more widely available in simple setups. We'll stick to a basic implementation first.
