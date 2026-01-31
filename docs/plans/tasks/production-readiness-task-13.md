# Task: Phase 5b - Interface Wrappers + Unit Tests for storage, server, sshkey Services

Metadata:
- Dependencies: task-12 (CheckDestroy must be complete for existing test files in these services)
- Provides: Interface wrappers + unit tests for storage, server, sshkey services
- Size: Medium (9 new files: 3 api.go + 3 resource tests + 3 data source adjustments)

## Implementation Content
Create service-specific interface wrappers (`_api.go` files) for storage, server, and sshkey services. Each interface declares only the methods used by that service's resource and data source. Update resource/data source struct fields from `*govpsie.Client` to the narrow interface type. Create unit tests with mock implementations.

Note: sshkey and storage already have acceptance tests. This task adds the interface layer and NEW unit tests (not acceptance tests).

## Target Files
- [ ] `internal/services/storage/storage_api.go` (NEW: StorageAPI interface)
- [ ] `internal/services/server/server_api.go` (NEW: ServerAPI interface)
- [ ] `internal/services/sshkey/sshkey_api.go` (NEW: SshkeyAPI interface)
- [ ] `internal/services/storage/storage_resource.go` (update struct field type + Configure)
- [ ] `internal/services/storage/storage_data_source.go` (update struct field type + Configure)
- [ ] `internal/services/server/server_resource.go` (update struct field type + Configure)
- [ ] `internal/services/server/server_data_source.go` (update struct field type + Configure)
- [ ] `internal/services/sshkey/sshkey_resource.go` (update struct field type + Configure)
- [ ] `internal/services/sshkey/sshkey_data_source.go` (update struct field type + Configure)

For unit tests (only for server, as sshkey and storage already have acceptance tests):
- [ ] `internal/services/server/server_resource_test.go` (NEW: mock ServerAPI, test Create/Read/Delete)

## Implementation Steps
### 1. Confirm Skill Constraints (Red Phase)
- [ ] Read each resource and data source file to catalog all SDK method calls
- [ ] Identify the sub-service field on `govpsie.Client` used by each service
- [ ] Verify the interface wrapper pattern from Design Doc

### 2. Green Phase - Interface Wrappers
For each service (storage, server, sshkey):
- [ ] Create `<service>_api.go` defining a narrow interface with only the methods called by that service
- [ ] Update resource struct: change `client *govpsie.Client` to `client <ServiceAPI>`
- [ ] Update data source struct: same change
- [ ] Update Configure methods: extract sub-service from full client
  ```go
  // Example pattern:
  type StorageAPI interface {
      ListStorages(ctx, opt) ([]govpsie.Storage, error)
      CreateStorage(ctx, req) error
      DeleteStorage(ctx, id) error
      // ... only methods actually used
  }
  ```

### 3. Green Phase - Unit Tests
- [ ] Create `server_resource_test.go` with mock ServerAPI implementation
- [ ] Write table-driven tests for Create, Read, Delete operations
- [ ] Verify mock returns expected values and resource state is populated correctly

### 4. Refactor Phase
- [ ] Ensure all existing acceptance tests still compile (sshkey, storage)
- [ ] Confirm interfaces are narrow (only methods used)
- [ ] Run `go build -v .` to verify all concrete types satisfy interfaces

## Completion Criteria
- [ ] 3 `_api.go` interface files created
- [ ] Resource and data source structs use interface types
- [ ] `go build -v .` succeeds (compile-time interface satisfaction proof)
- [ ] `go test ./internal/services/server/ -count=1` passes (unit tests)
- [ ] Existing acceptance tests in storage and sshkey still compile
- [ ] Verification level: L2 (unit tests pass) + L3 (build success)

## Verification Steps
```bash
# 1. Verify interface files exist
ls internal/services/storage/storage_api.go
ls internal/services/server/server_api.go
ls internal/services/sshkey/sshkey_api.go

# 2. Build verification (proves interfaces are satisfied)
go build -v .

# 3. Run unit tests
go test ./internal/services/server/ -count=1
go test ./internal/services/storage/ -count=1
go test ./internal/services/sshkey/ -count=1
```

## Notes
- Impact scope: Struct field types change from concrete to interface; Configure methods change extraction pattern
- Constraints: Interfaces must wrap sub-service fields (e.g., `client.Storage`), NOT the full `*govpsie.Client`
- Constraints: Do not modify existing acceptance test logic
- AC Coverage: AC-5.1 (partial), AC-5.3 (partial)
