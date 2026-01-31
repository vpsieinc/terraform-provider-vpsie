# Overall Design Document: Production Readiness Hardening

Generation Date: 2026-01-31
Target Plan Document: production-readiness-workplan.md

## Project Overview

### Purpose and Goals
Harden the VPSie Terraform provider to meet HashiCorp production quality standards by addressing six categories of gaps: credential exposure, local build dependencies, silenced errors, missing schema documentation and validation, insufficient test coverage, and incomplete import support. All changes must be backward-compatible and executed in strict priority order.

### Background and Context
The provider (27 resources, 24 data sources, 20 service modules) has been audited against HashiCorp production quality standards. The audit found: 5 sensitive fields exposed in plaintext, a local filesystem `replace` directive in `go.mod`, 12 silently discarded errors in firewall code, zero attribute descriptions or validators, test coverage in only 6 of 20 services, and 7 resources missing `ImportState`.

## Task Division Design

### Division Policy
Tasks are divided using a **horizontal slice** approach following the strict phase ordering mandated by the work plan. Within Phase 4 (51 files) and Phase 5 (~40 files), tasks are further split by service groupings to keep each task within the 1-5 file medium size limit. Phases 1, 2, and 3 are each small enough for a single task.

- Verification levels follow the implementation-approach skill:
  - Phases 1-3: L3 (build success) + L2 (test compilation)
  - Phase 4: L3 (build + generate)
  - Phase 5: L2 (all unit tests pass)
  - Phase 6: L2 (all tests pass) + L1 (full lifecycle works)

### Inter-task Relationship Map
```
Task 01: Phase 1 - Sensitive field markings (5 files)
  |
Task 02: Phase 2 - Remove SDK replace directive (1 file)
  |
Task 03: Phase 3 - Fix firewall error suppression (2 files)
  |
Task 04: Phase 4a - Validator dependency + provider descriptions (2 files)
  |
Task 05: Phase 4b - Descriptions/validators: server, sshkey, storage (8 files split into batch)
  |
Task 06: Phase 4c - Descriptions/validators: kubernetes, loadbalancer, firewall (8 files)
  |
Task 07: Phase 4d - Descriptions/validators: domain, dns_record, reverse_dns (5 files)
  |
Task 08: Phase 4e - Descriptions/validators: backup, snapshot (8 files)
  |
Task 09: Phase 4f - Descriptions/validators: vpc, gateway, fip (7 files)
  |
Task 10: Phase 4g - Descriptions/validators: project, script, monitoring, accesstoken, image (10 files)
  |
Task 11: Phase 4h - Descriptions/validators: bucket, datacenter, ip + doc regen (4 files + generate)
  |
Phase 4 Completion: Verify all 51 files have descriptions
  |
Task 12: Phase 5a - CheckDestroy for 6 existing test files (6 files)
  |
Task 13: Phase 5b - Interface wrappers + unit tests: storage, server, sshkey (9 files)
  |
Task 14: Phase 5c - Interface wrappers + unit tests: kubernetes, loadbalancer, firewall (9 files)
  |
Task 15: Phase 5d - Interface wrappers + unit tests: domain, backup, snapshot (9 files)
  |
Task 16: Phase 5e - Interface wrappers + unit tests: vpc, gateway, fip, monitoring (12 files)
  |
Task 17: Phase 5f - Interface wrappers + unit tests: project, script, accesstoken, bucket, datacenter, ip, image (14 files)
  |
Phase 5 Completion: Verify all 20 services have tests
  |
Task 18: Phase 6 - ImportState implementations (7 files)
  |
Phase 6 Completion: Final quality gate
```

### Interface Change Impact Analysis
| Existing Interface | New Interface | Conversion Required | Corresponding Task |
|-------------------|---------------|-------------------|-------------------|
| `*govpsie.Client` in resource structs | Service-specific interface | Yes (Phase 5) | Tasks 13-17 |
| Schema attributes (no descriptions) | Schema attributes + MarkdownDescription | No (additive) | Tasks 05-11 |
| Schema attributes (no validators) | Schema attributes + Validators | No (additive) | Tasks 05-11 |
| Resources without ImportState | Resources with ImportState | No (additive) | Task 18 |

### Common Processing Points
- MarkdownDescription pattern: applied identically across all 51 files (Phase 4)
- Validator pattern: `stringvalidator.LengthAtLeast(1)` for all Required strings (Phase 4)
- Interface wrapper pattern: narrow interface per service wrapping `client.<SubService>` (Phase 5)
- Unit test pattern: table-driven tests with mock interface implementation (Phase 5)
- CheckDestroy pattern: construct `govpsie.Client` from env var (Phase 5)
- ImportState passthrough pattern: `resource.ImportStatePassthroughID` (Phase 6)
- Composite ID pattern: split on `/`, validate part count, set attributes (Phase 6)

## Implementation Considerations

### Principles to Maintain Throughout
1. All changes must be backward-compatible (no state migrations)
2. Validators must be permissive (`LengthAtLeast(1)`) not restrictive
3. Interfaces must be narrow (only methods used by that service)
4. No modifications to the `govpsie` SDK itself
5. Strict phase ordering: each phase completes before next begins

### Risks and Countermeasures
- Risk: Published SDK differs from local copy
  Countermeasure: Full build verification after Phase 2; re-add replace temporarily if needed
- Risk: Validators reject previously valid configurations
  Countermeasure: Use only permissive validators; test with known valid configs
- Risk: Phase 4 (51 files) takes longer than estimated
  Countermeasure: Systematic approach with common patterns first; service-specific second

### Impact Scope Management
- Allowed change scope: `internal/provider/provider.go`, all files under `internal/services/`, `go.mod`, `go.sum`, `tests/`, `docs/`
- No-change areas: `govpsie` SDK, `.github/`, CI/CD configuration, `main.go`, `internal/acctest/`
