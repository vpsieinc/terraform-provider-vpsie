package bucket

import (
	"context"

	"github.com/vpsie/govpsie"
)

// BucketAPI defines the subset of govpsie.BucketService methods
// used by the bucket resource and data source in this provider.
type BucketAPI interface {
	Create(ctx context.Context, createReq *govpsie.CreateBucketReq) error
	Get(ctx context.Context, id string) (*govpsie.Bucket, error)
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Bucket, error)
	Delete(ctx context.Context, buckId, reason, note string) error
	ToggleFileListing(ctx context.Context, bucketId string, fileListing bool) (bool, error)
}
