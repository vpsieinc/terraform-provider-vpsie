package datacenter

import (
	"context"

	"github.com/vpsie/govpsie"
)

// DataCenterAPI defines the subset of govpsie.DataCenterService methods
// used by the datacenter data source in this provider.
type DataCenterAPI interface {
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.DataCenter, error)
}
