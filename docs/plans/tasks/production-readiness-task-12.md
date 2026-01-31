# Task: Phase 5a - Add CheckDestroy to All 6 Existing Test Files

Metadata:
- Dependencies: Phase 4 completion (all description/validator tasks done)
- Provides: CheckDestroy functions in all 6 existing acceptance test files
- Size: Medium (6 files)

## Implementation Content
Add `CheckDestroy` functions to all 6 existing acceptance test files. Each CheckDestroy function constructs a `govpsie.Client` directly from the `VPSIE_ACCESS_TOKEN` environment variable (the `acctest` package does not expose a provider instance) and calls the appropriate SDK List method to verify the resource no longer exists.

## Target Files
- [ ] `internal/services/storage/storage_resource_test.go` (add testAccCheckStorageResourceDestroy)
- [ ] `internal/services/sshkey/sshkey_resource_test.go` (add testAccCheckSshkeyResourceDestroy)
- [ ] `internal/services/script/script_resource_test.go` (add testAccCheckScriptResourceDestroy)
- [ ] `internal/services/project/project_resource_test.go` (add testAccCheckProjectResourceDestroy)
- [ ] `internal/services/domain/domain_resource_test.go` (add testAccCheckDomainResourceDestroy)
- [ ] `internal/services/accesstoken/accesstoken_resource_test.go` (add testAccCheckAccessTokenResourceDestroy)

## Implementation Steps
### 1. Confirm Skill Constraints
- [ ] Read each existing test file to understand the current TestCase structure
- [ ] Verify which SDK List/Get methods are available for each service
- [ ] Confirm the `govpsie.Client` construction pattern from env var

### 2. Implementation
For each test file:
- [ ] Create a `testAccCheck<Service>ResourceDestroy` function that:
  1. Constructs a `govpsie.Client` from `VPSIE_ACCESS_TOKEN` env var
  2. Iterates over `s.RootModule().Resources` filtering by resource type
  3. Calls the appropriate SDK List/Get method to check if resource still exists
  4. Returns error if resource still exists
- [ ] Add `CheckDestroy: testAccCheck<Service>ResourceDestroy` to each `resource.TestCase`

### 3. Verify Skill Fidelity
- [ ] All 6 test files have CheckDestroy functions
- [ ] Each CheckDestroy constructs its own client (does not rely on provider instance)

## Completion Criteria
- [ ] All 6 test files have `CheckDestroy` functions
- [ ] Each TestCase references its CheckDestroy function
- [ ] `go build -v .` succeeds
- [ ] `go test ./internal/services/storage/ -count=0` (compile check, no TF_ACC)
- [ ] Verification level: L3 (build success)

## Verification Steps
```bash
# 1. Verify CheckDestroy present in all 6 files
grep -rn "CheckDestroy" internal/services/storage/storage_resource_test.go
grep -rn "CheckDestroy" internal/services/sshkey/sshkey_resource_test.go
grep -rn "CheckDestroy" internal/services/script/script_resource_test.go
grep -rn "CheckDestroy" internal/services/project/project_resource_test.go
grep -rn "CheckDestroy" internal/services/domain/domain_resource_test.go
grep -rn "CheckDestroy" internal/services/accesstoken/accesstoken_resource_test.go
# Expected: Each file returns at least one match

# 2. Build verification
go build -v .
```

## Notes
- Impact scope: Test files only; no production code changes
- Constraints: Must construct `govpsie.Client` directly from env var (acctest does not expose provider)
- AC Coverage: AC-5.2
- The CheckDestroy pattern:
  ```go
  func testAccCheck<Service>ResourceDestroy(s *terraform.State) error {
      token := os.Getenv("VPSIE_ACCESS_TOKEN")
      client, _ := govpsie.NewClient(token)
      for _, rs := range s.RootModule().Resources {
          if rs.Type != "vpsie_<resource>" { continue }
          // Call SDK to check resource still exists
          // Return error if found
      }
      return nil
  }
  ```
