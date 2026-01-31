package ip

import (
	"context"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockIPAPI implements IPAPI for unit testing.
type mockIPAPI struct {
	ListPublicIPsFn  func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
	ListPrivateIPsFn func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
	ListAllIPsFn     func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error)
}

func (m *mockIPAPI) ListPublicIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
	return m.ListPublicIPsFn(ctx, options)
}

func (m *mockIPAPI) ListPrivateIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
	return m.ListPrivateIPsFn(ctx, options)
}

func (m *mockIPAPI) ListAllIPs(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
	return m.ListAllIPsFn(ctx, options)
}

// Compile-time check: mockIPAPI satisfies IPAPI.
var _ IPAPI = &mockIPAPI{}

func TestUnitIPAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockIPAPI{
		ListPublicIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
			return []govpsie.IP{}, nil
		},
		ListPrivateIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
			return []govpsie.IP{}, nil
		},
		ListAllIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
			return []govpsie.IP{}, nil
		},
	}

	var api IPAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy IPAPI interface")
	}
}

func TestUnitIPAPI_ListPublicIPs(t *testing.T) {
	mock := &mockIPAPI{
		ListPublicIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
			return []govpsie.IP{
				{IP: "1.2.3.4", Type: "public", Hostname: "server1"},
				{IP: "5.6.7.8", Type: "public", Hostname: "server2"},
			}, nil
		},
	}

	ips, err := mock.ListPublicIPs(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(ips))
	}
	if ips[0].IP != "1.2.3.4" {
		t.Fatalf("expected IP '1.2.3.4', got %q", ips[0].IP)
	}
	if ips[0].Hostname != "server1" {
		t.Fatalf("expected Hostname 'server1', got %q", ips[0].Hostname)
	}
}

func TestUnitIPAPI_ListPrivateIPs(t *testing.T) {
	mock := &mockIPAPI{
		ListPrivateIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
			return []govpsie.IP{
				{IP: "10.0.0.1", Type: "private", Hostname: "server1"},
			}, nil
		},
	}

	ips, err := mock.ListPrivateIPs(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) != 1 {
		t.Fatalf("expected 1 IP, got %d", len(ips))
	}
	if ips[0].Type != "private" {
		t.Fatalf("expected type 'private', got %q", ips[0].Type)
	}
}

func TestUnitIPAPI_ListAllIPs(t *testing.T) {
	mock := &mockIPAPI{
		ListAllIPsFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.IP, error) {
			return []govpsie.IP{
				{IP: "1.2.3.4", Type: "public"},
				{IP: "10.0.0.1", Type: "private"},
			}, nil
		},
	}

	ips, err := mock.ListAllIPs(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ips) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(ips))
	}
}
