# Task: Phase 4c - Descriptions and Validators for kubernetes, loadbalancer, firewall Services

Metadata:
- Dependencies: task-04 (validator dependency must be available)
- Provides: Complete MarkdownDescription + validators for kubernetes, loadbalancer, firewall services
- Size: Medium (9 files)

## Implementation Content
Add `MarkdownDescription` to every attribute and `Validators` to Required string attributes in the kubernetes, loadbalancer, and firewall service files (both resources and data sources). All files within each service are grouped together.

## Target Files
- [ ] `internal/services/kubernetes/kubernetes_resource.go` (descriptions + validators)
- [ ] `internal/services/kubernetes/kubernetes_group_resource.go` (descriptions + validators)
- [ ] `internal/services/kubernetes/kubernetes_data_source.go` (descriptions only)
- [ ] `internal/services/kubernetes/kubernetes_group_data_source.go` (descriptions only)
- [ ] `internal/services/loadbalancer/loadbalancer_resource.go` (descriptions + validators)
- [ ] `internal/services/loadbalancer/loadbalancer_data_source.go` (descriptions only)
- [ ] `internal/services/firewall/firewall_resource.go` (descriptions + validators)
- [ ] `internal/services/firewall/firewall_attachment_resource.go` (descriptions + validators)
- [ ] `internal/services/firewall/firewall_data_source.go` (descriptions only)

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
- [ ] Add `MarkdownDescription` to every attribute (descriptions only, no validators for Computed)

### 3. Verify Skill Fidelity
- [ ] Every attribute in all 9 files has a non-empty MarkdownDescription
- [ ] All Required string attributes have LengthAtLeast(1)

## Completion Criteria
- [ ] All 9 files have `MarkdownDescription` on every attribute
- [ ] All Required string attributes have validators
- [ ] `go build -v .` succeeds
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Spot-check descriptions
grep -c "MarkdownDescription" internal/services/kubernetes/kubernetes_resource.go
grep -c "MarkdownDescription" internal/services/loadbalancer/loadbalancer_resource.go
grep -c "MarkdownDescription" internal/services/firewall/firewall_resource.go
grep -c "MarkdownDescription" internal/services/firewall/firewall_data_source.go

# 2. Spot-check validators
grep -c "LengthAtLeast" internal/services/kubernetes/kubernetes_resource.go
grep -c "Validators:" internal/services/loadbalancer/loadbalancer_resource.go

# 3. Build
go build -v .
```

## Notes
- Impact scope: Schema metadata only; no CRUD logic changes
- Constraints: Use only permissive validators
- AC Coverage: AC-4.1 (partial), AC-4.3 (partial)
