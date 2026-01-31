package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockDomainAPI implements DomainAPI for unit testing.
type mockDomainAPI struct {
	CreateDomainFn        func(ctx context.Context, createReq *govpsie.CreateDomainRequest) error
	ListDomainsFn         func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Domain, error)
	DeleteDomainFn        func(ctx context.Context, domainIdentifier, reason, note string) error
	CreateDnsRecordFn     func(ctx context.Context, createReq govpsie.CreateDnsRecordReq) error
	UpdateDnsRecordFn     func(ctx context.Context, updateReq *govpsie.UpdateDnsRecordReq) error
	DeleteDnsRecordFn     func(ctx context.Context, domainIdentifier string, record *govpsie.Record) error
	AddReverseFn          func(ctx context.Context, reverseReq *govpsie.ReverseRequest) error
	ListReversePTRRecordsFn func(ctx context.Context) ([]govpsie.ReversePTR, error)
	UpdateReverseFn       func(ctx context.Context, reverseReq *govpsie.ReverseRequest) error
	DeleteReverseFn       func(ctx context.Context, ip, vmIdentifier string) error
}

func (m *mockDomainAPI) CreateDomain(ctx context.Context, createReq *govpsie.CreateDomainRequest) error {
	return m.CreateDomainFn(ctx, createReq)
}

func (m *mockDomainAPI) ListDomains(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Domain, error) {
	return m.ListDomainsFn(ctx, options)
}

func (m *mockDomainAPI) DeleteDomain(ctx context.Context, domainIdentifier, reason, note string) error {
	return m.DeleteDomainFn(ctx, domainIdentifier, reason, note)
}

func (m *mockDomainAPI) CreateDnsRecord(ctx context.Context, createReq govpsie.CreateDnsRecordReq) error {
	return m.CreateDnsRecordFn(ctx, createReq)
}

func (m *mockDomainAPI) UpdateDnsRecord(ctx context.Context, updateReq *govpsie.UpdateDnsRecordReq) error {
	return m.UpdateDnsRecordFn(ctx, updateReq)
}

func (m *mockDomainAPI) DeleteDnsRecord(ctx context.Context, domainIdentifier string, record *govpsie.Record) error {
	return m.DeleteDnsRecordFn(ctx, domainIdentifier, record)
}

func (m *mockDomainAPI) AddReverse(ctx context.Context, reverseReq *govpsie.ReverseRequest) error {
	return m.AddReverseFn(ctx, reverseReq)
}

func (m *mockDomainAPI) ListReversePTRRecords(ctx context.Context) ([]govpsie.ReversePTR, error) {
	return m.ListReversePTRRecordsFn(ctx)
}

func (m *mockDomainAPI) UpdateReverse(ctx context.Context, reverseReq *govpsie.ReverseRequest) error {
	return m.UpdateReverseFn(ctx, reverseReq)
}

func (m *mockDomainAPI) DeleteReverse(ctx context.Context, ip, vmIdentifier string) error {
	return m.DeleteReverseFn(ctx, ip, vmIdentifier)
}

// Compile-time check: mockDomainAPI satisfies DomainAPI.
var _ DomainAPI = &mockDomainAPI{}

func TestUnitDomainAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockDomainAPI{
		CreateDomainFn:    func(ctx context.Context, createReq *govpsie.CreateDomainRequest) error { return nil },
		ListDomainsFn:     func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Domain, error) { return nil, nil },
		DeleteDomainFn:    func(ctx context.Context, domainIdentifier, reason, note string) error { return nil },
		CreateDnsRecordFn: func(ctx context.Context, createReq govpsie.CreateDnsRecordReq) error { return nil },
		UpdateDnsRecordFn: func(ctx context.Context, updateReq *govpsie.UpdateDnsRecordReq) error { return nil },
		DeleteDnsRecordFn: func(ctx context.Context, domainIdentifier string, record *govpsie.Record) error { return nil },
		AddReverseFn:      func(ctx context.Context, reverseReq *govpsie.ReverseRequest) error { return nil },
		ListReversePTRRecordsFn: func(ctx context.Context) ([]govpsie.ReversePTR, error) { return nil, nil },
		UpdateReverseFn:   func(ctx context.Context, reverseReq *govpsie.ReverseRequest) error { return nil },
		DeleteReverseFn:   func(ctx context.Context, ip, vmIdentifier string) error { return nil },
	}

	var api DomainAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy DomainAPI interface")
	}
}

func TestUnitDomainAPI_GetDomainByName(t *testing.T) {
	tests := []struct {
		name        string
		domainName  string
		domains     []govpsie.Domain
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "domain found by name",
			domainName: "example.com",
			domains: []govpsie.Domain{
				{DomainName: "other.com", Identifier: "id-1"},
				{DomainName: "example.com", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "domain not found",
			domainName: "missing.com",
			domains: []govpsie.Domain{
				{DomainName: "other.com", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty domain list",
			domainName:  "example.com",
			domains:     []govpsie.Domain{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDomainAPI{
				ListDomainsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Domain, error) {
					return tt.domains, nil
				},
			}

			r := &domainResource{client: mock}
			domain, err := r.GetDomainByName(context.Background(), tt.domainName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && domain == nil {
				t.Fatal("expected domain to be non-nil when found")
			}
			if tt.expectFound && domain.DomainName != tt.domainName {
				t.Fatalf("expected domain name %q, got %q", tt.domainName, domain.DomainName)
			}
		})
	}
}

func TestUnitDomainAPI_ListDomainsError(t *testing.T) {
	mock := &mockDomainAPI{
		ListDomainsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Domain, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &domainResource{client: mock}
	_, err := r.GetDomainByName(context.Background(), "test.com")
	if err == nil {
		t.Fatal("expected error from ListDomains failure, got nil")
	}
}
