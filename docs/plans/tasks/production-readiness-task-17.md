# Task: Phase 5f - Interface Wrappers + Unit Tests for project, script, accesstoken, bucket, datacenter, ip, image Services

Metadata:
- Dependencies: task-13 (establishes interface wrapper pattern)
- Provides: Interface wrappers + unit tests for remaining 7 services; completes Phase 5
- Size: Medium-Large (7 services: new _api.go files + new _test.go files + existing file updates)

## Implementation Content
Create service-specific interface wrappers for the remaining 7 services: project, script, accesstoken, bucket, datacenter, ip, image. Create unit tests with mock implementations for services that do not yet have tests. Note: project, script, and accesstoken already have acceptance tests but need interface wrappers.

## Target Files

### Interface wrappers (7 new files)
- [ ] `internal/services/project/project_api.go` (NEW: ProjectAPI interface)
- [ ] `internal/services/script/script_api.go` (NEW: ScriptAPI interface)
- [ ] `internal/services/accesstoken/accesstoken_api.go` (NEW: AccesstokenAPI interface)
- [ ] `internal/services/bucket/bucket_api.go` (NEW: BucketAPI interface)
- [ ] `internal/services/datacenter/datacenter_api.go` (NEW: DatacenterAPI interface)
- [ ] `internal/services/ip/ip_api.go` (NEW: IpAPI interface)
- [ ] `internal/services/image/image_api.go` (NEW: ImageAPI interface)

### Existing files to update (struct fields + Configure)
- [ ] `internal/services/project/project_resource.go`
- [ ] `internal/services/project/project_data_source.go`
- [ ] `internal/services/script/script_resource.go`
- [ ] `internal/services/script/script_data_source.go`
- [ ] `internal/services/accesstoken/accesstoken_resource.go`
- [ ] `internal/services/accesstoken/accesstoken_datasource.go`
- [ ] `internal/services/bucket/bucket_resource.go`
- [ ] `internal/services/bucket/bucket_datasource.go`
- [ ] `internal/services/datacenter/datacenter_datasource.go`
- [ ] `internal/services/ip/ip_datasource.go`
- [ ] `internal/services/image/image_resource.go`
- [ ] `internal/services/image/image_data_source.go`

### New unit test files (4 services without tests)
- [ ] `internal/services/bucket/bucket_resource_test.go` (NEW: unit tests)
- [ ] `internal/services/datacenter/datacenter_datasource_test.go` (NEW: unit tests)
- [ ] `internal/services/ip/ip_datasource_test.go` (NEW: unit tests)
- [ ] `internal/services/image/image_resource_test.go` (NEW: unit tests)

## Implementation Steps
### 1. Red Phase
- [ ] Read each resource/data source to catalog all SDK method calls
- [ ] Write failing unit tests for untested services

### 2. Green Phase
- [ ] Create interface files with narrow method sets
- [ ] Update struct fields and Configure methods
- [ ] Create mock implementations in test files
- [ ] Implement tests to pass

### 3. Refactor Phase
- [ ] Verify interfaces are minimal
- [ ] Confirm existing acceptance tests still compile (project, script, accesstoken)
- [ ] Build verification

### 4. Full Suite Verification
- [ ] Run `go test ./... -count=1` to verify ALL 20 services pass
- [ ] This is the final Phase 5 implementation task

## Completion Criteria
- [ ] 7 `_api.go` interface files created
- [ ] All struct fields use interface types
- [ ] 4 new unit test files created
- [ ] `go build -v .` succeeds
- [ ] `go test ./... -count=1` passes (ALL services, proving 20/20 coverage)
- [ ] Existing acceptance tests in project, script, accesstoken still compile
- [ ] Verification level: L2 (all unit tests pass) + L3 (build success)

## Verification Steps
```bash
# 1. Verify all 20 interface files exist
find internal/services -name "*_api.go" | wc -l
# Expected: 20 (or appropriately grouped)

# 2. Build verification
go build -v .

# 3. Run ALL unit tests
go test ./... -count=1
# Expected: All pass; output shows test runs in all 20 service packages

# 4. Verify test files exist for previously untested services
ls internal/services/bucket/bucket_resource_test.go
ls internal/services/datacenter/datacenter_datasource_test.go
ls internal/services/ip/ip_datasource_test.go
ls internal/services/image/image_resource_test.go
```

## Notes
- Impact scope: Struct field types + new files across 7 services
- Constraints: Do not break existing acceptance tests for project, script, accesstoken
- This task completes Phase 5 implementation (verification follows in phase completion check)
- AC Coverage: AC-5.1 (complete), AC-5.3 (complete)
