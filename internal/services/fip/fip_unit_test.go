package fip

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockFipAPI implements FipAPI for unit testing.
type mockFipAPI struct {
	CreateFloatingIPFn   func(ctx context.Context, vmIdentifier, dcIdentifier, ipType string) error
	UnassignFloatingIPFn func(ctx context.Context, id string) error
}

func (m *mockFipAPI) CreateFloatingIP(ctx context.Context, vmIdentifier, dcIdentifier, ipType string) error {
	return m.CreateFloatingIPFn(ctx, vmIdentifier, dcIdentifier, ipType)
}

func (m *mockFipAPI) UnassignFloatingIP(ctx context.Context, id string) error {
	return m.UnassignFloatingIPFn(ctx, id)
}

// mockFipIPAPI implements FipIPAPI for unit testing.
type mockFipIPAPI struct {
	ListAllIPsFn    func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
	ListPublicIPsFn func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
}

func (m *mockFipIPAPI) ListAllIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
	return m.ListAllIPsFn(ctx, options)
}

func (m *mockFipIPAPI) ListPublicIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
	return m.ListPublicIPsFn(ctx, options)
}

// Compile-time checks: mocks satisfy interfaces.
var _ FipAPI = &mockFipAPI{}
var _ FipIPAPI = &mockFipIPAPI{}

func TestUnitFipAPI_MockSatisfiesInterface(t *testing.T) {
	fipMock := &mockFipAPI{
		CreateFloatingIPFn:   func(ctx context.Context, vmIdentifier, dcIdentifier, ipType string) error { return nil },
		UnassignFloatingIPFn: func(ctx context.Context, id string) error { return nil },
	}

	var fipAPI FipAPI = fipMock
	if fipAPI == nil {
		t.Fatal("expected mock to satisfy FipAPI interface")
	}

	ipMock := &mockFipIPAPI{
		ListAllIPsFn:    func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) { return nil, nil },
		ListPublicIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) { return nil, nil },
	}

	var ipAPI FipIPAPI = ipMock
	if ipAPI == nil {
		t.Fatal("expected mock to satisfy FipIPAPI interface")
	}
}

func TestUnitFipAPI_CreateFloatingIP(t *testing.T) {
	tests := []struct {
		name      string
		createErr error
		expectErr bool
	}{
		{
			name:      "create succeeds",
			createErr: nil,
			expectErr: false,
		},
		{
			name:      "create fails",
			createErr: fmt.Errorf("api error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFipAPI{
				CreateFloatingIPFn: func(ctx context.Context, vmIdentifier, dcIdentifier, ipType string) error {
					return tt.createErr
				},
			}

			err := mock.CreateFloatingIP(context.Background(), "vm-1", "dc-1", "ipv4")
			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUnitFipAPI_UnassignFloatingIP(t *testing.T) {
	tests := []struct {
		name      string
		unassignErr error
		expectErr   bool
	}{
		{
			name:        "unassign succeeds",
			unassignErr: nil,
			expectErr:   false,
		},
		{
			name:        "unassign fails",
			unassignErr: fmt.Errorf("api error"),
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFipAPI{
				UnassignFloatingIPFn: func(ctx context.Context, id string) error {
					return tt.unassignErr
				},
			}

			err := mock.UnassignFloatingIP(context.Background(), "ip-1")
			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUnitFipIPAPI_ListAllIPs(t *testing.T) {
	tests := []struct {
		name      string
		ips       []govpsie.IP
		listErr   error
		expectLen int
		expectErr bool
	}{
		{
			name: "returns IPs",
			ips: []govpsie.IP{
				{IP: "10.0.0.1", ID: 1},
				{IP: "10.0.0.2", ID: 2},
			},
			listErr:   nil,
			expectLen: 2,
			expectErr: false,
		},
		{
			name:      "list fails",
			ips:       nil,
			listErr:   fmt.Errorf("api error"),
			expectLen: 0,
			expectErr: true,
		},
		{
			name:      "empty list",
			ips:       []govpsie.IP{},
			listErr:   nil,
			expectLen: 0,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockFipIPAPI{
				ListAllIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
					return tt.ips, tt.listErr
				},
			}

			ips, err := mock.ListAllIPs(context.Background(), nil)
			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && len(ips) != tt.expectLen {
				t.Fatalf("expected %d IPs, got %d", tt.expectLen, len(ips))
			}
		})
	}
}
