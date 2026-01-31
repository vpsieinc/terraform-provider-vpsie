package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockStorageAPI implements StorageAPI for unit testing.
type mockStorageAPI struct {
	CreateVolumeFn func(ctx context.Context, req *govpsie.StorageCreateRequest) error
	ListAllFn      func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error)
	DeleteFn       func(ctx context.Context, storageIdentifier string) error
	UpdateNameFn   func(ctx context.Context, storageIdentifier, name string) error
	UpdateSizeFn   func(ctx context.Context, storageIdentifier, size string) error
}

func (m *mockStorageAPI) CreateVolume(ctx context.Context, req *govpsie.StorageCreateRequest) error {
	return m.CreateVolumeFn(ctx, req)
}

func (m *mockStorageAPI) ListAll(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error) {
	return m.ListAllFn(ctx, options)
}

func (m *mockStorageAPI) Delete(ctx context.Context, storageIdentifier string) error {
	return m.DeleteFn(ctx, storageIdentifier)
}

func (m *mockStorageAPI) UpdateName(ctx context.Context, storageIdentifier, name string) error {
	return m.UpdateNameFn(ctx, storageIdentifier, name)
}

func (m *mockStorageAPI) UpdateSize(ctx context.Context, storageIdentifier, size string) error {
	return m.UpdateSizeFn(ctx, storageIdentifier, size)
}

// Compile-time check: mockStorageAPI satisfies StorageAPI.
var _ StorageAPI = &mockStorageAPI{}

func TestUnitStorageAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockStorageAPI{
		CreateVolumeFn: func(ctx context.Context, req *govpsie.StorageCreateRequest) error {
			return nil
		},
		ListAllFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error) {
			return []govpsie.Storage{}, nil
		},
		DeleteFn: func(ctx context.Context, storageIdentifier string) error {
			return nil
		},
		UpdateNameFn: func(ctx context.Context, storageIdentifier, name string) error {
			return nil
		},
		UpdateSizeFn: func(ctx context.Context, storageIdentifier, size string) error {
			return nil
		},
	}

	var api StorageAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitStorageAPI_GetVolumeByName(t *testing.T) {
	tests := []struct {
		name        string
		volumeName  string
		volumes     []govpsie.Storage
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "volume found by name",
			volumeName: "test-vol",
			volumes: []govpsie.Storage{
				{Name: "other-vol", Identifier: "id-1"},
				{Name: "test-vol", Identifier: "id-2", ID: 10},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "volume not found",
			volumeName: "missing-vol",
			volumes: []govpsie.Storage{
				{Name: "other-vol", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty volume list",
			volumeName:  "test-vol",
			volumes:     []govpsie.Storage{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockStorageAPI{
				ListAllFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error) {
					return tt.volumes, nil
				},
			}

			r := &storageResource{client: mock}
			vol, err := r.GetVolumeByName(t.Context(), tt.volumeName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && vol == nil {
				t.Fatal("expected volume to be non-nil when found")
			}
			if tt.expectFound && vol.Name != tt.volumeName {
				t.Fatalf("expected name %q, got %q", tt.volumeName, vol.Name)
			}
		})
	}
}

func TestUnitStorageAPI_GetVolumeByIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		identifier  string
		volumes     []govpsie.Storage
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "volume found by identifier",
			identifier: "id-2",
			volumes: []govpsie.Storage{
				{Name: "vol-1", Identifier: "id-1"},
				{Name: "vol-2", Identifier: "id-2", ID: 20},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "volume not found by identifier",
			identifier: "id-missing",
			volumes: []govpsie.Storage{
				{Name: "vol-1", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockStorageAPI{
				ListAllFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error) {
					return tt.volumes, nil
				},
			}

			r := &storageResource{client: mock}
			vol, err := r.GetVolumeByIdentifier(t.Context(), tt.identifier)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && vol == nil {
				t.Fatal("expected volume to be non-nil when found")
			}
			if tt.expectFound && vol.Identifier != tt.identifier {
				t.Fatalf("expected identifier %q, got %q", tt.identifier, vol.Identifier)
			}
		})
	}
}

func TestUnitStorageAPI_ListAllError(t *testing.T) {
	mock := &mockStorageAPI{
		ListAllFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &storageResource{client: mock}
	_, err := r.GetVolumeByName(t.Context(), "test")
	if err == nil {
		t.Fatal("expected error from ListAll failure, got nil")
	}

	_, err = r.GetVolumeByIdentifier(t.Context(), "test-id")
	if err == nil {
		t.Fatal("expected error from ListAll failure, got nil")
	}
}
