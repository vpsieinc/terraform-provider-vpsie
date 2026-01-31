package loadbalancer

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockLoadbalancerAPI implements LoadbalancerAPI for unit testing.
type mockLoadbalancerAPI struct {
	ListLBsFn            func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.LB, error)
	GetLBFn              func(ctx context.Context, lbID string) (*govpsie.LBDetails, error)
	CreateLBFn           func(ctx context.Context, createLBReq *govpsie.CreateLBReq) error
	DeleteLBFn           func(ctx context.Context, lbID, reason, note string) error
	AddLBRuleFn          func(ctx context.Context, addRuleReq *govpsie.AddRuleReq) error
	DeleteLBRuleFn       func(ctx context.Context, ruleID string) error
	UpdateLBDomainFn     func(ctx context.Context, domainUpdateReq *govpsie.DomainUpdateReq) error
	UpdateDomainBackendFn func(ctx context.Context, domainId string, backends []govpsie.Backend) error
	UpdateLBRulesFn      func(ctx context.Context, ruleUpdateReq *govpsie.RuleUpdateReq) error
}

func (m *mockLoadbalancerAPI) ListLBs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.LB, error) {
	return m.ListLBsFn(ctx, options)
}

func (m *mockLoadbalancerAPI) GetLB(ctx context.Context, lbID string) (*govpsie.LBDetails, error) {
	return m.GetLBFn(ctx, lbID)
}

func (m *mockLoadbalancerAPI) CreateLB(ctx context.Context, createLBReq *govpsie.CreateLBReq) error {
	return m.CreateLBFn(ctx, createLBReq)
}

func (m *mockLoadbalancerAPI) DeleteLB(ctx context.Context, lbID, reason, note string) error {
	return m.DeleteLBFn(ctx, lbID, reason, note)
}

func (m *mockLoadbalancerAPI) AddLBRule(ctx context.Context, addRuleReq *govpsie.AddRuleReq) error {
	return m.AddLBRuleFn(ctx, addRuleReq)
}

func (m *mockLoadbalancerAPI) DeleteLBRule(ctx context.Context, ruleID string) error {
	return m.DeleteLBRuleFn(ctx, ruleID)
}

func (m *mockLoadbalancerAPI) UpdateLBDomain(ctx context.Context, domainUpdateReq *govpsie.DomainUpdateReq) error {
	return m.UpdateLBDomainFn(ctx, domainUpdateReq)
}

func (m *mockLoadbalancerAPI) UpdateDomainBackend(ctx context.Context, domainId string, backends []govpsie.Backend) error {
	return m.UpdateDomainBackendFn(ctx, domainId, backends)
}

func (m *mockLoadbalancerAPI) UpdateLBRules(ctx context.Context, ruleUpdateReq *govpsie.RuleUpdateReq) error {
	return m.UpdateLBRulesFn(ctx, ruleUpdateReq)
}

// Compile-time check: mockLoadbalancerAPI satisfies LoadbalancerAPI.
var _ LoadbalancerAPI = &mockLoadbalancerAPI{}

func TestUnitLoadbalancerAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockLoadbalancerAPI{
		ListLBsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.LB, error) {
			return []govpsie.LB{}, nil
		},
		GetLBFn: func(ctx context.Context, lbID string) (*govpsie.LBDetails, error) {
			return &govpsie.LBDetails{}, nil
		},
		CreateLBFn: func(ctx context.Context, createLBReq *govpsie.CreateLBReq) error {
			return nil
		},
		DeleteLBFn: func(ctx context.Context, lbID, reason, note string) error {
			return nil
		},
		AddLBRuleFn: func(ctx context.Context, addRuleReq *govpsie.AddRuleReq) error {
			return nil
		},
		DeleteLBRuleFn: func(ctx context.Context, ruleID string) error {
			return nil
		},
		UpdateLBDomainFn: func(ctx context.Context, domainUpdateReq *govpsie.DomainUpdateReq) error {
			return nil
		},
		UpdateDomainBackendFn: func(ctx context.Context, domainId string, backends []govpsie.Backend) error {
			return nil
		},
		UpdateLBRulesFn: func(ctx context.Context, ruleUpdateReq *govpsie.RuleUpdateReq) error {
			return nil
		},
	}

	var api LoadbalancerAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy LoadbalancerAPI interface")
	}
}

