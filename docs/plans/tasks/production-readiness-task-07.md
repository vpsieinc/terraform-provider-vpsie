# Task: Phase 4d - Descriptions and Validators for domain Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for domain, dns_record, reverse_dns services
- Size: Medium (4 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the domain service files. The `dns_record_resource.go` gets a `OneOf` validator for the record `type` field (A, AAAA, CNAME, MX, TXT, NS, SRV).

## Target Files
- [ ] `internal/services/domain/domain_resource.go` (descriptions + validators)
- [ ] `internal/services/domain/dns_record_resource.go` (descriptions + validators, including OneOf for record type)
- [ ] `internal/services/domain/reverse_dns_resource.go` (descriptions + validators)
- [ ] `internal/services/domain/domain_data_source.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each target file to catalog all schema attributes
- [ ] Identify Required vs Optional vs Computed attributes
- [ ] Confirm DNS record type enum values from API documentation

### 2. Implementation
For each resource file:
- [ ] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [ ] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [ ] Add `stringvalidator.OneOf("A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV")` to dns_record `type` attribute
- [ ] Import validator packages as needed

For the data source file:
- [ ] Add `MarkdownDescription` to every attribute (descriptions only)

### 3. Verify Skill Fidelity
- [ ] Every attribute in all 4 files has a non-empty MarkdownDescription
- [ ] All Required string attributes have LengthAtLeast(1)
- [ ] dns_record type has OneOf validator

## Completion Criteria
- [ ] All 4 files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have validators
- [ ] dns_record `type` has `OneOf` validator
- [ ] `go build -v .` succeeds
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Spot-check descriptions
grep -c "MarkdownDescription" internal/services/domain/domain_resource.go
grep -c "MarkdownDescription" internal/services/domain/dns_record_resource.go
grep -c "MarkdownDescription" internal/services/domain/reverse_dns_resource.go
grep -c "MarkdownDescription" internal/services/domain/domain_data_source.go

# 2. Verify OneOf validator
grep "OneOf" internal/services/domain/dns_record_resource.go

# 3. Build
go build -v .
```

## Notes
- Impact scope: Schema metadata only; no CRUD logic changes
- Constraints: Use only permissive validators; OneOf values must match API-accepted values
- AC Coverage: AC-4.1 (partial), AC-4.3 (partial), AC-4.4 (partial -- OneOf validator)
