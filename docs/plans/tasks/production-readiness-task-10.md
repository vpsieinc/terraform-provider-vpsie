# Task: Phase 4g - Descriptions and Validators for project, script, monitoring, accesstoken, image Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for project, script, monitoring, accesstoken, image services
- Size: Medium (10 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the project, script, monitoring, accesstoken, and image service files (both resources and data sources). The monitoring_rule_resource gets a `OneOf` validator for the `condition` field.

## Target Files
- [ ] `internal/services/project/project_resource.go` (descriptions + validators)
- [ ] `internal/services/project/project_data_source.go` (descriptions only)
- [ ] `internal/services/script/script_resource.go` (descriptions + validators)
- [ ] `internal/services/script/script_data_source.go` (descriptions only)
- [ ] `internal/services/monitoring/monitoring_rule_resource.go` (descriptions + validators, including OneOf for condition)
- [ ] `internal/services/monitoring/monitoring_rule_datasource.go` (descriptions only)
- [ ] `internal/services/accesstoken/accesstoken_resource.go` (descriptions + validators)
- [ ] `internal/services/accesstoken/accesstoken_datasource.go` (descriptions only)
- [ ] `internal/services/image/image_resource.go` (descriptions + validators)
- [ ] `internal/services/image/image_data_source.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each target file to catalog all schema attributes
- [ ] Identify Required vs Optional vs Computed attributes
- [ ] Confirm monitoring condition enum values from API documentation

### 2. Implementation
For each resource file:
- [ ] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [ ] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [ ] Add `stringvalidator.OneOf(...)` to monitoring `condition` attribute (values from API docs)
- [ ] Add appropriate int64 validators for numeric Required fields
- [ ] Import validator packages as needed

For each data source file:
- [ ] Add `MarkdownDescription` to every attribute (descriptions only)

### 3. Verify Skill Fidelity
- [ ] Every attribute in all 10 files has a non-empty MarkdownDescription
- [ ] All Required string attributes have LengthAtLeast(1)
- [ ] monitoring condition has OneOf validator

## Completion Criteria
- [ ] All 10 files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have validators
- [ ] monitoring `condition` has `OneOf` validator
- [ ] `go build -v .` succeeds
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Spot-check descriptions
grep -c "MarkdownDescription" internal/services/project/project_resource.go
grep -c "MarkdownDescription" internal/services/script/script_resource.go
grep -c "MarkdownDescription" internal/services/monitoring/monitoring_rule_resource.go
grep -c "MarkdownDescription" internal/services/accesstoken/accesstoken_resource.go
grep -c "MarkdownDescription" internal/services/image/image_resource.go

# 2. Verify OneOf validator
grep "OneOf" internal/services/monitoring/monitoring_rule_resource.go

# 3. Build
go build -v .
```

## Notes
- Impact scope: Schema metadata only; no CRUD logic changes
- Constraints: Use only permissive validators; OneOf values must match API-accepted values
- AC Coverage: AC-4.1 (partial), AC-4.3 (partial), AC-4.4 (partial -- OneOf for condition)
