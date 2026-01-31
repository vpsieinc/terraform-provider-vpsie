# Phase 6 Completion: ImportState + Final Quality Assurance

## Phase Summary
Phase 6 adds `ImportState` to the 7 remaining resources, implements integration and E2E test skeletons, and performs the final quality gate verification across all 26 acceptance criteria.

## Task Completion Checklist
- [ ] Task 18: ImportState implementations + integration/E2E tests + final QA

## E2E Verification Procedures (from Design Doc)

### 1. Verify ImportState interface assertions compile
```bash
go build -v .
```
Expected: Build succeeds

### 2. Count ImportState implementations
```bash
grep -rn "ResourceWithImportState" internal/services/ | wc -l
```
Expected: 27 (20 existing + 7 new)

### 3. Verify composite ID error messages
```bash
grep -rn "Expected import identifier with format" internal/services/
```
Expected: 4 results (one per composite ID resource)

### 4. Run full test suite
```bash
go test ./... -count=1
```
Expected: All tests pass

### 5. Full build and module verification
```bash
go build -v .
go mod verify
go generate ./...
```
Expected: All succeed with zero errors

### 6. Final acceptance criteria checklist
```bash
# AC-1.x: Sensitive fields
grep -rn "Sensitive:\s*true" internal/provider/ internal/services/sshkey/ internal/services/server/ | wc -l
# Expected: 5 new + pre-existing = 11+

# AC-2.x: No replace directives
grep "replace" go.mod
# Expected: zero results

# AC-3.x: No discarded errors
grep -rn '_, _ :=' internal/services/firewall/
# Expected: zero results

# AC-4.x: Descriptions
grep -rL "MarkdownDescription" internal/services/*/*.go | grep -v _test | grep -v _api
# Expected: zero results (all files have descriptions)

# AC-5.x: Test files
find internal/services -name "*_test.go" | wc -l
# Expected: 20+ test files

# AC-6.x: ImportState
grep -rn "ResourceWithImportState" internal/services/ | wc -l
# Expected: 27
```

## Phase Completion Criteria
- [ ] 7 resources have `ImportState` implementations
- [ ] 27/27 resources now support `terraform import`
- [ ] `TestAccImportStateCompositeID` integration test implemented and compiles
- [ ] `TestAccImportStateCompositeID_InvalidFormat` integration test implemented and compiles
- [ ] `TestE2EProductionReadinessFullCycle` E2E test implemented and compiles
- [ ] `go build -v .` succeeds
- [ ] `go test ./... -count=1` passes (all unit tests)
- [ ] All 26 Design Doc acceptance criteria verified
- [ ] Independently committable

## Overall Project Completion Criteria
- [ ] All 6 phases completed in strict order
- [ ] Each phase's operational verification procedures executed and passed
- [ ] All 26 Design Doc acceptance criteria satisfied
- [ ] All 5 test skeletons (4 integration + 1 E2E) implemented
- [ ] `go build -v .` succeeds
- [ ] `go mod verify` reports all modules verified
- [ ] `go test ./... -count=1` passes (all unit tests)
- [ ] `go generate ./...` succeeds (docs regenerated)
- [ ] Zero `replace` directives in `go.mod`
- [ ] Zero `_, _ :=` patterns in firewall code
- [ ] 27/27 resources support `terraform import`
- [ ] 20/20 service modules have unit tests
- [ ] 6/6 existing acceptance tests have `CheckDestroy`
- [ ] 100% schema attribute description coverage
