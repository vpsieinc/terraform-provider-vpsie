package kubernetes

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockKubernetesAPI implements KubernetesAPI for unit testing.
type mockKubernetesAPI struct {
	CreateFn         func(ctx context.Context, createReq *govpsie.CreateK8sReq) error
	ListFn           func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error)
	GetFn            func(ctx context.Context, identifier string) (*govpsie.K8s, error)
	DeleteFn         func(ctx context.Context, identifier, reason, note string) error
	AddSlaveFn       func(ctx context.Context, identifier string) error
	RemoveSlaveFn    func(ctx context.Context, identifier string) error
	ListK8sGroupsFn  func(ctx context.Context, identifier string) ([]govpsie.K8sGroup, error)
	CreateK8sGroupFn func(ctx context.Context, createReq *govpsie.CreateK8sGroupReq) error
	DeleteK8sGroupFn func(ctx context.Context, groupId string, reason, note string) error
	AddNodeFn        func(ctx context.Context, identifier, nodeType string, groupId int) error
	RemoveNodeFn     func(ctx context.Context, identifier, nodeType string, groupId int) error
}

func (m *mockKubernetesAPI) Create(ctx context.Context, createReq *govpsie.CreateK8sReq) error {
	return m.CreateFn(ctx, createReq)
}

func (m *mockKubernetesAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error) {
	return m.ListFn(ctx, options)
}

func (m *mockKubernetesAPI) Get(ctx context.Context, identifier string) (*govpsie.K8s, error) {
	return m.GetFn(ctx, identifier)
}

func (m *mockKubernetesAPI) Delete(ctx context.Context, identifier, reason, note string) error {
	return m.DeleteFn(ctx, identifier, reason, note)
}

func (m *mockKubernetesAPI) AddSlave(ctx context.Context, identifier string) error {
	return m.AddSlaveFn(ctx, identifier)
}

func (m *mockKubernetesAPI) RemoveSlave(ctx context.Context, identifier string) error {
	return m.RemoveSlaveFn(ctx, identifier)
}

func (m *mockKubernetesAPI) ListK8sGroups(ctx context.Context, identifier string) ([]govpsie.K8sGroup, error) {
	return m.ListK8sGroupsFn(ctx, identifier)
}

func (m *mockKubernetesAPI) CreateK8sGroup(ctx context.Context, createReq *govpsie.CreateK8sGroupReq) error {
	return m.CreateK8sGroupFn(ctx, createReq)
}

func (m *mockKubernetesAPI) DeleteK8sGroup(ctx context.Context, groupId string, reason, note string) error {
	return m.DeleteK8sGroupFn(ctx, groupId, reason, note)
}

func (m *mockKubernetesAPI) AddNode(ctx context.Context, identifier, nodeType string, groupId int) error {
	return m.AddNodeFn(ctx, identifier, nodeType, groupId)
}

func (m *mockKubernetesAPI) RemoveNode(ctx context.Context, identifier, nodeType string, groupId int) error {
	return m.RemoveNodeFn(ctx, identifier, nodeType, groupId)
}

// Compile-time check: mockKubernetesAPI satisfies KubernetesAPI.
var _ KubernetesAPI = &mockKubernetesAPI{}

func TestUnitKubernetesAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockKubernetesAPI{
		CreateFn: func(ctx context.Context, createReq *govpsie.CreateK8sReq) error {
			return nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error) {
			return []govpsie.ListK8s{}, nil
		},
		GetFn: func(ctx context.Context, identifier string) (*govpsie.K8s, error) {
			return &govpsie.K8s{}, nil
		},
		DeleteFn: func(ctx context.Context, identifier, reason, note string) error {
			return nil
		},
		AddSlaveFn: func(ctx context.Context, identifier string) error {
			return nil
		},
		RemoveSlaveFn: func(ctx context.Context, identifier string) error {
			return nil
		},
		ListK8sGroupsFn: func(ctx context.Context, identifier string) ([]govpsie.K8sGroup, error) {
			return []govpsie.K8sGroup{}, nil
		},
		CreateK8sGroupFn: func(ctx context.Context, createReq *govpsie.CreateK8sGroupReq) error {
			return nil
		},
		DeleteK8sGroupFn: func(ctx context.Context, groupId string, reason, note string) error {
			return nil
		},
		AddNodeFn: func(ctx context.Context, identifier, nodeType string, groupId int) error {
			return nil
		},
		RemoveNodeFn: func(ctx context.Context, identifier, nodeType string, groupId int) error {
			return nil
		},
	}

	var api KubernetesAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitKubernetesAPI_CheckResourceStatus(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		clusters    []govpsie.ListK8s
		getResult   *govpsie.K8s
		expectFound bool
		expectErr   bool
	}{
		{
			name:        "cluster found by name",
			clusterName: "test-cluster",
			clusters: []govpsie.ListK8s{
				{ClusterName: "other-cluster", Identifier: "id-1"},
				{ClusterName: "test-cluster", Identifier: "id-2"},
			},
			getResult: &govpsie.K8s{
				ClusterName: "test-cluster",
				Identifier:  "id-2",
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:        "cluster not found",
			clusterName: "missing-cluster",
			clusters: []govpsie.ListK8s{
				{ClusterName: "other-cluster", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   false,
		},
		{
			name:        "empty cluster list",
			clusterName: "test-cluster",
			clusters:    []govpsie.ListK8s{},
			expectFound: false,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKubernetesAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error) {
					return tt.clusters, nil
				},
				GetFn: func(ctx context.Context, identifier string) (*govpsie.K8s, error) {
					if tt.getResult != nil {
						return tt.getResult, nil
					}
					return nil, fmt.Errorf("not found")
				},
			}

			r := &kubernetesResource{client: mock}
			k8s, found, err := r.checkResourceStatus(t.Context(), tt.clusterName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tt.expectFound {
				t.Fatalf("expected found=%v, got %v", tt.expectFound, found)
			}
			if tt.expectFound && k8s == nil {
				t.Fatal("expected k8s to be non-nil when found")
			}
			if tt.expectFound && k8s.ClusterName != tt.clusterName {
				t.Fatalf("expected cluster name %q, got %q", tt.clusterName, k8s.ClusterName)
			}
		})
	}
}

