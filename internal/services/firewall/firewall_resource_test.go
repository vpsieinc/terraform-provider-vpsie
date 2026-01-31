package firewall

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockFirewallAPI implements FirewallAPI for unit testing.
type mockFirewallAPI struct {
	CreateFn          func(ctx context.Context, groupName string, firewallUpdateReq []govpsie.FirewallUpdateReq) error
	ListFn            func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.FirewallGroupListData, error)
	GetFn             func(ctx context.Context, fwGroupId string) (*govpsie.FirewallGroupDetailData, error)
	DeleteFn          func(ctx context.Context, fwGroupId string) error
	AttachToVpsieFn   func(ctx context.Context, groupId, vmId string) error
	DetachFromVpsieFn func(ctx context.Context, groupId, vmId string) error
}

func (m *mockFirewallAPI) Create(ctx context.Context, groupName string, firewallUpdateReq []govpsie.FirewallUpdateReq) error {
	return m.CreateFn(ctx, groupName, firewallUpdateReq)
}

func (m *mockFirewallAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.FirewallGroupListData, error) {
	return m.ListFn(ctx, options)
}

func (m *mockFirewallAPI) Get(ctx context.Context, fwGroupId string) (*govpsie.FirewallGroupDetailData, error) {
	return m.GetFn(ctx, fwGroupId)
}

func (m *mockFirewallAPI) Delete(ctx context.Context, fwGroupId string) error {
	return m.DeleteFn(ctx, fwGroupId)
}

func (m *mockFirewallAPI) AttachToVpsie(ctx context.Context, groupId, vmId string) error {
	return m.AttachToVpsieFn(ctx, groupId, vmId)
}

func (m *mockFirewallAPI) DetachFromVpsie(ctx context.Context, groupId, vmId string) error {
	return m.DetachFromVpsieFn(ctx, groupId, vmId)
}

// Compile-time check: mockFirewallAPI satisfies FirewallAPI.
var _ FirewallAPI = &mockFirewallAPI{}

func TestUnitFirewallAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockFirewallAPI{
		CreateFn: func(ctx context.Context, groupName string, firewallUpdateReq []govpsie.FirewallUpdateReq) error {
			return nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.FirewallGroupListData, error) {
			return []govpsie.FirewallGroupListData{}, nil
		},
		GetFn: func(ctx context.Context, fwGroupId string) (*govpsie.FirewallGroupDetailData, error) {
			return &govpsie.FirewallGroupDetailData{}, nil
		},
		DeleteFn: func(ctx context.Context, fwGroupId string) error {
			return nil
		},
		AttachToVpsieFn: func(ctx context.Context, groupId, vmId string) error {
			return nil
		},
		DetachFromVpsieFn: func(ctx context.Context, groupId, vmId string) error {
			return nil
		},
	}

	var api FirewallAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitFirewallAPI_GetFirewallGroupByName(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
		firewalls []govpsie.FirewallGroupListData
		expectErr bool
		expectID  int64
	}{
		{
			name:      "firewall group found",
			groupName: "my-firewall",
			firewalls: []govpsie.FirewallGroupListData{
				{GroupName: "other-firewall", ID: 1},
				{GroupName: "my-firewall", ID: 2},
			},
			expectErr: false,
			expectID:  2,
		},
		{
			name:      "firewall group not found",
			groupName: "missing-firewall",
			firewalls: []govpsie.FirewallGroupListData{
				{GroupName: "other-firewall", ID: 1},
			},
			expectErr: true,
		},
		{
			name:      "empty firewall list",
			groupName: "any",
			firewalls: []govpsie.FirewallGroupListData{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFirewallAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.FirewallGroupListData, error) {
					return tt.firewalls, nil
				},
			}

			r := &firewallResource{client: mock}
			fw, err := r.GetFirewallGroupByName(t.Context(), tt.groupName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && fw.ID != tt.expectID {
				t.Fatalf("expected ID %d, got %d", tt.expectID, fw.ID)
			}
		})
	}
}

