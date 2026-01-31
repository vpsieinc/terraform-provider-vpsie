package image

import (
	"context"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockImageAPI implements ImageAPI for unit testing.
type mockImageAPI struct {
	CreateImagesFn func(ctx context.Context, dcIdentifier, imageName, imageUrl string) error
	GetImageFn     func(ctx context.Context, imageIdentifier string) (*govpsie.CustomImage, error)
	ListFn         func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.CustomImage, error)
	DeleteImageFn  func(ctx context.Context, imageIdentifier string) error
}

func (m *mockImageAPI) CreateImages(ctx context.Context, dcIdentifier, imageName, imageUrl string) error {
	return m.CreateImagesFn(ctx, dcIdentifier, imageName, imageUrl)
}

func (m *mockImageAPI) GetImage(ctx context.Context, imageIdentifier string) (*govpsie.CustomImage, error) {
	return m.GetImageFn(ctx, imageIdentifier)
}

func (m *mockImageAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.CustomImage, error) {
	return m.ListFn(ctx, options)
}

func (m *mockImageAPI) DeleteImage(ctx context.Context, imageIdentifier string) error {
	return m.DeleteImageFn(ctx, imageIdentifier)
}

// Compile-time check: mockImageAPI satisfies ImageAPI.
var _ ImageAPI = &mockImageAPI{}

func TestUnitImageAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockImageAPI{
		CreateImagesFn: func(ctx context.Context, dcIdentifier, imageName, imageUrl string) error {
			return nil
		},
		GetImageFn: func(ctx context.Context, imageIdentifier string) (*govpsie.CustomImage, error) {
			return &govpsie.CustomImage{}, nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.CustomImage, error) {
			return []govpsie.CustomImage{}, nil
		},
		DeleteImageFn: func(ctx context.Context, imageIdentifier string) error {
			return nil
		},
	}

	var api ImageAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitImageAPI_CheckResourceStatus(t *testing.T) {
	tests := []struct {
		name        string
		imageLabel  string
		images      []govpsie.CustomImage
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "image found by label",
			imageLabel: "test-image",
			images: []govpsie.CustomImage{
				{ImageLabel: "other-image", Identifier: "id-1"},
				{ImageLabel: "test-image", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "image not found",
			imageLabel: "missing-image",
			images: []govpsie.CustomImage{
				{ImageLabel: "other-image", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   false,
		},
		{
			name:        "empty image list",
			imageLabel:  "test-image",
			images:      []govpsie.CustomImage{},
			expectFound: false,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockImageAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.CustomImage, error) {
					return tt.images, nil
				},
			}

			r := &imageResource{client: mock}
			image, found, err := r.checkResourceStatus(t.Context(), tt.imageLabel)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tt.expectFound {
				t.Fatalf("expected found=%v, got %v", tt.expectFound, found)
			}
			if tt.expectFound && image == nil {
				t.Fatal("expected image to be non-nil when found")
			}
			if tt.expectFound && image.ImageLabel != tt.imageLabel {
				t.Fatalf("expected label %q, got %q", tt.imageLabel, image.ImageLabel)
			}
		})
	}
}

func TestUnitImageAPI_DeleteImage(t *testing.T) {
	var calledWith string

	mock := &mockImageAPI{
		DeleteImageFn: func(ctx context.Context, imageIdentifier string) error {
			calledWith = imageIdentifier
			return nil
		},
	}

	err := mock.DeleteImage(t.Context(), "image-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith != "image-id-123" {
		t.Fatalf("expected identifier 'image-id-123', got %q", calledWith)
	}
}
