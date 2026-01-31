package gateway

import (
	"context"

	"github.com/vpsie/govpsie"
)

// GatewayAPI defines the subset of govpsie.GatewayService methods
// used by the gateway resource and data source.
type GatewayAPI interface {
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Gateway, error)
	Create(ctx context.Context, createReq *govpsie.CreateGatewayReq) error
	Delete(ctx context.Context, ipId int) error
	Get(ctx context.Context, id int64) (*govpsie.Gateway, error)
	AttachVM(ctx context.Context, id int64, vms []string, ignoreLegacyVms int64) error
	DetachVM(ctx context.Context, id int64, mapping_id []int64) error
}
