# Task: Phase 4f - Descriptions and Validators for vpc, gateway, fip Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for vpc, gateway, fip services
- Size: Medium (7 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the vpc, gateway, and fip service files (both resources and data sources). The fip resource gets a `OneOf` validator for `ip_type` (ipv4, ipv6).

## Target Files
- [ ] `internal/services/vpc/vpc_resource.go` (descriptions + validators)
- [ ] `internal/services/vpc/vpc_server_assignment_resource.go` (descriptions + validators)
- [ ] `internal/services/vpc/vpc_data_source.go` (descriptions only)
- [ ] `internal/services/gateway/gateway_resource.go` (descriptions + validators)
- [ ] `internal/services/gateway/gateway_data_source.go` (descriptions only)
- [ ] `internal/services/fip/fip_resource.go` (descriptions + validators, including OneOf for ip_type)
- [ ] `internal/services/fip/fip_datasource.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each target file to catalog all schema attributes
- [ ] Identify Required vs Optional vs Computed attributes
- [ ] Confirm fip ip_type enum values

### 2. Implementation
For each resource file:
- [ ] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [ ] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [ ] Add `stringvalidator.OneOf("ipv4", "ipv6")` to fip `ip_type` attribute
- [ ] Add appropriate int64 validators for numeric Required fields
- [ ] Import validator packages as needed

For each data source file:
- [ ] Add `MarkdownDescription` to every attribute (descriptions only)

### 3. Verify Skill Fidelity
- [ ] Every attribute in all 7 files has a non-empty MarkdownDescription
- [ ] All Required string attributes have LengthAtLeast(1)
- [ ] fip ip_type has OneOf validator

## Completion Criteria
- [ ] All 7 files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have validators
- [ ] fip `ip_type` has `OneOf` validator
- [ ] `go build -v .` succeeds
- [ ] Verification level: L3 (build success)

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
