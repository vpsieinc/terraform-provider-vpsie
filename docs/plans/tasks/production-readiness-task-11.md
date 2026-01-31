# Task: Phase 4h - Descriptions for bucket, datacenter, ip + storage_snapshot_data_resource + Doc Regeneration

Metadata:
- Dependencies: tasks 05-10 (all description/validator tasks must be complete)
- Provides: 100% description coverage across all 51 files; regenerated docs
- Size: Small (4 files + go generate)

## Implementation Content
Add `MarkdownDescription` to the remaining data source files (bucket, datacenter, ip, storage_snapshot_data_resource) and run `go generate ./...` to regenerate documentation from the completed schemas. This is the final Phase 4 task.

## Target Files
- [ ] `internal/services/bucket/bucket_resource.go` (descriptions + validators)
- [ ] `internal/services/bucket/bucket_datasource.go` (descriptions only)
- [ ] `internal/services/datacenter/datacenter_datasource.go` (descriptions only)
- [ ] `internal/services/ip/ip_datasource.go` (descriptions only)
- [ ] `internal/services/storage/storage_snapshot_data_resource.go` (descriptions only -- confirmed data source despite filename)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each target file to catalog all schema attributes
- [ ] Confirm `storage_snapshot_data_resource.go` is a data source (not a resource)
- [ ] Verify all 51 files across tasks 04-11 have been covered

### 2. Implementation
- [ ] Add `MarkdownDescription` to every attribute in `bucket_resource.go` (descriptions + validators for Required strings)
- [ ] Add `MarkdownDescription` to every attribute in `bucket_datasource.go` (descriptions only)
- [ ] Add `MarkdownDescription` to every attribute in `datacenter_datasource.go` (descriptions only)
- [ ] Add `MarkdownDescription` to every attribute in `ip_datasource.go` (descriptions only)
- [ ] Add `MarkdownDescription` to every attribute in `storage_snapshot_data_resource.go` (descriptions only)
- [ ] Run `go generate ./...` to regenerate docs

### 3. Verify Skill Fidelity
- [ ] Spot-check 5 resource files and 5 data source files for description completeness across all Phase 4 tasks
- [ ] Verify `docs/` directory reflects all attribute descriptions

## Completion Criteria
- [ ] All 5 files in this task have `MarkdownDescription` on every attribute
- [ ] `go build -v .` succeeds
- [ ] `go generate ./...` succeeds and `docs/` directory is updated
- [ ] 100% of the 51 resource/data source files now have descriptions
- [ ] Verification level: L3 (build success + generate success)

## Verification Steps
```bash
# 1. Verify remaining files have descriptions
grep -c "MarkdownDescription" internal/services/bucket/bucket_resource.go
grep -c "MarkdownDescription" internal/services/bucket/bucket_datasource.go
grep -c "MarkdownDescription" internal/services/datacenter/datacenter_datasource.go
grep -c "MarkdownDescription" internal/services/ip/ip_datasource.go
grep -c "MarkdownDescription" internal/services/storage/storage_snapshot_data_resource.go

# 2. Full coverage check: find files WITHOUT MarkdownDescription
grep -rL "MarkdownDescription" internal/services/*/*.go | grep -v _test | grep -v _api
# Expected: zero results (all files have descriptions)

# 3. Build and generate
go build -v .
go generate ./...

# 4. Verify docs directory updated
ls -la docs/resources/ docs/data-sources/
```

## Notes
- Impact scope: Schema metadata + documentation generation
- Constraints: `storage_snapshot_data_resource.go` is a data source despite the misleading filename
- AC Coverage: AC-4.1 (complete), AC-4.2, AC-4.5
- This task completes Phase 4
