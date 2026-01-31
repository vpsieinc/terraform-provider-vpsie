package domain

import (
	"context"

	"github.com/vpsie/govpsie"
)

// DomainAPI defines the subset of govpsie.DomainService methods
// used by the domain, dns_record, and reverse_dns resources and the domain data source.
type DomainAPI interface {
	CreateDomain(ctx context.Context, createReq *govpsie.CreateDomainRequest) error
	ListDomains(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Domain, error)
	DeleteDomain(ctx context.Context, domainIdentifier, reason, note string) error
	CreateDnsRecord(ctx context.Context, createReq govpsie.CreateDnsRecordReq) error
	UpdateDnsRecord(ctx context.Context, updateReq *govpsie.UpdateDnsRecordReq) error
	DeleteDnsRecord(ctx context.Context, domainIdentifier string, record *govpsie.Record) error
	AddReverse(ctx context.Context, reverseReq *govpsie.ReverseRequest) error
	ListReversePTRRecords(ctx context.Context) ([]govpsie.ReversePTR, error)
	UpdateReverse(ctx context.Context, reverseReq *govpsie.ReverseRequest) error
	DeleteReverse(ctx context.Context, ip, vmIdentifier string) error
}
