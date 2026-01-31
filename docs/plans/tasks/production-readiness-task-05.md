# Task: Phase 4b - Descriptions and Validators for server, sshkey, storage Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for server, sshkey, storage services
- Size: Medium (8 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the server, sshkey, and storage service files (both resources and data sources). Storage resource also gets `OneOf` validator for `disk_format` (EXT4, XFS).

## Target Files
- [ ] `internal/services/server/server_resource.go` (descriptions + validators)
- [ ] `internal/services/server/server_data_source.go` (descriptions only)
- [ ] `internal/services/sshkey/sshkey_resource.go` (descriptions + validators)
- [ ] `internal/services/sshkey/sshkey_data_source.go` (descriptions only)
- [ ] `internal/services/storage/storage_resource.go` (descriptions + validators, including OneOf for disk_format)
- [ ] `internal/services/storage/storage_attachment_resource.go` (descriptions + validators)
- [ ] `internal/services/storage/storage_snapshot_resource.go` (descriptions + validators)
- [ ] `internal/services/storage/storage_data_source.go` (descriptions only)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each target file to catalog all schema attributes
- [ ] Identify Required vs Optional vs Computed attributes
- [ ] Identify enum fields needing OneOf validators

### 2. Implementation
For each resource file:
- [ ] Add `MarkdownDescription` to every attribute (1-2 sentence descriptions)
- [ ] Add `stringvalidator.LengthAtLeast(1)` to all Required string attributes
- [ ] Add `stringvalidator.OneOf("EXT4", "XFS")` to storage `disk_format` attribute
- [ ] Add appropriate int64 validators for numeric Required fields (e.g., `int64validator.AtLeast(1)`)
- [ ] Import validator packages as needed

For each data source file:
- [ ] Add `MarkdownDescription` to every attribute (descriptions only, no validators for Computed)

### 3. Verify Skill Fidelity
- [ ] Every attribute in all 8 files has a non-empty MarkdownDescription
- [ ] All Required string attributes have LengthAtLeast(1)
- [ ] storage disk_format has OneOf validator

## Completion Criteria
- [ ] All 8 files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have validators
- [ ] `go build -v .` succeeds
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Spot-check descriptions
grep -c "MarkdownDescription" internal/services/storage/storage_resource.go
grep -c "MarkdownDescription" internal/services/sshkey/sshkey_resource.go
grep -c "MarkdownDescription" internal/services/server/server_data_source.go

# 2. Spot-check validators
grep -c "Validators:" internal/services/storage/storage_resource.go
grep -c "LengthAtLeast" internal/services/sshkey/sshkey_resource.go
grep "OneOf" internal/services/storage/storage_resource.go

# 3. Build
go build -v .
```

## Notes
- Impact scope: Schema metadata only; no CRUD logic changes
- Constraints: Use only permissive validators; do not use restrictive regex patterns
- `storage_snapshot_data_resource.go` is NOT in this task (it is in task-11 with bucket/datacenter/ip)
- AC Coverage: AC-4.1 (partial), AC-4.3 (partial), AC-4.4 (partial)
