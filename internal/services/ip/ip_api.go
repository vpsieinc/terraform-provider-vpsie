package ip

import (
	"context"

	"github.com/vpsie/govpsie"
)

// IPAPI defines the subset of govpsie.IPsService methods
// used by the IP data source in this provider.
type IPAPI interface {
	ListPublicIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
	ListPrivateIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
	ListAllIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
}
