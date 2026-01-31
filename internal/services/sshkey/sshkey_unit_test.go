package sshkey

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockSshkeyAPI implements SshkeyAPI for unit testing.
type mockSshkeyAPI struct {
	CreateFn func(ctx context.Context, privateKey, name string) error
	ListFn   func(ctx context.Context) ([]govpsie.SShKey, error)
	GetFn    func(ctx context.Context, sshKeyIdentifier string) (*govpsie.SShKey, error)
	DeleteFn func(ctx context.Context, sshKeyIdentifier string) error
}

func (m *mockSshkeyAPI) Create(ctx context.Context, privateKey, name string) error {
	return m.CreateFn(ctx, privateKey, name)
}

func (m *mockSshkeyAPI) List(ctx context.Context) ([]govpsie.SShKey, error) {
	return m.ListFn(ctx)
}

func (m *mockSshkeyAPI) Get(ctx context.Context, sshKeyIdentifier string) (*govpsie.SShKey, error) {
	return m.GetFn(ctx, sshKeyIdentifier)
}

func (m *mockSshkeyAPI) Delete(ctx context.Context, sshKeyIdentifier string) error {
	return m.DeleteFn(ctx, sshKeyIdentifier)
}

// Compile-time check: mockSshkeyAPI satisfies SshkeyAPI.
var _ SshkeyAPI = &mockSshkeyAPI{}

func TestUnitSshkeyAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockSshkeyAPI{
		CreateFn: func(ctx context.Context, privateKey, name string) error {
			return nil
		},
		ListFn: func(ctx context.Context) ([]govpsie.SShKey, error) {
			return []govpsie.SShKey{}, nil
		},
		GetFn: func(ctx context.Context, sshKeyIdentifier string) (*govpsie.SShKey, error) {
			return &govpsie.SShKey{}, nil
		},
		DeleteFn: func(ctx context.Context, sshKeyIdentifier string) error {
			return nil
		},
	}

	var api SshkeyAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy SshkeyAPI interface")
	}
}

func TestUnitSshkeyAPI_GetSshkeyByName(t *testing.T) {
	tests := []struct {
		name        string
		sshkeyName  string
		sshkeys     []govpsie.SShKey
		expectFound bool
		expectErr   bool
	}{
		{
			name:       "sshkey found by name",
			sshkeyName: "my-key",
			sshkeys: []govpsie.SShKey{
				{Name: "other-key", Id: 1},
				{Name: "my-key", Id: 2, CreatedOn: "2024-01-01"},
			},
			expectFound: true,
			expectErr:   false,
		},
		{
			name:       "sshkey not found",
			sshkeyName: "missing-key",
			sshkeys: []govpsie.SShKey{
				{Name: "other-key", Id: 1},
			},
			expectFound: false,
			expectErr:   true,
		},
		{
			name:        "empty sshkey list",
			sshkeyName:  "my-key",
			sshkeys:     []govpsie.SShKey{},
			expectFound: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSshkeyAPI{
				ListFn: func(ctx context.Context) ([]govpsie.SShKey, error) {
					return tt.sshkeys, nil
				},
			}

			r := &sshkeyResource{client: mock}
			key, err := r.GetSshkeyByName(context.Background(), tt.sshkeyName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectFound && key == nil {
				t.Fatal("expected sshkey to be non-nil when found")
			}
			if tt.expectFound && key.Name != tt.sshkeyName {
				t.Fatalf("expected name %q, got %q", tt.sshkeyName, key.Name)
			}
		})
	}
}

func TestUnitSshkeyAPI_ListError(t *testing.T) {
	mock := &mockSshkeyAPI{
		ListFn: func(ctx context.Context) ([]govpsie.SShKey, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &sshkeyResource{client: mock}
	_, err := r.GetSshkeyByName(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error from List failure, got nil")
	}
}

func TestUnitSshkeyAPI_CreateAndGet(t *testing.T) {
	var createCalled bool
	var createKey, createName string

	mock := &mockSshkeyAPI{
		CreateFn: func(ctx context.Context, privateKey, name string) error {
			createCalled = true
			createKey = privateKey
			createName = name
			return nil
		},
		GetFn: func(ctx context.Context, sshKeyIdentifier string) (*govpsie.SShKey, error) {
			if sshKeyIdentifier != "my-key" {
				return nil, fmt.Errorf("sshkey %s not found", sshKeyIdentifier)
			}
			return &govpsie.SShKey{
				Name:       "my-key",
				Id:         5,
				PrivateKey: "ssh-ed25519 AAAA...",
				CreatedOn:  "2024-01-01",
				CreatedBy:  "admin",
				UserId:     1,
			}, nil
		},
	}

	err := mock.Create(context.Background(), "ssh-ed25519 AAAA...", "my-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !createCalled {
		t.Fatal("expected Create to be called")
	}
	if createKey != "ssh-ed25519 AAAA..." {
		t.Fatalf("expected privateKey 'ssh-ed25519 AAAA...', got %q", createKey)
	}
	if createName != "my-key" {
		t.Fatalf("expected name 'my-key', got %q", createName)
	}

	key, err := mock.Get(context.Background(), "my-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.Name != "my-key" {
		t.Fatalf("expected name 'my-key', got %q", key.Name)
	}
	if key.Id != 5 {
		t.Fatalf("expected Id 5, got %d", key.Id)
	}
}
