package gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockGatewayAPI implements GatewayAPI for unit testing.
type mockGatewayAPI struct {
	ListFn     func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Gateway, error)
	CreateFn   func(ctx context.Context, createReq *govpsie.CreateGatewayReq) error
	DeleteFn   func(ctx context.Context, ipId int) error
	GetFn      func(ctx context.Context, id int64) (*govpsie.Gateway, error)
	AttachVMFn func(ctx context.Context, id int64, vms []string, ignoreLegacyVms int64) error
	DetachVMFn func(ctx context.Context, id int64, mapping_id []int64) error
}

func (m *mockGatewayAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Gateway, error) {
	return m.ListFn(ctx, options)
}

func (m *mockGatewayAPI) Create(ctx context.Context, createReq *govpsie.CreateGatewayReq) error {
	return m.CreateFn(ctx, createReq)
}

func (m *mockGatewayAPI) Delete(ctx context.Context, ipId int) error {
	return m.DeleteFn(ctx, ipId)
}

func (m *mockGatewayAPI) Get(ctx context.Context, id int64) (*govpsie.Gateway, error) {
	return m.GetFn(ctx, id)
}

func (m *mockGatewayAPI) AttachVM(ctx context.Context, id int64, vms []string, ignoreLegacyVms int64) error {
	return m.AttachVMFn(ctx, id, vms, ignoreLegacyVms)
}

func (m *mockGatewayAPI) DetachVM(ctx context.Context, id int64, mapping_id []int64) error {
	return m.DetachVMFn(ctx, id, mapping_id)
}

// Compile-time check: mockGatewayAPI satisfies GatewayAPI.
var _ GatewayAPI = &mockGatewayAPI{}

func TestUnitGatewayAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockGatewayAPI{
		ListFn:     func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Gateway, error) { return nil, nil },
		CreateFn:   func(ctx context.Context, createReq *govpsie.CreateGatewayReq) error { return nil },
		DeleteFn:   func(ctx context.Context, ipId int) error { return nil },
		GetFn:      func(ctx context.Context, id int64) (*govpsie.Gateway, error) { return nil, nil },
		AttachVMFn: func(ctx context.Context, id int64, vms []string, ignoreLegacyVms int64) error { return nil },
		DetachVMFn: func(ctx context.Context, id int64, mapping_id []int64) error { return nil },
	}

	var api GatewayAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy GatewayAPI interface")
	}
}

func TestUnitGatewayAPI_CreateAndReturnGateway(t *testing.T) {
	tests := []struct {
		name        string
		existing    []govpsie.Gateway
		afterCreate []govpsie.Gateway
		createErr   error
		expectFound bool
		expectErr   bool
	}{
		{
			name:     "new gateway found after creation",
			existing: []govpsie.Gateway{{IP: "10.0.0.1"}},
			afterCreate: []govpsie.Gateway{
				{IP: "10.0.0.1"},
				{IP: "10.0.0.2", DcIdentifier: "dc-1"},
			},
			createErr:   nil,
			expectFound: true,
			expectErr:   false,
		},
		{
			name:        "create fails",
			existing:    []govpsie.Gateway{},
			afterCreate: []govpsie.Gateway{},
			createErr:   fmt.Errorf("api error"),
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "gateway not found after creation",
			existing:    []govpsie.Gateway{{IP: "10.0.0.1"}},
			afterCreate: []govpsie.Gateway{{IP: "10.0.0.1"}},
			createErr:   nil,
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			mock := &mockGatewayAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Gateway, error) {
					callCount++
					if callCount == 1 {
						return tt.existing, nil
					}
					return tt.afterCreate, nil
				},
				CreateFn: func(ctx context.Context, createReq *govpsie.CreateGatewayReq) error {
					return tt.createErr
				},
			}

			r := &gatewayResource{client: mock}
			gw, err := r.CreateAndReturnGateway(context.Background(), "ipv4", "dc-1")

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && gw == nil {
				t.Fatal("expected gateway to be non-nil when found")
			}
			if tt.expectFound && gw.IP != "10.0.0.2" {
				t.Fatalf("expected new gateway IP %q, got %q", "10.0.0.2", gw.IP)
			}
		})
	}
}

func TestUnitGatewayAPI_ListError(t *testing.T) {
	mock := &mockGatewayAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Gateway, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &gatewayResource{client: mock}
	_, err := r.CreateAndReturnGateway(context.Background(), "ipv4", "dc-1")
	if err == nil {
		t.Fatal("expected error from List failure, got nil")
	}
}
