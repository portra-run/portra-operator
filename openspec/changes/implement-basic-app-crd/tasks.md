# Tasks: Implement Basic Functional App CRD

## Task List

- [x] **CRD Definition:** Update `api/v1/app_types.go` with the new fields: `Image`, `ContainerPort`, `Domains`, `TLS`, `Env`, `Replicas`, `Resources`.
- [x] **Code Generation:** Run `make generate` and `make manifests` to update the CRD manifests and generated code.
- [x] **Controller Logic - Deployment:** Implement the reconciliation logic to create or update a Kubernetes `Deployment`.
- [x] **Controller Logic - Service:** Implement the reconciliation logic to create or update a Kubernetes `Service`.
- [x] **Controller Logic - Ingress:** Implement the reconciliation logic to create or update a Kubernetes `Ingress` if `domains` are provided.
- [x] **Controller Logic - Ownership:** Ensure all created resources have an `OwnerReference` pointing to the `App` CR.
- [x] **Controller Logic - Status:** Implement logic to update the `AppStatus` based on the status of its children.
- [x] **Testing - Integration:** Add tests in `internal/controller/app_controller_test.go` to verify the controller's behavior.
- [x] **Validation:** Verify the new `App` CRD by applying a sample YAML to a local Kind cluster (if available).
