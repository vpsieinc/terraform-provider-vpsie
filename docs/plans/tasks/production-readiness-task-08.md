# Task: Phase 4e - Descriptions and Validators for backup, snapshot Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for backup and snapshot services
- Size: Medium (8 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the backup and snapshot service files (both resources and data sources).

## Target Files
- [ ] `internal/services/backup/backup_resource.go` (descriptions + validators)
- [ ] `internal/services/backup/backup_policy_resource.go` (descriptions + validators)
- [ ] `internal/services/backup/backup_data_source.go` (descriptions only)
- [ ] `internal/services/backup/backup_policy_datasource.go` (descriptions only)
- [ ] `internal/services/snapshot/server_snapshot_resource.go` (descriptions + validators)
- [ ] `internal/services/snapshot/snapshot_policy_resource.go` (descriptions + validators)
- [ ] `internal/services/snapshot/server_snapshot_data_source.go` (descriptions only)
- [ ] `internal/services/snapshot/snapshot_policy_datasource.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each target file to catalog all schema attributes
- [ ] Identify Required vs Optional vs Computed attributes

### 2. Implementation
For each resource file:
- [ ] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [ ] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [ ] Add appropriate int64 validators for numeric Required fields
- [ ] Import validator packages as needed

For each data source file:
- [ ] Add `MarkdownDescription` to every attribute (descriptions only)

### 3. Verify Skill Fidelity
- [ ] Every attribute in all 8 files has a non-empty MarkdownDescription
- [ ] All Required string attributes have LengthAtLeast(1)

## Completion Criteria
- [ ] All 8 files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have validators
- [ ] `go build -v .` succeeds
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Spot-check descriptions
grep -c "MarkdownDescription" internal/services/backup/backup_resource.go
grep -c "MarkdownDescription" internal/services/backup/backup_policy_resource.go
grep -c "MarkdownDescription" internal/services/snapshot/server_snapshot_resource.go

# 2. Spot-check validators
grep -c "LengthAtLeast" internal/services/backup/backup_resource.go
grep -c "Validators:" internal/services/snapshot/server_snapshot_resource.go

# 3. Build
go build -v .
```

## Notes
- Impact scope: Schema metadata only; no CRUD logic changes
- Constraints: Use only permissive validators
- AC Coverage: AC-4.1 (partial), AC-4.3 (partial)
