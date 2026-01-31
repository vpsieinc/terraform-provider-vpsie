# Task: Phase 6 - ImportState Implementations + Integration/E2E Tests + Final QA

Metadata:
- Dependencies: Phase 5 completion (all interface wrappers and tests must be done)
- Provides: ImportState for 7 remaining resources; integration/E2E tests; final quality gate
- Size: Medium (7 resource files + 2 test files)

## Implementation Content
Add `ImportState` to the 7 resources that lack it: 3 with standard passthrough IDs and 4 with composite `/`-separated IDs. Implement integration test skeletons for composite ID import. Implement the E2E test skeleton. Perform final quality gate verification.

## Target Files

### Standard ImportState (3 resources)
- [ ] `internal/services/accesstoken/accesstoken_resource.go` (ImportState with `path.Root("identifier")`)
- [ ] `internal/services/fip/fip_resource.go` (ImportState with `path.Root("id")`)
- [ ] `internal/services/monitoring/monitoring_rule_resource.go` (ImportState with `path.Root("identifier")`)

### Composite ID ImportState (4 resources)
- [ ] `internal/services/firewall/firewall_attachment_resource.go` (format: `<group_id>/<vm_identifier>`)
- [ ] `internal/services/vpc/vpc_server_assignment_resource.go` (format: `<vm_identifier>/<vpc_id>`, vpc_id needs int64 parse)
- [ ] `internal/services/domain/reverse_dns_resource.go` (format: `<vm_identifier>/<ip>`)
- [ ] `internal/services/domain/dns_record_resource.go` (format: `<domain_identifier>/<type>/<name>`)

### Integration + E2E Tests
- [ ] `tests/integration_test.go` (implement TestAccImportStateCompositeID + InvalidFormat)
- [ ] `tests/e2e_test.go` (implement TestE2EProductionReadinessFullCycle)

## Implementation Steps

### 1. Confirm Skill Constraints
- [ ] Verify import ID fields and struct names match Design Doc for all 7 resources
- [ ] Read each resource file to confirm the primary ID attribute name

### 2. Standard ImportState Implementation
For each of the 3 standard resources:
- [ ] Add `resource.ResourceWithImportState` interface assertion: `var _ resource.ResourceWithImportState = &<Resource>{}`
- [ ] Add `ImportState` method using `resource.ImportStatePassthroughID` with the appropriate `path.Root()`

### 3. Composite ID ImportState Implementation
For each of the 4 composite resources:
- [ ] Add `resource.ResourceWithImportState` interface assertion
- [ ] Implement custom `ImportState` method:
  1. Split `req.ID` on `/`
  2. Validate correct number of parts (2 for firewall_attachment, vpc_server_assignment, reverse_dns; 3 for dns_record)
  3. Validate no empty segments
  4. Set each attribute from parsed parts using `resp.Diagnostics.Append(resp.State.SetAttribute(...)...)`
  5. For vpc_server_assignment: parse vpc_id as int64
  6. On validation failure: add error diagnostic with format hint
    ```
    "Expected import identifier with format: <format_string>"
    ```

### 4. Integration Tests
- [ ] Implement `TestAccImportStateCompositeID`:
  - Create DNS record resource
  - Import with 3-part composite ID
  - Verify state attributes match
- [ ] Implement `TestAccImportStateCompositeID_InvalidFormat`:
  - Attempt import with malformed ID
  - Verify ExpectError matches format hint
- [ ] Complete remaining TODO items in `testAccDnsRecordConfig` and `testAccImportStateIDFunc`

### 5. E2E Test
- [ ] Implement `TestE2EProductionReadinessFullCycle`:
  - Complete `testE2EProductionReadinessConfig_phase1()` (already partially done)
  - Complete `testE2EProductionReadinessConfig_phase3()`
  - Verify all 4 test steps work together

### 6. Final Quality Gate
- [ ] Verify all 26 Design Doc acceptance criteria:
  - AC-1.x: 5/5 sensitive fields (grep verification)
  - AC-2.x: Zero replace directives; build and module verify pass
  - AC-3.x: Zero `_, _ :=` in firewall code
  - AC-4.x: 100% description coverage; validators on all Required attributes
  - AC-5.x: 20/20 services with unit tests; 6/6 acceptance tests with CheckDestroy
  - AC-6.x: 27/27 resources support ImportState
- [ ] Run final quality checks:
  ```bash
  go build -v .
  go mod verify
  go test ./... -count=1
  go generate ./...
  ```

### 7. Verify Skill Fidelity
- [ ] Final review of all changes against Design Doc specifications
- [ ] Confirm no divergence from ADR decisions

## Completion Criteria
- [ ] 7 resources have `ImportState` implementations
- [ ] 27/27 resources now support `terraform import`
- [ ] Integration tests implemented and compile
- [ ] E2E test implemented and compiles
- [ ] `go build -v .` succeeds
- [ ] `go test ./... -count=1` passes
- [ ] All 26 Design Doc acceptance criteria verified
- [ ] Verification level: L1 (full lifecycle works) + L2 (all tests pass) + L3 (build success)

## Verification Steps
```bash
# 1. Build verification
go build -v .

# 2. Count ImportState implementations
grep -rn "ResourceWithImportState" internal/services/ | wc -l
# Expected: 27 (20 existing + 7 new)

# 3. Verify composite ID error messages
grep -rn "Expected import identifier with format" internal/services/
# Expected: 4 results

# 4. Run full test suite
go test ./... -count=1

# 5. Full build and module verification
go mod verify
go generate ./...

# 6. Final AC checklist
grep -rn "Sensitive:\s*true" internal/provider/ internal/services/sshkey/ internal/services/server/ | wc -l
grep "replace" go.mod
grep -rn '_, _ :=' internal/services/firewall/
grep -rL "MarkdownDescription" internal/services/*/*.go | grep -v _test | grep -v _api
find internal/services -name "*_test.go" | wc -l
grep -rn "ResourceWithImportState" internal/services/ | wc -l
```

## Notes
- Impact scope: 7 resource files (ImportState method added) + 2 test files
- Constraints: Standard ImportState uses passthrough; composite uses `/` separator per ADR Decision 5
- Composite ID error messages must include the expected format string
- For vpc_server_assignment, vpc_id requires `strconv.ParseInt` conversion
- AC Coverage: AC-6.1 through AC-6.8 (all Priority 6 criteria)
