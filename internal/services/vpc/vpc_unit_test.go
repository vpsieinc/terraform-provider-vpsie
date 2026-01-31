package vpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockVpcAPI implements VpcAPI for unit testing.
type mockVpcAPI struct {
	CreateVpcFn        func(ctx context.Context, createReq *govpsie.CreateVpcReq) error
	ListFn             func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VPC, error)
	GetFn              func(ctx context.Context, id string) (*govpsie.VPC, error)
	DeleteVpcFn        func(ctx context.Context, vpcId, reason, note string) error
	AssignServerFn     func(ctx context.Context, assignReq *govpsie.AssignServerReq) error
	ReleasePrivateIPFn func(ctx context.Context, vmIdentifer string, privateIpId int) error
}

func (m *mockVpcAPI) CreateVpc(ctx context.Context, createReq *govpsie.CreateVpcReq) error {
	return m.CreateVpcFn(ctx, createReq)
}

func (m *mockVpcAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VPC, error) {
	return m.ListFn(ctx, options)
}

func (m *mockVpcAPI) Get(ctx context.Context, id string) (*govpsie.VPC, error) {
	return m.GetFn(ctx, id)
}

func (m *mockVpcAPI) DeleteVpc(ctx context.Context, vpcId, reason, note string) error {
	return m.DeleteVpcFn(ctx, vpcId, reason, note)
}

func (m *mockVpcAPI) AssignServer(ctx context.Context, assignReq *govpsie.AssignServerReq) error {
	return m.AssignServerFn(ctx, assignReq)
}

func (m *mockVpcAPI) ReleasePrivateIP(ctx context.Context, vmIdentifer string, privateIpId int) error {
	return m.ReleasePrivateIPFn(ctx, vmIdentifer, privateIpId)
}

// mockIPLookupAPI implements IPLookupAPI for unit testing.
type mockIPLookupAPI struct {
	ListPrivateIPsFn func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
}

func (m *mockIPLookupAPI) ListPrivateIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
	return m.ListPrivateIPsFn(ctx, options)
}

// Compile-time checks: mocks satisfy interfaces.
var _ VpcAPI = &mockVpcAPI{}
var _ IPLookupAPI = &mockIPLookupAPI{}

func TestUnitVpcAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockVpcAPI{
		CreateVpcFn:        func(ctx context.Context, createReq *govpsie.CreateVpcReq) error { return nil },
		ListFn:             func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VPC, error) { return nil, nil },
		GetFn:              func(ctx context.Context, id string) (*govpsie.VPC, error) { return nil, nil },
		DeleteVpcFn:        func(ctx context.Context, vpcId, reason, note string) error { return nil },
		AssignServerFn:     func(ctx context.Context, assignReq *govpsie.AssignServerReq) error { return nil },
		ReleasePrivateIPFn: func(ctx context.Context, vmIdentifer string, privateIpId int) error { return nil },
	}

	var api VpcAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitVpcAPI_GetVpcByName(t *testing.T) {
	tests := []struct {
		name        string
		vpcName     string
		vpcs        []govpsie.VPC
		expectFound bool
		expectErr   bool
	}{
		{
			name:    "vpc found by name",
			vpcName: "test-vpc",
			vpcs: []govpsie.VPC{
				{Name: "other-vpc"},
				{Name: "test-vpc"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:    "vpc not found",
			vpcName: "missing-vpc",
			vpcs: []govpsie.VPC{
				{Name: "other-vpc"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty vpc list",
			vpcName:     "test-vpc",
			vpcs:        []govpsie.VPC{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockVpcAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VPC, error) {
					return tt.vpcs, nil
				},
			}

			r := &vpcResource{client: mock}
			vpc, err := r.GetVpcByName(t.Context(), tt.vpcName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && vpc == nil {
				t.Fatal("expected vpc to be non-nil when found")
			}
			if tt.expectFound && vpc.Name != tt.vpcName {
				t.Fatalf("expected name %q, got %q", tt.vpcName, vpc.Name)
			}
		})
	}
}

func TestUnitVpcAPI_ListError(t *testing.T) {
	mock := &mockVpcAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VPC, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &vpcResource{client: mock}
	_, err := r.GetVpcByName(t.Context(), "test")
	if err == nil {
		t.Fatal("expected error from List failure, got nil")
	}
}
