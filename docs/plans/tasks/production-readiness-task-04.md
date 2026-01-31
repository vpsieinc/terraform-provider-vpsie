# Task: Phase 4a - Add Validators Dependency + Provider Descriptions

Metadata:
- Dependencies: task-03 (Phase 3 must be complete)
- Provides: Validator library available for all subsequent Phase 4 tasks; provider-level schema documented
- Size: Small (2 files: go.mod, provider.go)

## Implementation Content
Add the `terraform-plugin-framework-validators` dependency to `go.mod` and add `MarkdownDescription` plus `Sensitive` validator awareness to the provider-level schema in `provider.go`. This establishes the foundation for all subsequent Phase 4 description/validator tasks.

## Target Files
- [ ] `go.mod` (add terraform-plugin-framework-validators dependency)
- [ ] `internal/provider/provider.go` (add MarkdownDescription to provider schema and access_token attribute)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Verify validator library compatibility with terraform-plugin-framework v1.17.0
- [ ] Read `internal/provider/provider.go` to understand current schema structure

### 2. Implementation
- [ ] Add validator dependency:
  ```bash
  go get github.com/hashicorp/terraform-plugin-framework-validators@v0.17.0
  go mod tidy
  ```
- [ ] Add `MarkdownDescription` to provider-level schema in `provider.go`:
  - Provider schema description (what the VPSie provider does)
  - `access_token` attribute description (API authentication token)
- [ ] Verify build succeeds with new dependency

### 3. Verify Skill Fidelity
- [ ] Validator library version is compatible
- [ ] Provider schema has MarkdownDescription

## Completion Criteria
- [ ] `terraform-plugin-framework-validators` present in `go.mod`
- [ ] `go build -v .` succeeds
- [ ] Provider schema has `MarkdownDescription` on the provider and `access_token` attribute
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Verify validator dependency
grep "terraform-plugin-framework-validators" go.mod

# 2. Verify provider description
grep "MarkdownDescription" internal/provider/provider.go

# 3. Build verification
go build -v .
```

## Notes
- Impact scope: go.mod dependency + provider.go schema metadata only
- Constraints: Do not add validators to provider.go yet (access_token is Optional, no validator needed)
- AC Coverage: AC-4.1 (partial), AC-4.2 (partial)
- This task MUST complete before tasks 05-11 can begin