func TestUnitKubernetesAPI_Get(t *testing.T) {
	expectedK8s := &govpsie.K8s{
		ClusterName: "test-cluster",
		Identifier:  "test-id-123",
		Cpu:         4,
		Ram:         8192,
	}

	mock := &mockKubernetesAPI{
		GetFn: func(ctx context.Context, identifier string) (*govpsie.K8s, error) {
			if identifier != "test-id-123" {
				t.Fatalf("expected identifier 'test-id-123', got %q", identifier)
			}
			return expectedK8s, nil
		},
	}

	k8s, err := mock.Get(t.Context(), "test-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k8s.ClusterName != "test-cluster" {
		t.Fatalf("expected cluster name 'test-cluster', got %q", k8s.ClusterName)
	}
	if k8s.Cpu != 4 {
		t.Fatalf("expected cpu 4, got %d", k8s.Cpu)
	}
}

func TestUnitKubernetesAPI_Delete(t *testing.T) {
	var calledWith struct {
		identifier string
		reason     string
		note       string
	}

	mock := &mockKubernetesAPI{
		DeleteFn: func(ctx context.Context, identifier, reason, note string) error {
			calledWith.identifier = identifier
			calledWith.reason = reason
			calledWith.note = note
			return nil
		},
	}

	err := mock.Delete(t.Context(), "cluster-id", "terraform", "terraform")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith.identifier != "cluster-id" {
		t.Fatalf("expected identifier 'cluster-id', got %q", calledWith.identifier)
	}
	if calledWith.reason != "terraform" {
		t.Fatalf("expected reason 'terraform', got %q", calledWith.reason)
	}
	if calledWith.note != "terraform" {
		t.Fatalf("expected note 'terraform', got %q", calledWith.note)
	}
}

func TestUnitKubernetesAPI_GetGroupByName(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
		clusters  []govpsie.ListK8s
		groups    map[string][]govpsie.K8sGroup
		expectErr bool
		expectID  int64
	}{
		{
			name:      "group found",
			groupName: "my-group",
			clusters: []govpsie.ListK8s{
				{Identifier: "cluster-1"},
			},
			groups: map[string][]govpsie.K8sGroup{
				"cluster-1": {
					{GroupName: "other-group", ID: 1},
					{GroupName: "my-group", ID: 2},
				},
			},
			expectErr: false,
			expectID:  2,
		},
		{
			name:      "group not found",
			groupName: "missing-group",
			clusters: []govpsie.ListK8s{
				{Identifier: "cluster-1"},
			},
			groups: map[string][]govpsie.K8sGroup{
				"cluster-1": {
					{GroupName: "other-group", ID: 1},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKubernetesAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error) {
					return tt.clusters, nil
				},
				ListK8sGroupsFn: func(ctx context.Context, identifier string) ([]govpsie.K8sGroup, error) {
					return tt.groups[identifier], nil
				},
			}

			r := &kubernetesGroupResource{client: mock}
			group, err := r.GetKubernetesGroupByName(t.Context(), tt.groupName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && group.ID != tt.expectID {
				t.Fatalf("expected group ID %d, got %d", tt.expectID, group.ID)
			}
		})
	}
}

func TestUnitKubernetesAPI_GetGroupByIdentifier(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		clusters   []govpsie.ListK8s
		groups     map[string][]govpsie.K8sGroup
		expectErr  bool
		expectName string
	}{
		{
			name:       "group found by identifier",
			identifier: "group-id-2",
			clusters: []govpsie.ListK8s{
				{Identifier: "cluster-1"},
			},
			groups: map[string][]govpsie.K8sGroup{
				"cluster-1": {
					{Identifier: "group-id-1", GroupName: "first"},
					{Identifier: "group-id-2", GroupName: "second"},
				},
			},
			expectErr:  false,
			expectName: "second",
		},
		{
			name:       "group not found by identifier",
			identifier: "missing-id",
			clusters: []govpsie.ListK8s{
				{Identifier: "cluster-1"},
			},
			groups: map[string][]govpsie.K8sGroup{
				"cluster-1": {
					{Identifier: "group-id-1", GroupName: "first"},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKubernetesAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error) {
					return tt.clusters, nil
				},
				ListK8sGroupsFn: func(ctx context.Context, identifier string) ([]govpsie.K8sGroup, error) {
					return tt.groups[identifier], nil
				},
			}

			r := &kubernetesGroupResource{client: mock}
			group, err := r.GetKubernetesGroupByIdentifier(t.Context(), tt.identifier)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && group.GroupName != tt.expectName {
				t.Fatalf("expected group name %q, got %q", tt.expectName, group.GroupName)
			}
		})
	}
}
