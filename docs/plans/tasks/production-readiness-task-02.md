# Task: Phase 2 - Remove SDK Replace Directive

Metadata:
- Dependencies: task-01 (Phase 1 must be complete)
- Provides: Clean `go.mod` buildable from any machine
- Size: Small (1 file: go.mod, go.sum auto-updated)

## Implementation Content
Remove the local filesystem `replace` directive from `go.mod` (line 16: `replace github.com/vpsie/govpsie => /Users/zozo/projects/govpsie`) and run `go mod tidy` to reconcile checksums with the published SDK.

## Target Files
- [ ] `go.mod` (remove replace directive)
- [ ] `go.sum` (auto-updated by go mod tidy)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Verify `go.mod` line 16 contains the replace directive
- [ ] Confirm published SDK version `v0.0.0-20241020152435-33a7b18a901e` is present in the require block

### 2. Implementation
- [ ] Delete line 16 from `go.mod`: `replace github.com/vpsie/govpsie => /Users/zozo/projects/govpsie`
- [ ] Run `go mod tidy` to update `go.sum` with published module checksums
- [ ] Run `go mod verify` to confirm all modules verified

### 3. Verify Skill Fidelity
- [ ] Confirm `go.mod` has zero replace directives
- [ ] Confirm build uses published module (not local path)

## Completion Criteria
- [ ] `go.mod` contains zero `replace` directives
- [ ] `go mod tidy` completes without errors
- [ ] `go build -v .` succeeds (fetches from module proxy)
- [ ] `go mod verify` reports "all modules verified"
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Verify replace directive removed
grep -n "replace" go.mod
# Expected: Zero results

# 2. Module tidy and verify
go mod tidy
go mod verify
# Expected: "all modules verified"

# 3. Clean build
go build -v .
# Expected: Build succeeds
```

## Notes
- Impact scope: Build system only; no code changes
- Constraints: If published SDK fails, re-add replace temporarily and file SDK issue
- AC Coverage: AC-2.1, AC-2.2, AC-2.3
- Risk: Published SDK may have transitive dependency differences; `go mod tidy` resolves this
