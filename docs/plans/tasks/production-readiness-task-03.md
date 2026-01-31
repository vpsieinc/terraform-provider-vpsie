# Task: Phase 3 - Fix Firewall Error Suppression

Metadata:
- Dependencies: task-02 (Phase 2 must be complete)
- Provides: Proper diagnostic propagation for all 12 ListValueFrom calls in firewall code
- Size: Small (2 implementation files + 1 integration test)

## Implementation Content
Replace all 12 silently discarded `ListValueFrom` errors in firewall code with proper diagnostic propagation. The pattern `_, _ :=` becomes `destList, destDiags :=` followed by `resp.Diagnostics.Append(destDiags...)` and an early return guard. Implement the `TestAccFirewallErrorPropagation` integration test skeleton.

## Target Files
- [ ] `internal/services/firewall/firewall_resource.go` (8 ListValueFrom calls: Create lines 394-395, 421-422; Read lines 508-509, 535-536)
- [ ] `internal/services/firewall/firewall_data_source.go` (4 ListValueFrom calls: Read lines 332-333, 359-360)
- [ ] `tests/integration_test.go` (implement TestAccFirewallErrorPropagation)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read both firewall files and verify all 12 `ListValueFrom` call locations match the plan
- [ ] Confirm the current pattern is `_, _ :=` (discarded diagnostics)

### 2. Implementation
- [ ] Fix `firewall_resource.go` Create method (lines 394-395 InBound, lines 421-422 OutBound):
  - Replace `_, _ :=` with diagnostic capture variables
  - Append diagnostics to `resp.Diagnostics`
  - Add `if resp.Diagnostics.HasError() { return }` guard after each pair
- [ ] Fix `firewall_resource.go` Read method (lines 508-509 InBound, lines 535-536 OutBound):
  - Same pattern as Create
- [ ] Fix `firewall_data_source.go` Read method (lines 332-333 InBound, lines 359-360 OutBound):
  - Same pattern as above
- [ ] Implement `TestAccFirewallErrorPropagation` in `tests/integration_test.go`:
  - Complete the TODO items for firewall config with inbound/outbound rules
  - Complete `testAccFirewallWithRulesConfig()` helper
  - Verify CRUD succeeds without silent failures

### 3. Verify Skill Fidelity
- [ ] Confirm zero `_, _ :=` patterns remain in firewall code
- [ ] Each error pair includes `if resp.Diagnostics.HasError() { return }` guard

## Completion Criteria
- [ ] All 12 `ListValueFrom` calls have diagnostic capture and propagation
- [ ] `go build -v .` succeeds
- [ ] `grep -rn '_, _ :=' internal/services/firewall/` returns zero results
- [ ] `grep -rn '_, _\s*=' internal/services/firewall/` returns zero results
- [ ] `TestAccFirewallErrorPropagation` integration test implemented and compiles
- [ ] Verification level: L3 (build success) + L2 (test compiles)

## Verification Steps
```bash
# 1. Verify no discarded errors remain
grep -rn '_, _ :=' internal/services/firewall/
grep -rn '_, _\s*=' internal/services/firewall/
# Expected: Zero results for both

# 2. Verify error propagation pattern
grep -c "resp.Diagnostics.Append" internal/services/firewall/firewall_resource.go
grep -c "resp.Diagnostics.Append" internal/services/firewall/firewall_data_source.go
# Expected: Resource shows increase of 8; data source shows increase of 4

# 3. Build verification
go build -v .

# 4. Verify test compiles
go test ./tests/ -run TestAccFirewallErrorPropagation -count=0
```

## Notes
- Impact scope: Only firewall error handling; no changes to firewall business logic
- Constraints: Do not change the ListValueFrom call arguments or return types
- AC Coverage: AC-3.1, AC-3.2, AC-3.3, AC-3.4
- The error propagation pattern per pair:
  ```go
  destList, destDiags := types.ListValueFrom(...)
  resp.Diagnostics.Append(destDiags...)
  sourceList, sourceDiags := types.ListValueFrom(...)
  resp.Diagnostics.Append(sourceDiags...)
  if resp.Diagnostics.HasError() {
      return
  }
  ```
