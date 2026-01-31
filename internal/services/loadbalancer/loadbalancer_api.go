package loadbalancer

import (
	"context"

	"github.com/vpsie/govpsie"
)

// LoadbalancerAPI defines the subset of govpsie.LBsService methods
// used by the loadbalancer resource and data source in this provider.
type LoadbalancerAPI interface {
	ListLBs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.LB, error)
	GetLB(ctx context.Context, lbID string) (*govpsie.LBDetails, error)
	CreateLB(ctx context.Context, createLBReq *govpsie.CreateLBReq) error
	DeleteLB(ctx context.Context, lbID, reason, note string) error
	AddLBRule(ctx context.Context, addRuleReq *govpsie.AddRuleReq) error
	DeleteLBRule(ctx context.Context, ruleID string) error
	UpdateLBDomain(ctx context.Context, domainUpdateReq *govpsie.DomainUpdateReq) error
	UpdateDomainBackend(ctx context.Context, domainId string, backends []govpsie.Backend) error
	UpdateLBRules(ctx context.Context, ruleUpdateReq *govpsie.RuleUpdateReq) error
}
