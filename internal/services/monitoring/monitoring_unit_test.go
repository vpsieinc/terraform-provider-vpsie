package monitoring

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockMonitoringAPI implements MonitoringAPI for unit testing.
type mockMonitoringAPI struct {
	CreateRuleFn                  func(ctx context.Context, createReq *govpsie.CreateMonitoringRuleReq) error
	ListMonitoringRuleFn          func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.MonitoringRule, error)
	ToggleMonitoringRuleStatusFn  func(ctx context.Context, status, ruleIdentifier string) error
	DeleteMonitoringRuleFn        func(ctx context.Context, ruleIdentifier string) error
}

func (m *mockMonitoringAPI) CreateRule(ctx context.Context, createReq *govpsie.CreateMonitoringRuleReq) error {
	return m.CreateRuleFn(ctx, createReq)
}

func (m *mockMonitoringAPI) ListMonitoringRule(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.MonitoringRule, error) {
	return m.ListMonitoringRuleFn(ctx, options)
}

func (m *mockMonitoringAPI) ToggleMonitoringRuleStatus(ctx context.Context, status, ruleIdentifier string) error {
	return m.ToggleMonitoringRuleStatusFn(ctx, status, ruleIdentifier)
}

func (m *mockMonitoringAPI) DeleteMonitoringRule(ctx context.Context, ruleIdentifier string) error {
	return m.DeleteMonitoringRuleFn(ctx, ruleIdentifier)
}

// Compile-time check: mockMonitoringAPI satisfies MonitoringAPI.
var _ MonitoringAPI = &mockMonitoringAPI{}

func TestUnitMonitoringAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockMonitoringAPI{
		CreateRuleFn:                 func(ctx context.Context, createReq *govpsie.CreateMonitoringRuleReq) error { return nil },
		ListMonitoringRuleFn:         func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.MonitoringRule, error) { return nil, nil },
		ToggleMonitoringRuleStatusFn: func(ctx context.Context, status, ruleIdentifier string) error { return nil },
		DeleteMonitoringRuleFn:       func(ctx context.Context, ruleIdentifier string) error { return nil },
	}

	var api MonitoringAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy MonitoringAPI interface")
	}
}

func TestUnitMonitoringAPI_GetRuleByName(t *testing.T) {
	tests := []struct {
		name        string
		ruleName    string
		rules       []govpsie.MonitoringRule
		expectFound bool
		expectErr   bool
	}{
		{
			name:     "rule found by name",
			ruleName: "cpu-alert",
			rules: []govpsie.MonitoringRule{
				{RuleName: "memory-alert", Identifier: "id-1"},
				{RuleName: "cpu-alert", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:     "rule not found",
			ruleName: "missing-rule",
			rules: []govpsie.MonitoringRule{
				{RuleName: "cpu-alert", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty rule list",
			ruleName:    "cpu-alert",
			rules:       []govpsie.MonitoringRule{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockMonitoringAPI{
				ListMonitoringRuleFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.MonitoringRule, error) {
					return tt.rules, nil
				},
			}

			r := &monitoringRuleResource{client: mock}
			rule, err := r.GetRuleByName(context.Background(), tt.ruleName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && rule == nil {
				t.Fatal("expected rule to be non-nil when found")
			}
			if tt.expectFound && rule.RuleName != tt.ruleName {
				t.Fatalf("expected rule name %q, got %q", tt.ruleName, rule.RuleName)
			}
		})
	}
}

func TestUnitMonitoringAPI_ListMonitoringRuleError(t *testing.T) {
	mock := &mockMonitoringAPI{
		ListMonitoringRuleFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.MonitoringRule, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &monitoringRuleResource{client: mock}
	_, err := r.GetRuleByName(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error from ListMonitoringRule failure, got nil")
	}
}
