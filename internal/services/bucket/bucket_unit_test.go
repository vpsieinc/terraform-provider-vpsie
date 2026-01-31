package bucket

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockBucketAPI implements BucketAPI for unit testing.
type mockBucketAPI struct {
	CreateFn            func(ctx context.Context, createReq *govpsie.CreateBucketReq) error
	GetFn               func(ctx context.Context, id string) (*govpsie.Bucket, error)
	ListFn              func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Bucket, error)
	DeleteFn            func(ctx context.Context, buckId, reason, note string) error
	ToggleFileListingFn func(ctx context.Context, bucketId string, fileListing bool) (bool, error)
}

func (m *mockBucketAPI) Create(ctx context.Context, createReq *govpsie.CreateBucketReq) error {
	return m.CreateFn(ctx, createReq)
}

func (m *mockBucketAPI) Get(ctx context.Context, id string) (*govpsie.Bucket, error) {
	return m.GetFn(ctx, id)
}

func (m *mockBucketAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Bucket, error) {
	return m.ListFn(ctx, options)
}

func (m *mockBucketAPI) Delete(ctx context.Context, buckId, reason, note string) error {
	return m.DeleteFn(ctx, buckId, reason, note)
}

func (m *mockBucketAPI) ToggleFileListing(ctx context.Context, bucketId string, fileListing bool) (bool, error) {
	return m.ToggleFileListingFn(ctx, bucketId, fileListing)
}

// Compile-time check: mockBucketAPI satisfies BucketAPI.
var _ BucketAPI = &mockBucketAPI{}

func TestUnitBucketAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockBucketAPI{
		CreateFn: func(ctx context.Context, createReq *govpsie.CreateBucketReq) error {
			return nil
		},
		GetFn: func(ctx context.Context, id string) (*govpsie.Bucket, error) {
			return &govpsie.Bucket{}, nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Bucket, error) {
			return []govpsie.Bucket{}, nil
		},
		DeleteFn: func(ctx context.Context, buckId, reason, note string) error {
			return nil
		},
		ToggleFileListingFn: func(ctx context.Context, bucketId string, fileListing bool) (bool, error) {
			return true, nil
		},
	}

	var api BucketAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitBucketAPI_GetBucketByName(t *testing.T) {
	tests := []struct {
		name       string
		bucketName string
		buckets    []govpsie.Bucket
		expectErr  bool
	}{
		{
			name:       "bucket found by name",
			bucketName: "test-bucket",
			buckets: []govpsie.Bucket{
				{BucketName: "other-bucket", Identifier: "id-1"},
				{BucketName: "test-bucket", Identifier: "id-2"},
			},
			expectErr: false,
		},
		{
			name:       "bucket not found",
			bucketName: "missing-bucket",
			buckets: []govpsie.Bucket{
				{BucketName: "other-bucket", Identifier: "id-1"},
			},
			expectErr: true,
		},
		{
			name:       "empty bucket list",
			bucketName: "test-bucket",
			buckets:    []govpsie.Bucket{},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockBucketAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Bucket, error) {
					return tt.buckets, nil
				},
			}

			r := &bucketResource{client: mock}
			bucket, err := r.GetBucketByName(t.Context(), tt.bucketName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && bucket == nil {
				t.Fatal("expected bucket to be non-nil")
			}
			if !tt.expectErr && bucket.BucketName != tt.bucketName {
				t.Fatalf("expected name %q, got %q", tt.bucketName, bucket.BucketName)
			}
		})
	}
}

func TestUnitBucketAPI_Delete(t *testing.T) {
	var calledWith struct {
		buckId string
		reason string
		note   string
	}

	mock := &mockBucketAPI{
		DeleteFn: func(ctx context.Context, buckId, reason, note string) error {
			calledWith.buckId = buckId
			calledWith.reason = reason
			calledWith.note = note
			return nil
		},
	}

	err := mock.Delete(t.Context(), "bucket-id-123", "test-reason", "test-note")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith.buckId != "bucket-id-123" {
		t.Fatalf("expected buckId 'bucket-id-123', got %q", calledWith.buckId)
	}
	if calledWith.reason != "test-reason" {
		t.Fatalf("expected reason 'test-reason', got %q", calledWith.reason)
	}
}

func TestUnitBucketAPI_ListError(t *testing.T) {
	mock := &mockBucketAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Bucket, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &bucketResource{client: mock}
	_, err := r.GetBucketByName(t.Context(), "any-bucket")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "api error" {
		t.Fatalf("expected 'api error', got %q", err.Error())
	}
}