func TestUnitLoadbalancerAPI_CheckResourceStatus(t *testing.T) {
	tests := []struct {
		name        string
		lbName      string
		lbs         []govpsie.LB
		getResult   *govpsie.LBDetails
		expectFound bool
		expectErr   bool
	}{
		{
			name:   "lb found by name",
			lbName: "test-lb",
			lbs: []govpsie.LB{
				{LBName: "other-lb", Identifier: "id-1"},
				{LBName: "test-lb", Identifier: "id-2"},
			},
			getResult: &govpsie.LBDetails{
				LBName:     "test-lb",
				Identifier: "id-2",
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:   "lb not found",
			lbName: "missing-lb",
			lbs: []govpsie.LB{
				{LBName: "other-lb", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   false,
		},
		{
			name:        "empty lb list",
			lbName:      "test-lb",
			lbs:         []govpsie.LB{},
			expectFound: false,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockLoadbalancerAPI{
				ListLBsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.LB, error) {
					return tt.lbs, nil
				},
				GetLBFn: func(ctx context.Context, lbID string) (*govpsie.LBDetails, error) {
					if tt.getResult != nil {
						return tt.getResult, nil
					}
					return nil, fmt.Errorf("not found")
				},
			}

			r := &loadbalancerResource{client: mock}
			lb, found, err := r.checkResourceStatus(context.Background(), tt.lbName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tt.expectFound {
				t.Fatalf("expected found=%v, got %v", tt.expectFound, found)
			}
			if tt.expectFound && lb == nil {
				t.Fatal("expected lb to be non-nil when found")
			}
			if tt.expectFound && lb.LBName != tt.lbName {
				t.Fatalf("expected lb name %q, got %q", tt.lbName, lb.LBName)
			}
		})
	}
}

func TestUnitLoadbalancerAPI_GetLB(t *testing.T) {
	expectedLB := &govpsie.LBDetails{
		LBName:     "test-lb",
		Identifier: "test-id-123",
		Traffic:    100,
		DefaultIP:  "10.0.0.1",
	}

	mock := &mockLoadbalancerAPI{
		GetLBFn: func(ctx context.Context, lbID string) (*govpsie.LBDetails, error) {
			if lbID != "test-id-123" {
				t.Fatalf("expected lbID 'test-id-123', got %q", lbID)
			}
			return expectedLB, nil
		},
	}

	lb, err := mock.GetLB(context.Background(), "test-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lb.LBName != "test-lb" {
		t.Fatalf("expected lb name 'test-lb', got %q", lb.LBName)
	}
	if lb.Traffic != 100 {
		t.Fatalf("expected traffic 100, got %d", lb.Traffic)
	}
	if lb.DefaultIP != "10.0.0.1" {
		t.Fatalf("expected default IP '10.0.0.1', got %q", lb.DefaultIP)
	}
}

func TestUnitLoadbalancerAPI_DeleteLB(t *testing.T) {
	var calledWith struct {
		lbID   string
		reason string
		note   string
	}

	mock := &mockLoadbalancerAPI{
		DeleteLBFn: func(ctx context.Context, lbID, reason, note string) error {
			calledWith.lbID = lbID
			calledWith.reason = reason
			calledWith.note = note
			return nil
		},
	}

	err := mock.DeleteLB(context.Background(), "lb-id-1", "terraform provider", "terraform provider")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith.lbID != "lb-id-1" {
		t.Fatalf("expected lbID 'lb-id-1', got %q", calledWith.lbID)
	}
	if calledWith.reason != "terraform provider" {
		t.Fatalf("expected reason 'terraform provider', got %q", calledWith.reason)
	}
	if calledWith.note != "terraform provider" {
		t.Fatalf("expected note 'terraform provider', got %q", calledWith.note)
	}
}

func TestUnitLoadbalancerAPI_ListLBsError(t *testing.T) {
	mock := &mockLoadbalancerAPI{
		ListLBsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.LB, error) {
			return nil, fmt.Errorf("API error: connection refused")
		},
	}

	r := &loadbalancerResource{client: mock}
	_, found, err := r.checkResourceStatus(context.Background(), "any-lb")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if found {
		t.Fatal("expected found=false when error occurs")
	}
	if err.Error() != "API error: connection refused" {
		t.Fatalf("expected specific error message, got %q", err.Error())
	}
}
