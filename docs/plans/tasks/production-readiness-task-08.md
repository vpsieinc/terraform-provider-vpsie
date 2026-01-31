# Task: Phase 4e - Descriptions and Validators for backup, snapshot Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for backup and snapshot services
- Size: Medium (8 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the backup and snapshot service files (both resources and data sources).

## Target Files
- [x] `internal/services/backup/backup_resource.go` (descriptions + validators)
- [x] `internal/services/backup/backup_policy_resource.go` (descriptions + validators)
- [x] `internal/services/backup/backup_data_source.go` (descriptions only)
- [x] `internal/services/backup/backup_policy_datasource.go` (descriptions only)
- [x] `internal/services/snapshot/server_snapshot_resource.go` (descriptions + validators)
- [x] `internal/services/snapshot/snapshot_policy_resource.go` (descriptions + validators)
- [x] `internal/services/snapshot/server_snapshot_data_source.go` (descriptions only)
- [x] `internal/services/snapshot/snapshot_policy_datasource.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [x] Read each target file to catalog all schema attributes
- [x] Identify Required vs Optional vs Computed attributes

### 2. Implementation
For each resource file:
- [x] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [x] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [x] Add appropriate int64 validators for numeric Required fields
- [x] Import validator packages as needed

For each data source file:
- [x] Add `MarkdownDescription` to every attribute (descriptions only)

### 3. Verify Skill Fidelity
- [x] Every attribute in all 8 files has a non-empty MarkdownDescription
- [x] All Required string attributes have LengthAtLeast(1)

## Completion Criteria
- [x] All 8 files have `MarkdownDescription` on every attribute
- [x] All Required string attributes have validators
- [x] `go build -v .` succeeds
- [x] Verification level: L3 (build success)

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
