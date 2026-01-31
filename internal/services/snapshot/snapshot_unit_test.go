package snapshot

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockSnapshotAPI implements SnapshotAPI for unit testing.
type mockSnapshotAPI struct {
	CreateFn               func(ctx context.Context, name, vmIdentifier, note string) error
	ListFn                 func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Snapshot, error)
	GetFn                  func(ctx context.Context, buckupIdentifier string) (*govpsie.Snapshot, error)
	UpdateFn               func(ctx context.Context, snapshotIdentifier, newNote string) error
	DeleteFn               func(ctx context.Context, snapshotIdentifier, reason, note string) error
	CreateSnapShotPolicyFn func(ctx context.Context, createReq *govpsie.CreateSnapShotPolicyReq) error
	GetSnapShotPolicyFn    func(ctx context.Context, identifier string) (*govpsie.SnapShotPolicy, error)
	ListSnapShotPoliciesFn func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.SnapShotPolicyListDetail, error)
	AttachSnapShotPolicyFn func(ctx context.Context, policyId string, vms []string) error
	DetachSnapShotPolicyFn func(ctx context.Context, policyId string, vms []string) error
	DeleteSnapShotPolicyFn func(ctx context.Context, policyId, identifier string) error
}

func (m *mockSnapshotAPI) Create(ctx context.Context, name, vmIdentifier, note string) error {
	return m.CreateFn(ctx, name, vmIdentifier, note)
}

func (m *mockSnapshotAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Snapshot, error) {
	return m.ListFn(ctx, options)
}

func (m *mockSnapshotAPI) Get(ctx context.Context, buckupIdentifier string) (*govpsie.Snapshot, error) {
	return m.GetFn(ctx, buckupIdentifier)
}

func (m *mockSnapshotAPI) Update(ctx context.Context, snapshotIdentifier, newNote string) error {
	return m.UpdateFn(ctx, snapshotIdentifier, newNote)
}

func (m *mockSnapshotAPI) Delete(ctx context.Context, snapshotIdentifier, reason, note string) error {
	return m.DeleteFn(ctx, snapshotIdentifier, reason, note)
}

func (m *mockSnapshotAPI) CreateSnapShotPolicy(ctx context.Context, createReq *govpsie.CreateSnapShotPolicyReq) error {
	return m.CreateSnapShotPolicyFn(ctx, createReq)
}

func (m *mockSnapshotAPI) GetSnapShotPolicy(ctx context.Context, identifier string) (*govpsie.SnapShotPolicy, error) {
	return m.GetSnapShotPolicyFn(ctx, identifier)
}

func (m *mockSnapshotAPI) ListSnapShotPolicies(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.SnapShotPolicyListDetail, error) {
	return m.ListSnapShotPoliciesFn(ctx, options)
}

func (m *mockSnapshotAPI) AttachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error {
	return m.AttachSnapShotPolicyFn(ctx, policyId, vms)
}

func (m *mockSnapshotAPI) DetachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error {
	return m.DetachSnapShotPolicyFn(ctx, policyId, vms)
}

func (m *mockSnapshotAPI) DeleteSnapShotPolicy(ctx context.Context, policyId, identifier string) error {
	return m.DeleteSnapShotPolicyFn(ctx, policyId, identifier)
}

// Compile-time check: mockSnapshotAPI satisfies SnapshotAPI.
var _ SnapshotAPI = &mockSnapshotAPI{}

func TestUnitSnapshotAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockSnapshotAPI{
		CreateFn:               func(ctx context.Context, name, vmIdentifier, note string) error { return nil },
		ListFn:                 func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Snapshot, error) { return nil, nil },
		GetFn:                  func(ctx context.Context, buckupIdentifier string) (*govpsie.Snapshot, error) { return nil, nil },
		UpdateFn:               func(ctx context.Context, snapshotIdentifier, newNote string) error { return nil },
		DeleteFn:               func(ctx context.Context, snapshotIdentifier, reason, note string) error { return nil },
		CreateSnapShotPolicyFn: func(ctx context.Context, createReq *govpsie.CreateSnapShotPolicyReq) error { return nil },
		GetSnapShotPolicyFn:    func(ctx context.Context, identifier string) (*govpsie.SnapShotPolicy, error) { return nil, nil },
		ListSnapShotPoliciesFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.SnapShotPolicyListDetail, error) {
			return nil, nil
		},
		AttachSnapShotPolicyFn: func(ctx context.Context, policyId string, vms []string) error { return nil },
		DetachSnapShotPolicyFn: func(ctx context.Context, policyId string, vms []string) error { return nil },
		DeleteSnapShotPolicyFn: func(ctx context.Context, policyId, identifier string) error { return nil },
	}

	var api SnapshotAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitSnapshotAPI_GetSnapshotByName(t *testing.T) {
	tests := []struct {
		name         string
		snapshotName string
		snapshots    []govpsie.Snapshot
		expectFound  bool
		expectErr    bool
	}{
		{
			name:         "snapshot found by name",
			snapshotName: "test-snap",
			snapshots: []govpsie.Snapshot{
				{Name: "other-snap", Identifier: "id-1"},
				{Name: "test-snap", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:         "snapshot not found",
			snapshotName: "missing-snap",
			snapshots: []govpsie.Snapshot{
				{Name: "other-snap", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:         "empty snapshot list",
			snapshotName: "test-snap",
			snapshots:    []govpsie.Snapshot{},
			expectFound:  false,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSnapshotAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Snapshot, error) {
					return tt.snapshots, nil
				},
			}

			r := &serverSnapshotResource{client: mock}
			snap, err := r.GetSnapshotByName(t.Context(), tt.snapshotName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && snap == nil {
				t.Fatal("expected snapshot to be non-nil when found")
			}
			if tt.expectFound && snap.Name != tt.snapshotName {
				t.Fatalf("expected name %q, got %q", tt.snapshotName, snap.Name)
			}
		})
	}
}

func TestUnitSnapshotAPI_GetPolicyByName(t *testing.T) {
	tests := []struct {
		name        string
		policyName  string
		policies    []govpsie.SnapShotPolicyListDetail
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "policy found by name",
			policyName: "daily-snap-policy",
			policies: []govpsie.SnapShotPolicyListDetail{
				{Name: "weekly-snap-policy", Identifier: "id-1"},
				{Name: "daily-snap-policy", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "policy not found",
			policyName: "missing-policy",
			policies: []govpsie.SnapShotPolicyListDetail{
				{Name: "weekly-snap-policy", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty policy list",
			policyName:  "daily-snap-policy",
			policies:    []govpsie.SnapShotPolicyListDetail{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSnapshotAPI{
				ListSnapShotPoliciesFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.SnapShotPolicyListDetail, error) {
					return tt.policies, nil
				},
			}

			r := &snapshotPolicyResource{client: mock}
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

func TestUnitSnapshotAPI_ListError(t *testing.T) {
	mock := &mockSnapshotAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Snapshot, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &serverSnapshotResource{client: mock}
	_, err := r.GetSnapshotByName(t.Context(), "test")
	if err == nil {
		t.Fatal("expected error from List failure, got nil")
	}
}

func TestUnitSnapshotAPI_ListSnapShotPoliciesError(t *testing.T) {
	mock := &mockSnapshotAPI{
		ListSnapShotPoliciesFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.SnapShotPolicyListDetail, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &snapshotPolicyResource{client: mock}
	_, err := r.GetPolicyByName(t.Context(), "test")
	if err == nil {
		t.Fatal("expected error from ListSnapShotPolicies failure, got nil")
	}
}
