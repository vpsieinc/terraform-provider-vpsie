# Task: Phase 1 - Sensitive Field Markings

Metadata:
- Dependencies: None (first task)
- Provides: 5 sensitive field markings in provider and service schemas
- Size: Small (5 files + 1 integration test)

## Implementation Content
Add `Sensitive: true` to 5 credential/secret fields across the provider and service schemas: `access_token` in provider.go, `private_key` in sshkey resource and data source, `initial_password` in server resource and data source. Implement the `TestAccSensitiveFieldMasking` integration test skeleton from `tests/integration_test.go`.

## Target Files
- [x] `internal/provider/provider.go` (add Sensitive to access_token)
- [x] `internal/services/sshkey/sshkey_resource.go` (add Sensitive to private_key)
- [x] `internal/services/sshkey/sshkey_data_source.go` (add Sensitive to private_key)
- [x] `internal/services/server/server_data_source.go` (add Sensitive to initial_password)
- [x] `internal/services/server/server_resource.go` (add Sensitive to initial_password)
- [x] `tests/integration_test.go` (implement TestAccSensitiveFieldMasking)

## Implementation Steps
### 1. Confirm Skill Constraints
- [x] Read Design Doc section on sensitive field implementation pattern
- [x] Verify all 5 target files and line locations match the plan:
  - `provider.go` line 69-72 (access_token)
  - `sshkey_resource.go` line 73-78 (private_key)
  - `sshkey_data_source.go` line 66-67 (private_key)
  - `server_data_source.go` line 174-176 (initial_password)
  - `server_resource.go` line 312-317 (initial_password)

### 2. Implementation
- [x] Add `Sensitive: true` to `access_token` attribute in `internal/provider/provider.go`
- [x] Add `Sensitive: true` to `private_key` attribute in `internal/services/sshkey/sshkey_resource.go`
- [x] Add `Sensitive: true` to `private_key` attribute in `internal/services/sshkey/sshkey_data_source.go`
- [x] Add `Sensitive: true` to `initial_password` attribute in `internal/services/server/server_data_source.go`
- [x] Add `Sensitive: true` to `initial_password` attribute in `internal/services/server/server_resource.go`
- [x] Implement `TestAccSensitiveFieldMasking` in `tests/integration_test.go`:
  - Complete the TODO items for sshkey config with private_key
  - Complete the TODO items for sensitive masking verification
  - Complete `testAccSensitiveFieldConfig()` helper

### 3. Verify Skill Fidelity
- [x] Confirm all 5 fields match Design Doc specification exactly
- [x] No target sensitive fields remain without `Sensitive: true`

## Completion Criteria
- [x] All 5 fields have `Sensitive: true` in their schema definitions
- [x] `go build -v .` succeeds
- [x] `grep -rn "Sensitive:" internal/provider/provider.go internal/services/sshkey/ internal/services/server/` shows 5 new Sensitive markings
- [x] `TestAccSensitiveFieldMasking` integration test implemented and compiles
- [x] Verification level: L3 (build success) + L2 (test compiles)

## Verification Steps
```bash
# 1. Grep for all sensitive markings
grep -rn "Sensitive:" internal/provider/provider.go internal/services/sshkey/ internal/services/server/

# 2. Build verification
go build -v .

# 3. Verify each specific field
grep -A2 "access_token" internal/provider/provider.go | grep Sensitive
grep -A2 "private_key" internal/services/sshkey/sshkey_resource.go | grep Sensitive
grep -A2 "private_key" internal/services/sshkey/sshkey_data_source.go | grep Sensitive
grep -A2 "initial_password" internal/services/server/server_data_source.go | grep Sensitive
grep -A2 "initial_password" internal/services/server/server_resource.go | grep Sensitive

# 4. Verify test compiles
go test ./tests/ -run TestAccSensitiveFieldMasking -count=0
```

## Notes
- Impact scope: Only schema metadata changes; no behavioral changes to CRUD operations
- Constraints: Do not modify any CRUD logic, only schema definitions
- AC Coverage: AC-1.1, AC-1.2, AC-1.3, AC-1.4, AC-1.5
