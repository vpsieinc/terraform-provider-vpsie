package firewall

import (
	"context"

	"github.com/vpsie/govpsie"
)

// FirewallAPI defines the subset of govpsie.FirewallGroupService methods
// used by the firewall resource, attachment resource, and data source in this provider.
type FirewallAPI interface {
	Create(ctx context.Context, groupName string, firewallUpdateReq []govpsie.FirewallUpdateReq) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.FirewallGroupListData, error)
	Get(ctx context.Context, fwGroupId string) (*govpsie.FirewallGroupDetailData, error)
	Delete(ctx context.Context, fwGroupId string) error
	AttachToVpsie(ctx context.Context, groupId, vmId string) error
	DetachFromVpsie(ctx context.Context, groupId, vmId string) error
}
