# Task: Phase 5c - Interface Wrappers + Unit Tests for kubernetes, loadbalancer, firewall Services

Metadata:
- Dependencies: task-13 (establishes interface wrapper pattern)
- Provides: Interface wrappers + unit tests for kubernetes, loadbalancer, firewall services
- Size: Medium (new _api.go files + new _test.go files + existing file updates)

## Implementation Content
Create service-specific interface wrappers for kubernetes, loadbalancer, and firewall services. Create unit tests with mock implementations for all three (none currently have tests).

## Target Files
- [ ] `internal/services/kubernetes/kubernetes_api.go` (NEW: KubernetesAPI interface)
- [ ] `internal/services/loadbalancer/loadbalancer_api.go` (NEW: LoadbalancerAPI interface)
- [ ] `internal/services/firewall/firewall_api.go` (NEW: FirewallAPI interface)
- [ ] `internal/services/kubernetes/kubernetes_resource.go` (update struct field + Configure)
- [ ] `internal/services/kubernetes/kubernetes_group_resource.go` (update struct field + Configure)
- [ ] `internal/services/kubernetes/kubernetes_data_source.go` (update struct field + Configure)
- [ ] `internal/services/kubernetes/kubernetes_group_data_source.go` (update struct field + Configure)
- [ ] `internal/services/loadbalancer/loadbalancer_resource.go` (update struct field + Configure)
- [ ] `internal/services/loadbalancer/loadbalancer_data_source.go` (update struct field + Configure)
- [ ] `internal/services/firewall/firewall_resource.go` (update struct field + Configure)
- [ ] `internal/services/firewall/firewall_attachment_resource.go` (update struct field + Configure)
- [ ] `internal/services/firewall/firewall_data_source.go` (update struct field + Configure)
- [ ] `internal/services/kubernetes/kubernetes_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/loadbalancer/loadbalancer_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/firewall/firewall_resource_test.go` (NEW: unit tests with error propagation verification)

## Implementation Steps
### 1. Red Phase
- [ ] Read each resource/data source to catalog all SDK method calls
- [ ] Write failing unit tests for Create/Read/Delete operations

### 2. Green Phase
- [ ] Create interface files with narrow method sets
- [ ] Update struct fields and Configure methods
- [ ] Create mock implementations in test files
- [ ] Implement tests to pass

### 3. Refactor Phase
- [ ] Verify interfaces are minimal
- [ ] Firewall unit test should include test for ListValueFrom error propagation (AC-3.1 unit-level)
- [ ] Build verification

## Completion Criteria
- [ ] 3 `_api.go` interface files created
- [ ] All struct fields use interface types
- [ ] 3 new unit test files created
- [ ] `go build -v .` succeeds
- [ ] `go test ./internal/services/kubernetes/ -count=1` passes
- [ ] `go test ./internal/services/loadbalancer/ -count=1` passes
- [ ] `go test ./internal/services/firewall/ -count=1` passes
- [ ] Verification level: L2 (unit tests pass) + L3 (build success)

## Verification Steps
```bash
# 1. Verify interface files exist
ls internal/services/kubernetes/kubernetes_api.go
ls internal/services/loadbalancer/loadbalancer_api.go
ls internal/services/firewall/firewall_api.go

# 2. Build verification
go build -v .

# 3. Run unit tests
go test ./internal/services/kubernetes/ -count=1
go test ./internal/services/loadbalancer/ -count=1
go test ./internal/services/firewall/ -count=1
```

## Notes
- Impact scope: Struct field types + new files; no CRUD behavior changes
- Constraints: Interfaces wrap sub-service fields only
- Firewall test should verify error propagation from ListValueFrom (AC-3.1)
- AC Coverage: AC-5.1 (partial), AC-5.3 (partial)
