# Task: Phase 4f - Descriptions and Validators for vpc, gateway, fip Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for vpc, gateway, fip services
- Size: Medium (7 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the vpc, gateway, and fip service files (both resources and data sources). The fip resource gets a `OneOf` validator for `ip_type` (ipv4, ipv6).

## Target Files
- [x] `internal/services/vpc/vpc_resource.go` (descriptions + validators)
- [x] `internal/services/vpc/vpc_server_assignment_resource.go` (descriptions + validators)
- [x] `internal/services/vpc/vpc_data_source.go` (descriptions only)
- [x] `internal/services/gateway/gateway_resource.go` (descriptions + validators)
- [x] `internal/services/gateway/gateway_data_source.go` (descriptions only)
- [x] `internal/services/fip/fip_resource.go` (descriptions + validators, including OneOf for ip_type)
- [x] `internal/services/fip/fip_datasource.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [x] Read each target file to catalog all schema attributes
- [x] Identify Required vs Optional vs Computed attributes
- [x] Confirm fip ip_type enum values

### 2. Implementation
For each resource file:
- [x] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [x] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [x] Add `stringvalidator.OneOf("ipv4", "ipv6")` to fip `ip_type` attribute
- [x] Add appropriate int64 validators for numeric Required fields
- [x] Import validator packages as needed

For each data source file:
- [x] Add `MarkdownDescription` to every attribute (descriptions only)

### 3. Verify Skill Fidelity
- [x] Every attribute in all 7 files has a non-empty MarkdownDescription
- [x] All Required string attributes have LengthAtLeast(1)
- [x] fip ip_type has OneOf validator

## Completion Criteria
- [x] All 7 files have `MarkdownDescription` on every attribute
- [x] All Required string attributes have validators
- [x] fip `ip_type` has `OneOf` validator
- [x] `go build -v .` succeeds
- [x] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Spot-check descriptions
grep -c "MarkdownDescription" internal/services/vpc/vpc_resource.go
grep -c "MarkdownDescription" internal/services/gateway/gateway_resource.go
grep -c "MarkdownDescription" internal/services/fip/fip_resource.go

# 2. Verify OneOf validator
grep "OneOf" internal/services/fip/fip_resource.go

# 3. Build
go build -v .
```

## Notes
- Impact scope: Schema metadata only; no CRUD logic changes
- Constraints: Use only permissive validators
- AC Coverage: AC-4.1 (partial), AC-4.3 (partial), AC-4.4 (partial -- OneOf for ip_type)
