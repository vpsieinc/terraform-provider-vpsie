package monitoring

import (
	"context"

	"github.com/vpsie/govpsie"
)

// MonitoringAPI defines the subset of govpsie.MonitoringService methods
// used by the monitoring_rule resource and data source.
type MonitoringAPI interface {
	CreateRule(ctx context.Context, createReq *govpsie.CreateMonitoringRuleReq) error
	ListMonitoringRule(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.MonitoringRule, error)
	ToggleMonitoringRuleStatus(ctx context.Context, status, ruleIdentifier string) error
	DeleteMonitoringRule(ctx context.Context, ruleIdentifier string) error
}
