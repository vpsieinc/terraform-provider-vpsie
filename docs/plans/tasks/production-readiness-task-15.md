# Task: Phase 5d - Interface Wrappers + Unit Tests for domain, backup, snapshot Services

Metadata:
- Dependencies: task-13 (establishes interface wrapper pattern)
- Provides: Interface wrappers + unit tests for domain, backup, snapshot services
- Size: Medium (new _api.go files + new _test.go files + existing file updates)

## Implementation Content
Create service-specific interface wrappers for domain, backup, and snapshot services. Create unit tests with mock implementations. Note: domain already has acceptance tests but no interface wrapper or unit tests.

## Target Files
- [ ] `internal/services/domain/domain_api.go` (NEW: DomainAPI interface)
- [ ] `internal/services/backup/backup_api.go` (NEW: BackupAPI interface)
- [ ] `internal/services/snapshot/snapshot_api.go` (NEW: SnapshotAPI interface)
- [ ] `internal/services/domain/domain_resource.go` (update struct field + Configure)
- [ ] `internal/services/domain/dns_record_resource.go` (update struct field + Configure)
- [ ] `internal/services/domain/reverse_dns_resource.go` (update struct field + Configure)
- [ ] `internal/services/domain/domain_data_source.go` (update struct field + Configure)
- [ ] `internal/services/backup/backup_resource.go` (update struct field + Configure)
- [ ] `internal/services/backup/backup_policy_resource.go` (update struct field + Configure)
- [ ] `internal/services/backup/backup_data_source.go` (update struct field + Configure)
- [ ] `internal/services/backup/backup_policy_datasource.go` (update struct field + Configure)
- [ ] `internal/services/snapshot/server_snapshot_resource.go` (update struct field + Configure)
- [ ] `internal/services/snapshot/snapshot_policy_resource.go` (update struct field + Configure)
- [ ] `internal/services/snapshot/server_snapshot_data_source.go` (update struct field + Configure)
- [ ] `internal/services/snapshot/snapshot_policy_datasource.go` (update struct field + Configure)
- [ ] `internal/services/backup/backup_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/snapshot/server_snapshot_resource_test.go` (NEW: unit tests)

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
- [ ] Confirm existing domain acceptance tests still compile
- [ ] Build verification

## Completion Criteria
- [ ] 3 `_api.go` interface files created
- [ ] All struct fields use interface types
- [ ] 2 new unit test files created (backup, snapshot)
- [ ] `go build -v .` succeeds
- [ ] `go test ./internal/services/domain/ -count=1` passes
- [ ] `go test ./internal/services/backup/ -count=1` passes
- [ ] `go test ./internal/services/snapshot/ -count=1` passes
- [ ] Verification level: L2 (unit tests pass) + L3 (build success)

## Verification Steps
```bash
# 1. Verify interface files exist
ls internal/services/domain/domain_api.go
ls internal/services/backup/backup_api.go
ls internal/services/snapshot/snapshot_api.go

# 2. Build verification
go build -v .

# 3. Run unit tests
go test ./internal/services/domain/ -count=1
go test ./internal/services/backup/ -count=1
go test ./internal/services/snapshot/ -count=1
```

## Notes
- Impact scope: Struct field types + new files
- Constraints: domain already has acceptance tests -- do not break them
- AC Coverage: AC-5.1 (partial), AC-5.3 (partial)
