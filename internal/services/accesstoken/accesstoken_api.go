package accesstoken

import (
	"context"

	"github.com/vpsie/govpsie"
)

// AccessTokenAPI defines the subset of govpsie.AccessTokenService methods
// used by the access token resource and data source in this provider.
type AccessTokenAPI interface {
	Create(ctx context.Context, name, accessToken, expirationDate string) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.AccessToken, error)
	Update(ctx context.Context, accessTokenIdentifier, name, expirationDate string) error
	Delete(ctx context.Context, accessTokenIdentifier string) error
}
