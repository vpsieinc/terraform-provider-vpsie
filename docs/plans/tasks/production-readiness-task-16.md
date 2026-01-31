# Task: Phase 5e - Interface Wrappers + Unit Tests for vpc, gateway, fip, monitoring Services

Metadata:
- Dependencies: task-13 (establishes interface wrapper pattern)
- Provides: Interface wrappers + unit tests for vpc, gateway, fip, monitoring services
- Size: Medium (new _api.go files + new _test.go files + existing file updates)

## Implementation Content
Create service-specific interface wrappers for vpc, gateway, fip, and monitoring services. Create unit tests with mock implementations. None of these services currently have tests.

## Target Files
- [ ] `internal/services/vpc/vpc_api.go` (NEW: VpcAPI interface)
- [ ] `internal/services/gateway/gateway_api.go` (NEW: GatewayAPI interface)
- [ ] `internal/services/fip/fip_api.go` (NEW: FipAPI interface)
- [ ] `internal/services/monitoring/monitoring_api.go` (NEW: MonitoringAPI interface)
- [ ] `internal/services/vpc/vpc_resource.go` (update struct field + Configure)
- [ ] `internal/services/vpc/vpc_server_assignment_resource.go` (update struct field + Configure)
- [ ] `internal/services/vpc/vpc_data_source.go` (update struct field + Configure)
- [ ] `internal/services/gateway/gateway_resource.go` (update struct field + Configure)
- [ ] `internal/services/gateway/gateway_data_source.go` (update struct field + Configure)
- [ ] `internal/services/fip/fip_resource.go` (update struct field + Configure)
- [ ] `internal/services/fip/fip_datasource.go` (update struct field + Configure)
- [ ] `internal/services/monitoring/monitoring_rule_resource.go` (update struct field + Configure)
- [ ] `internal/services/monitoring/monitoring_rule_datasource.go` (update struct field + Configure)
- [ ] `internal/services/vpc/vpc_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/gateway/gateway_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/fip/fip_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/monitoring/monitoring_rule_resource_test.go` (NEW: unit tests)

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
- [ ] Build verification

## Completion Criteria
- [ ] 4 `_api.go` interface files created
- [ ] All struct fields use interface types
- [ ] 4 new unit test files created
- [ ] `go build -v .` succeeds
- [ ] `go test ./internal/services/vpc/ -count=1` passes
- [ ] `go test ./internal/services/gateway/ -count=1` passes
- [ ] `go test ./internal/services/fip/ -count=1` passes
- [ ] `go test ./internal/services/monitoring/ -count=1` passes
- [ ] Verification level: L2 (unit tests pass) + L3 (build success)

## Verification Steps
```bash
# 1. Verify interface files exist
ls internal/services/vpc/vpc_api.go
ls internal/services/gateway/gateway_api.go
ls internal/services/fip/fip_api.go
ls internal/services/monitoring/monitoring_api.go

# 2. Build verification
go build -v .

# 3. Run unit tests
go test ./internal/services/vpc/ -count=1
go test ./internal/services/gateway/ -count=1
go test ./internal/services/fip/ -count=1
go test ./internal/services/monitoring/ -count=1
```

## Notes
- Impact scope: Struct field types + new files
- Constraints: Interfaces wrap sub-service fields only
- AC Coverage: AC-5.1 (partial), AC-5.3 (partial)
