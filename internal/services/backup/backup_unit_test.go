package backup

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockBackupAPI implements BackupAPI for unit testing.
type mockBackupAPI struct {
	CreateBackupsFn      func(ctx context.Context, vmIdentifier, name, notes string) error
	ListFn               func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Backup, error)
	GetFn                func(ctx context.Context, identifier string) (*govpsie.Backup, error)
	RenameFn             func(ctx context.Context, backupIdentifier, newName string) error
	DeleteBackupFn       func(ctx context.Context, backupIdentifier, deleteReason, deleteNote string) error
	CreateBackupPolicyFn func(ctx context.Context, createReq *govpsie.CreateBackupPolicyReq) error
	GetBackupPolicyFn    func(ctx context.Context, identifier string) (*govpsie.BackupPolicy, error)
	ListBackupPoliciesFn func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.BackupPolicyListDetail, error)
	AttachBackupPolicyFn func(ctx context.Context, policyId string, vms []string) error
	DetachBackupPolicyFn func(ctx context.Context, policyId string, vms []string) error
	DeleteBackupPolicyFn func(ctx context.Context, policyId, identifier string) error
}

func (m *mockBackupAPI) CreateBackups(ctx context.Context, vmIdentifier, name, notes string) error {
	return m.CreateBackupsFn(ctx, vmIdentifier, name, notes)
}

func (m *mockBackupAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Backup, error) {
	return m.ListFn(ctx, options)
}

func (m *mockBackupAPI) Get(ctx context.Context, identifier string) (*govpsie.Backup, error) {
	return m.GetFn(ctx, identifier)
}

func (m *mockBackupAPI) Rename(ctx context.Context, backupIdentifier, newName string) error {
	return m.RenameFn(ctx, backupIdentifier, newName)
}

func (m *mockBackupAPI) DeleteBackup(ctx context.Context, backupIdentifier, deleteReason, deleteNote string) error {
	return m.DeleteBackupFn(ctx, backupIdentifier, deleteReason, deleteNote)
}

func (m *mockBackupAPI) CreateBackupPolicy(ctx context.Context, createReq *govpsie.CreateBackupPolicyReq) error {
	return m.CreateBackupPolicyFn(ctx, createReq)
}

func (m *mockBackupAPI) GetBackupPolicy(ctx context.Context, identifier string) (*govpsie.BackupPolicy, error) {
	return m.GetBackupPolicyFn(ctx, identifier)
}

func (m *mockBackupAPI) ListBackupPolicies(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.BackupPolicyListDetail, error) {
	return m.ListBackupPoliciesFn(ctx, options)
}

func (m *mockBackupAPI) AttachBackupPolicy(ctx context.Context, policyId string, vms []string) error {
	return m.AttachBackupPolicyFn(ctx, policyId, vms)
}

func (m *mockBackupAPI) DetachBackupPolicy(ctx context.Context, policyId string, vms []string) error {
	return m.DetachBackupPolicyFn(ctx, policyId, vms)
}

func (m *mockBackupAPI) DeleteBackupPolicy(ctx context.Context, policyId, identifier string) error {
	return m.DeleteBackupPolicyFn(ctx, policyId, identifier)
}

// Compile-time check: mockBackupAPI satisfies BackupAPI.
var _ BackupAPI = &mockBackupAPI{}

func TestUnitBackupAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockBackupAPI{
		CreateBackupsFn:      func(ctx context.Context, vmIdentifier, name, notes string) error { return nil },
		ListFn:               func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Backup, error) { return nil, nil },
		GetFn:                func(ctx context.Context, identifier string) (*govpsie.Backup, error) { return nil, nil },
		RenameFn:             func(ctx context.Context, backupIdentifier, newName string) error { return nil },
		DeleteBackupFn:       func(ctx context.Context, backupIdentifier, deleteReason, deleteNote string) error { return nil },
		CreateBackupPolicyFn: func(ctx context.Context, createReq *govpsie.CreateBackupPolicyReq) error { return nil },
		GetBackupPolicyFn:    func(ctx context.Context, identifier string) (*govpsie.BackupPolicy, error) { return nil, nil },
		ListBackupPoliciesFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.BackupPolicyListDetail, error) {
			return nil, nil
		},
		AttachBackupPolicyFn: func(ctx context.Context, policyId string, vms []string) error { return nil },
		DetachBackupPolicyFn: func(ctx context.Context, policyId string, vms []string) error { return nil },
		DeleteBackupPolicyFn: func(ctx context.Context, policyId, identifier string) error { return nil },
	}

	var api BackupAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitBackupAPI_GetBackupByName(t *testing.T) {
	tests := []struct {
		name        string
		backupName  string
		backups     []govpsie.Backup
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "backup found by name",
			backupName: "test-backup",
			backups: []govpsie.Backup{
				{Name: "other-backup", Identifier: "id-1"},
				{Name: "test-backup", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "backup not found",
			backupName: "missing-backup",
			backups: []govpsie.Backup{
				{Name: "other-backup", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty backup list",
			backupName:  "test-backup",
			backups:     []govpsie.Backup{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockBackupAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Backup, error) {
					return tt.backups, nil
				},
			}

			r := &backupResource{client: mock}
			backup, err := r.GetBackupByName(t.Context(), tt.backupName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && backup == nil {
				t.Fatal("expected backup to be non-nil when found")
			}
			if tt.expectFound && backup.Name != tt.backupName {
				t.Fatalf("expected name %q, got %q", tt.backupName, backup.Name)
			}
		})
	}
}

func TestUnitBackupAPI_GetPolicyByName(t *testing.T) {
	tests := []struct {
		name        string
		policyName  string
		policies    []govpsie.BackupPolicyListDetail
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "policy found by name",
			policyName: "daily-policy",
			policies: []govpsie.BackupPolicyListDetail{
				{Name: "weekly-policy", Identifier: "id-1"},
				{Name: "daily-policy", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "policy not found",
			policyName: "missing-policy",
			policies: []govpsie.BackupPolicyListDetail{
				{Name: "weekly-policy", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty policy list",
			policyName:  "daily-policy",
			policies:    []govpsie.BackupPolicyListDetail{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockBackupAPI{
				ListBackupPoliciesFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.BackupPolicyListDetail, error) {
					return tt.policies, nil
				},
			}

			r := &backupPolicyResource{client: mock}
			policy, err := r.GetPolicyByName(t.Context(), tt.policyName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && policy == nil {
				t.Fatal("expected policy to be non-nil when found")
			}
			if tt.expectFound && policy.Name != tt.policyName {
				t.Fatalf("expected name %q, got %q", tt.policyName, policy.Name)
			}
		})
	}
}

func TestUnitBackupAPI_ListError(t *testing.T) {
	mock := &mockBackupAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Backup, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &backupResource{client: mock}
	_, err := r.GetBackupByName(t.Context(), "test")
	if err == nil {
		t.Fatal("expected error from List failure, got nil")
	}
}

func TestUnitBackupAPI_ListBackupPoliciesError(t *testing.T) {
	mock := &mockBackupAPI{
		ListBackupPoliciesFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.BackupPolicyListDetail, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &backupPolicyResource{client: mock}
	_, err := r.GetPolicyByName(t.Context(), "test")
	if err == nil {
		t.Fatal("expected error from ListBackupPolicies failure, got nil")
	}
}
