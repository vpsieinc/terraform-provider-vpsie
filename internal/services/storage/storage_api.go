package storage

import (
	"context"

	"github.com/vpsie/govpsie"
)

// StorageAPI defines the subset of govpsie.StorageService methods
// used by the storage resource and data source in this provider.
type StorageAPI interface {
	CreateVolume(ctx context.Context, req *govpsie.StorageCreateRequest) error
	ListAll(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Storage, error)
	Delete(ctx context.Context, storageIdentifier string) error
	UpdateName(ctx context.Context, storageIdentifier, name string) error
	UpdateSize(ctx context.Context, storageIdentifier, size string) error
}
