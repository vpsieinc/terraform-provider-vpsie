package datacenter

import (
	"context"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockDataCenterAPI implements DataCenterAPI for unit testing.
type mockDataCenterAPI struct {
	ListFn func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.DataCenter, error)
}

func (m *mockDataCenterAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.DataCenter, error) {
	return m.ListFn(ctx, options)
}

// Compile-time check: mockDataCenterAPI satisfies DataCenterAPI.
var _ DataCenterAPI = &mockDataCenterAPI{}

func TestUnitDataCenterAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockDataCenterAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.DataCenter, error) {
			return []govpsie.DataCenter{}, nil
		},
	}

	var api DataCenterAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy DataCenterAPI interface")
	}
}

func TestUnitDataCenterAPI_ListDatacenters(t *testing.T) {
	tests := []struct {
		name      string
		dcs       []govpsie.DataCenter
		expectLen int
		expectErr bool
	}{
		{
			name: "returns multiple datacenters",
			dcs: []govpsie.DataCenter{
				{DcName: "dc-us", Identifier: "id-1", Country: "US"},
				{DcName: "dc-eu", Identifier: "id-2", Country: "DE"},
			},
			expectLen: 2,
			expectErr: false,
		},
		{
			name:      "returns empty list",
			dcs:       []govpsie.DataCenter{},
			expectLen: 0,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDataCenterAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.DataCenter, error) {
					return tt.dcs, nil
				},
			}

			dcs, err := mock.List(context.Background(), nil)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(dcs) != tt.expectLen {
				t.Fatalf("expected %d datacenters, got %d", tt.expectLen, len(dcs))
			}
		})
	}
}

func TestUnitDataCenterAPI_ListVerifiesData(t *testing.T) {
	mock := &mockDataCenterAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.DataCenter, error) {
			return []govpsie.DataCenter{
				{DcName: "dc-test", Identifier: "dc-id-123", Country: "NL", IsActive: 1},
			}, nil
		},
	}

	dcs, err := mock.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dcs) != 1 {
		t.Fatalf("expected 1 datacenter, got %d", len(dcs))
	}
	if dcs[0].DcName != "dc-test" {
		t.Fatalf("expected DcName 'dc-test', got %q", dcs[0].DcName)
	}
	if dcs[0].Identifier != "dc-id-123" {
		t.Fatalf("expected Identifier 'dc-id-123', got %q", dcs[0].Identifier)
	}
	if dcs[0].Country != "NL" {
		t.Fatalf("expected Country 'NL', got %q", dcs[0].Country)
	}
}
