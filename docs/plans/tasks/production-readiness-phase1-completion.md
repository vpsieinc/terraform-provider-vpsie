# Phase 1 Completion: Sensitive Field Markings

## Phase Summary
Phase 1 eliminates credential exposure in `terraform plan` output and CI logs for all 5 identified sensitive fields.

## Task Completion Checklist
- [ ] Task 01: Add Sensitive: true to all 5 fields + integration test

## E2E Verification Procedures (from Design Doc)

### 1. Grep verification
```bash
grep -rn "Sensitive:" internal/provider/provider.go internal/services/sshkey/ internal/services/server/
```
Expected: 5 new `Sensitive: true` lines (plus pre-existing in other files)

### 2. Build verification
```bash
go build -v .
```
Expected: Build succeeds with zero errors

### 3. Individual field verification
```bash
grep -A2 "access_token" internal/provider/provider.go | grep Sensitive
grep -A2 "private_key" internal/services/sshkey/sshkey_resource.go | grep Sensitive
grep -A2 "private_key" internal/services/sshkey/sshkey_data_source.go | grep Sensitive
grep -A2 "initial_password" internal/services/server/server_data_source.go | grep Sensitive
grep -A2 "initial_password" internal/services/server/server_resource.go | grep Sensitive
```
Expected: All 5 commands return lines containing `Sensitive`

## Phase Completion Criteria
- [ ] All 5 fields have `Sensitive: true` in their schema definitions
- [ ] `go build -v .` succeeds
- [ ] `TestAccSensitiveFieldMasking` integration test implemented and compiles
- [ ] Zero instances of listed sensitive fields without `Sensitive: true` in target files

## AC Coverage
- AC-1.1, AC-1.2, AC-1.3, AC-1.4, AC-1.5
