package vpc

import (
	"context"

	"github.com/vpsie/govpsie"
)

// VpcAPI defines the subset of govpsie.VPCService methods
// used by the vpc, vpc_server_assignment resources and the vpc data source.
type VpcAPI interface {
	CreateVpc(ctx context.Context, createReq *govpsie.CreateVpcReq) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VPC, error)
	Get(ctx context.Context, id string) (*govpsie.VPC, error)
	DeleteVpc(ctx context.Context, vpcId, reason, note string) error
	AssignServer(ctx context.Context, assignReq *govpsie.AssignServerReq) error
	ReleasePrivateIP(ctx context.Context, vmIdentifer string, privateIpId int) error
}

// IPLookupAPI defines the subset of govpsie.IPService methods
// used by the vpc_server_assignment resource to look up private IPs.
type IPLookupAPI interface {
	ListPrivateIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
}
