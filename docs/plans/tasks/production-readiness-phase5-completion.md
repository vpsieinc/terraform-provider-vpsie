# Phase 5 Completion: Test Coverage

## Phase Summary
Phase 5 creates service-specific interface wrappers, adds unit tests for all 14 untested services, and adds `CheckDestroy` to all 6 existing acceptance test files.

## Task Completion Checklist
- [ ] Task 12: CheckDestroy for 6 existing test files
- [ ] Task 13: Interface wrappers + unit tests for storage, server, sshkey
- [ ] Task 14: Interface wrappers + unit tests for kubernetes, loadbalancer, firewall
- [ ] Task 15: Interface wrappers + unit tests for domain, backup, snapshot
- [ ] Task 16: Interface wrappers + unit tests for vpc, gateway, fip, monitoring
- [ ] Task 17: Interface wrappers + unit tests for project, script, accesstoken, bucket, datacenter, ip, image

## E2E Verification Procedures (from Design Doc)

### 1. Verify interface files exist
```bash
find internal/services -name "*_api.go" | wc -l
```
Expected: 20 (one per service package)

### 2. Verify no compilation errors with interfaces
```bash
go build -v .
```
Expected: Build succeeds (proves sub-service types satisfy interfaces)

### 3. Run unit tests only (no acceptance)
```bash
go test ./... -count=1
```
Expected: All tests pass; output shows test runs in all 20 service packages

### 4. Verify CheckDestroy is present in existing tests
```bash
grep -rn "CheckDestroy" internal/services/storage/storage_resource_test.go
grep -rn "CheckDestroy" internal/services/sshkey/sshkey_resource_test.go
grep -rn "CheckDestroy" internal/services/script/script_resource_test.go
grep -rn "CheckDestroy" internal/services/project/project_resource_test.go
grep -rn "CheckDestroy" internal/services/domain/domain_resource_test.go
grep -rn "CheckDestroy" internal/services/accesstoken/accesstoken_resource_test.go
```
Expected: Each file returns at least one match

### 5. Count test files
```bash
find internal/services -name "*_test.go" | wc -l
```
Expected: 20+ test files

## Phase Completion Criteria
- [ ] 20 `_api.go` interface files created (one per service package)
- [ ] All 20 resource/data source struct types use sub-service interface fields
- [ ] 14 new unit test files created for previously untested services
- [ ] 6 existing acceptance test files have `CheckDestroy` functions
- [ ] `go build -v .` succeeds
- [ ] `go test ./... -count=1` passes (unit tests, no TF_ACC)
- [ ] Each interface declares only methods used by the service
- [ ] Each task group independently committable

## AC Coverage
- AC-5.1: go test ./... passes for all 20 service modules
- AC-5.2: Each existing acceptance test has CheckDestroy
- AC-5.3: Each interface wrapper declares only methods used by the service