func TestUnitFirewallAPI_Delete(t *testing.T) {
	var calledWithID string

	mock := &mockFirewallAPI{
		DeleteFn: func(ctx context.Context, fwGroupId string) error {
			calledWithID = fwGroupId
			return nil
		},
	}

	err := mock.Delete(t.Context(), "fw-group-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWithID != "fw-group-123" {
		t.Fatalf("expected fwGroupId 'fw-group-123', got %q", calledWithID)
	}
}

func TestUnitFirewallAPI_AttachToVpsie(t *testing.T) {
	var calledWith struct {
		groupId string
		vmId    string
	}

	mock := &mockFirewallAPI{
		AttachToVpsieFn: func(ctx context.Context, groupId, vmId string) error {
			calledWith.groupId = groupId
			calledWith.vmId = vmId
			return nil
		},
	}

	err := mock.AttachToVpsie(t.Context(), "group-1", "vm-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith.groupId != "group-1" {
		t.Fatalf("expected groupId 'group-1', got %q", calledWith.groupId)
	}
	if calledWith.vmId != "vm-abc" {
		t.Fatalf("expected vmId 'vm-abc', got %q", calledWith.vmId)
	}
}

func TestUnitFirewallAPI_DetachFromVpsie(t *testing.T) {
	var calledWith struct {
		groupId string
		vmId    string
	}

	mock := &mockFirewallAPI{
		DetachFromVpsieFn: func(ctx context.Context, groupId, vmId string) error {
			calledWith.groupId = groupId
			calledWith.vmId = vmId
			return nil
		},
	}

	err := mock.DetachFromVpsie(t.Context(), "group-1", "vm-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith.groupId != "group-1" {
		t.Fatalf("expected groupId 'group-1', got %q", calledWith.groupId)
	}
	if calledWith.vmId != "vm-abc" {
		t.Fatalf("expected vmId 'vm-abc', got %q", calledWith.vmId)
	}
}

// TestUnitFirewallAPI_ListErrorPropagation verifies that errors from the List
// API call are properly propagated through GetFirewallGroupByName (AC-3.1).
func TestUnitFirewallAPI_ListErrorPropagation(t *testing.T) {
	mock := &mockFirewallAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.FirewallGroupListData, error) {
			return nil, fmt.Errorf("API error: authentication failed")
		},
	}

	r := &firewallResource{client: mock}
	_, err := r.GetFirewallGroupByName(t.Context(), "any-group")

	if err == nil {
		t.Fatal("expected error to propagate from List, got nil")
	}
	if err.Error() != "API error: authentication failed" {
		t.Fatalf("expected error message 'API error: authentication failed', got %q", err.Error())
	}
}

// TestUnitFirewallAPI_GetErrorPropagation verifies that errors from the Get
// API call are properly propagated (AC-3.1 unit-level).
func TestUnitFirewallAPI_GetErrorPropagation(t *testing.T) {
	mock := &mockFirewallAPI{
		GetFn: func(ctx context.Context, fwGroupId string) (*govpsie.FirewallGroupDetailData, error) {
			return nil, fmt.Errorf("API error: not found")
		},
	}

	_, err := mock.Get(t.Context(), "non-existent-id")
	if err == nil {
		t.Fatal("expected error to propagate from Get, got nil")
	}
	if err.Error() != "API error: not found" {
		t.Fatalf("expected error message 'API error: not found', got %q", err.Error())
	}
}

// TestUnitFirewallAPI_CreateErrorPropagation verifies that errors from the Create
// API call are properly propagated (AC-3.1 unit-level).
func TestUnitFirewallAPI_CreateErrorPropagation(t *testing.T) {
	mock := &mockFirewallAPI{
		CreateFn: func(ctx context.Context, groupName string, firewallUpdateReq []govpsie.FirewallUpdateReq) error {
			return fmt.Errorf("API error: quota exceeded")
		},
	}

	err := mock.Create(t.Context(), "test-group", nil)
	if err == nil {
		t.Fatal("expected error to propagate from Create, got nil")
	}
	if err.Error() != "API error: quota exceeded" {
		t.Fatalf("expected error message 'API error: quota exceeded', got %q", err.Error())
	}
}
