package image

import (
	"context"

	"github.com/vpsie/govpsie"
)

// ImageAPI defines the subset of govpsie.ImagesService methods
// used by the image resource and data source in this provider.
type ImageAPI interface {
	CreateImages(ctx context.Context, dcIdentifier, imageName, imageUrl string) error
	GetImage(ctx context.Context, imageIdentifier string) (*govpsie.CustomImage, error)
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.CustomImage, error)
	DeleteImage(ctx context.Context, imageIdentifier string) error
}
