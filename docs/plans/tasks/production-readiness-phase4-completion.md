# Phase 4 Completion: Schema Validators and Descriptions

## Phase Summary
Phase 4 adds `MarkdownDescription` to every attribute in all 51 resource/data source files and adds validators to Required and meaningful Optional attributes.

## Task Completion Checklist
- [ ] Task 04: Validator dependency + provider descriptions
- [ ] Task 05: server, sshkey, storage services (8 files)
- [ ] Task 06: kubernetes, loadbalancer, firewall services (9 files)
- [ ] Task 07: domain services (4 files)
- [ ] Task 08: backup, snapshot services (8 files)
- [ ] Task 09: vpc, gateway, fip services (7 files)
- [ ] Task 10: project, script, monitoring, accesstoken, image services (10 files)
- [ ] Task 11: bucket, datacenter, ip + storage_snapshot_data_resource + doc regen (5 files)

## E2E Verification Procedures (from Design Doc)

### 1. Verify validator dependency
```bash
grep "terraform-plugin-framework-validators" /Users/zozo/projects/terraform-provider-vpsie/go.mod
```
Expected: Version entry present

### 2. Full description coverage check
```bash
# Find files WITHOUT MarkdownDescription (should be zero excluding test and api files)
grep -rL "MarkdownDescription" internal/services/*/*.go | grep -v _test | grep -v _api
```
Expected: Zero results

### 3. Spot-check descriptions (sample 3 files)
```bash
grep -c "MarkdownDescription" internal/services/storage/storage_resource.go
grep -c "MarkdownDescription" internal/services/sshkey/sshkey_resource.go
grep -c "MarkdownDescription" internal/services/server/server_data_source.go
```
Expected: Each file returns a count matching its number of attributes

### 4. Spot-check validators
```bash
grep -c "Validators:" internal/services/storage/storage_resource.go
grep -c "LengthAtLeast" internal/services/sshkey/sshkey_resource.go
```
Expected: Non-zero counts for Required attributes

### 5. Verify OneOf validators
```bash
grep "OneOf" internal/services/storage/storage_resource.go   # disk_format
grep "OneOf" internal/services/domain/dns_record_resource.go  # record type
grep "OneOf" internal/services/fip/fip_resource.go             # ip_type
grep "OneOf" internal/services/monitoring/monitoring_rule_resource.go  # condition
```
Expected: Each returns a match

### 6. Build and generate
```bash
go build -v .
go generate ./...
```
Expected: Both succeed

## Phase Completion Criteria
- [ ] `terraform-plugin-framework-validators` in `go.mod`
- [ ] All 27 resource files have `MarkdownDescription` on every attribute
- [ ] All 24 data source files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have `stringvalidator.LengthAtLeast(1)`
- [ ] Known enum fields have `stringvalidator.OneOf(...)` validators
- [ ] Positive-only numeric fields have `int64validator.AtLeast(1)` or `AtLeast(0)`
- [ ] `go build -v .` succeeds
- [ ] `go generate ./...` succeeds and `docs/` is updated
- [ ] Provider-level schema has `MarkdownDescription`

## AC Coverage
- AC-4.1: Every attribute in all 27 resources and 24 data sources has non-empty MarkdownDescription
- AC-4.2: terraform providers schema -json reports descriptions
- AC-4.3: Empty string for Required string attribute rejected at validate
- AC-4.4: Invalid value for validated Optional attribute rejected at validate
- AC-4.5: go generate updates docs with attribute descriptions
