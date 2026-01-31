package server

import (
	"context"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockServerAPI implements ServerAPI for unit testing.
type mockServerAPI struct {
	CreateServerFn          func(ctx context.Context, req *govpsie.CreateServerRequest) error
	ListFn                  func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VmData, error)
	GetServerByIdentifierFn func(ctx context.Context, identifierId string) (*govpsie.VmData, error)
	DeleteServerFn          func(ctx context.Context, identifierId, password, reason, note string) error
	ChangeHostNameFn        func(ctx context.Context, identifierId string, newHostname string) error
	StartServerFn           func(ctx context.Context, identifierId string) error
	StopServerFn            func(ctx context.Context, identifierId string) error
	LockFn                  func(ctx context.Context, identifierId string) error
	UnLockFn                func(ctx context.Context, identifierId string) error
	AddSshFn                func(ctx context.Context, identifierId, sshKeyIdentifier string) error
	AddScriptFn             func(ctx context.Context, identifierId, scriptIdentifier string) error
	ResizeServerFn          func(ctx context.Context, identifierId, cpu, ram string) error
}

func (m *mockServerAPI) CreateServer(ctx context.Context, req *govpsie.CreateServerRequest) error {
	return m.CreateServerFn(ctx, req)
}

func (m *mockServerAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VmData, error) {
	return m.ListFn(ctx, options)
}

func (m *mockServerAPI) GetServerByIdentifier(ctx context.Context, identifierId string) (*govpsie.VmData, error) {
	return m.GetServerByIdentifierFn(ctx, identifierId)
}

func (m *mockServerAPI) DeleteServer(ctx context.Context, identifierId, password, reason, note string) error {
	return m.DeleteServerFn(ctx, identifierId, password, reason, note)
}

func (m *mockServerAPI) ChangeHostName(ctx context.Context, identifierId string, newHostname string) error {
	return m.ChangeHostNameFn(ctx, identifierId, newHostname)
}

func (m *mockServerAPI) StartServer(ctx context.Context, identifierId string) error {
	return m.StartServerFn(ctx, identifierId)
}

func (m *mockServerAPI) StopServer(ctx context.Context, identifierId string) error {
	return m.StopServerFn(ctx, identifierId)
}

func (m *mockServerAPI) Lock(ctx context.Context, identifierId string) error {
	return m.LockFn(ctx, identifierId)
}

func (m *mockServerAPI) UnLock(ctx context.Context, identifierId string) error {
	return m.UnLockFn(ctx, identifierId)
}

func (m *mockServerAPI) AddSsh(ctx context.Context, identifierId, sshKeyIdentifier string) error {
	return m.AddSshFn(ctx, identifierId, sshKeyIdentifier)
}

func (m *mockServerAPI) AddScript(ctx context.Context, identifierId, scriptIdentifier string) error {
	return m.AddScriptFn(ctx, identifierId, scriptIdentifier)
}

func (m *mockServerAPI) ResizeServer(ctx context.Context, identifierId, cpu, ram string) error {
	return m.ResizeServerFn(ctx, identifierId, cpu, ram)
}

// Compile-time check: mockServerAPI satisfies ServerAPI.
var _ ServerAPI = &mockServerAPI{}

func TestUnitServerAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockServerAPI{
		CreateServerFn: func(ctx context.Context, req *govpsie.CreateServerRequest) error {
			return nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VmData, error) {
			return []govpsie.VmData{}, nil
		},
		GetServerByIdentifierFn: func(ctx context.Context, identifierId string) (*govpsie.VmData, error) {
			return &govpsie.VmData{}, nil
		},
		DeleteServerFn: func(ctx context.Context, identifierId, password, reason, note string) error {
			return nil
		},
		ChangeHostNameFn: func(ctx context.Context, identifierId string, newHostname string) error {
			return nil
		},
		StartServerFn: func(ctx context.Context, identifierId string) error {
			return nil
		},
		StopServerFn: func(ctx context.Context, identifierId string) error {
			return nil
		},
		LockFn: func(ctx context.Context, identifierId string) error {
			return nil
		},
		UnLockFn: func(ctx context.Context, identifierId string) error {
			return nil
		},
		AddSshFn: func(ctx context.Context, identifierId, sshKeyIdentifier string) error {
			return nil
		},
		AddScriptFn: func(ctx context.Context, identifierId, scriptIdentifier string) error {
			return nil
		},
		ResizeServerFn: func(ctx context.Context, identifierId, cpu, ram string) error {
			return nil
		},
	}

	var api ServerAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitServerAPI_CheckResourceStatus(t *testing.T) {
	tests := []struct {
		name        string
		hostname    string
		servers     []govpsie.VmData
		expectFound bool
		expectErr   bool
	}{
		{
			name:     "server found by hostname",
			hostname: "test-host",
			servers: []govpsie.VmData{
				{Hostname: "other-host", Identifier: "id-1"},
				{Hostname: "test-host", Identifier: "id-2"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:     "server not found",
			hostname: "missing-host",
			servers: []govpsie.VmData{
				{Hostname: "other-host", Identifier: "id-1"},
			},
			expectFound: false,
			expectErr:   false,
		},
		{
			name:        "empty server list",
			hostname:    "test-host",
			servers:     []govpsie.VmData{},
			expectFound: false,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockServerAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VmData, error) {
					return tt.servers, nil
				},
			}

			r := &serverResource{client: mock}
			server, found, err := r.checkResourceStatus(t.Context(), tt.hostname)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tt.expectFound {
				t.Fatalf("expected found=%v, got %v", tt.expectFound, found)
			}
			if tt.expectFound && server == nil {
				t.Fatal("expected server to be non-nil when found")
			}
			if tt.expectFound && server.Hostname != tt.hostname {
				t.Fatalf("expected hostname %q, got %q", tt.hostname, server.Hostname)
			}
		})
	}
}

func TestUnitServerAPI_GetServerByIdentifier(t *testing.T) {
	expectedServer := &govpsie.VmData{
		Hostname:   "test-host",
		Identifier: "test-id-123",
		ID:         42,
		State:      "running",
	}

	mock := &mockServerAPI{
		GetServerByIdentifierFn: func(ctx context.Context, identifierId string) (*govpsie.VmData, error) {
			if identifierId != "test-id-123" {
				t.Fatalf("expected identifier 'test-id-123', got %q", identifierId)
			}
			return expectedServer, nil
		},
	}

	server, err := mock.GetServerByIdentifier(t.Context(), "test-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if server.Hostname != "test-host" {
		t.Fatalf("expected hostname 'test-host', got %q", server.Hostname)
	}
	if server.ID != 42 {
		t.Fatalf("expected ID 42, got %d", server.ID)
	}
}

func TestUnitServerAPI_DeleteServer(t *testing.T) {
	var calledWith struct {
		identifier string
		password   string
		reason     string
		note       string
	}

	mock := &mockServerAPI{
		DeleteServerFn: func(ctx context.Context, identifierId, password, reason, note string) error {
			calledWith.identifier = identifierId
			calledWith.password = password
			calledWith.reason = reason
			calledWith.note = note
			return nil
		},
	}

	err := mock.DeleteServer(t.Context(), "server-id", "pass123", "testing", "test note")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith.identifier != "server-id" {
		t.Fatalf("expected identifier 'server-id', got %q", calledWith.identifier)
	}
	if calledWith.password != "pass123" {
		t.Fatalf("expected password 'pass123', got %q", calledWith.password)
	}
	if calledWith.reason != "testing" {
		t.Fatalf("expected reason 'testing', got %q", calledWith.reason)
	}
	if calledWith.note != "test note" {
		t.Fatalf("expected note 'test note', got %q", calledWith.note)
	}
}
