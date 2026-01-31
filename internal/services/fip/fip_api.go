package fip

import (
	"context"

	"github.com/vpsie/govpsie"
)

// FipAPI defines the subset of govpsie.FipService methods
// used by the floating IP resource.
type FipAPI interface {
	CreateFloatingIP(ctx context.Context, vmIdentifier, dcIdentifier, ipType string) error
	UnassignFloatingIP(ctx context.Context, id string) error
}

// FipIPAPI defines the subset of govpsie.IPService methods
// used by the floating IP resource and data source.
type FipIPAPI interface {
	ListAllIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
	ListPublicIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
}
