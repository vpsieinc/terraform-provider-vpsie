package accesstoken

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockAccessTokenAPI implements AccessTokenAPI for unit testing.
type mockAccessTokenAPI struct {
	CreateFn func(ctx context.Context, name, accessToken, expirationDate string) error
	ListFn   func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.AccessToken, error)
	UpdateFn func(ctx context.Context, accessTokenIdentifier, name, expirationDate string) error
	DeleteFn func(ctx context.Context, accessTokenIdentifier string) error
}

func (m *mockAccessTokenAPI) Create(ctx context.Context, name, accessToken, expirationDate string) error {
	return m.CreateFn(ctx, name, accessToken, expirationDate)
}

func (m *mockAccessTokenAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.AccessToken, error) {
	return m.ListFn(ctx, options)
}

func (m *mockAccessTokenAPI) Update(ctx context.Context, accessTokenIdentifier, name, expirationDate string) error {
	return m.UpdateFn(ctx, accessTokenIdentifier, name, expirationDate)
}

func (m *mockAccessTokenAPI) Delete(ctx context.Context, accessTokenIdentifier string) error {
	return m.DeleteFn(ctx, accessTokenIdentifier)
}

// Compile-time check: mockAccessTokenAPI satisfies AccessTokenAPI.
var _ AccessTokenAPI = &mockAccessTokenAPI{}

func TestUnitAccessTokenAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockAccessTokenAPI{
		CreateFn: func(ctx context.Context, name, accessToken, expirationDate string) error {
			return nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.AccessToken, error) {
			return []govpsie.AccessToken{}, nil
		},
		UpdateFn: func(ctx context.Context, accessTokenIdentifier, name, expirationDate string) error {
			return nil
		},
		DeleteFn: func(ctx context.Context, accessTokenIdentifier string) error {
			return nil
		},
	}

	var api AccessTokenAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitAccessTokenAPI_GetTokenByName(t *testing.T) {
	tests := []struct {
		name      string
		tokenName string
		tokens    []govpsie.AccessToken
		expectErr bool
	}{
		{
			name:      "token found by name",
			tokenName: "test-token",
			tokens: []govpsie.AccessToken{
				{Name: "other-token", AccessTokenIdentifier: "id-1"},
				{Name: "test-token", AccessTokenIdentifier: "id-2"},
			},
			expectErr: false,
		},
		{
			name:      "token not found",
			tokenName: "missing-token",
			tokens: []govpsie.AccessToken{
				{Name: "other-token", AccessTokenIdentifier: "id-1"},
			},
			expectErr: true,
		},
		{
			name:      "empty token list",
			tokenName: "test-token",
			tokens:    []govpsie.AccessToken{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAccessTokenAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.AccessToken, error) {
					return tt.tokens, nil
				},
			}

			r := &accessTokenResource{client: mock}
			token, err := r.GetTokenByName(t.Context(), tt.tokenName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && token == nil {
				t.Fatal("expected token to be non-nil")
			}
			if !tt.expectErr && token.Name != tt.tokenName {
				t.Fatalf("expected name %q, got %q", tt.tokenName, token.Name)
			}
		})
	}
}

func TestUnitAccessTokenAPI_Delete(t *testing.T) {
	var calledWith string

	mock := &mockAccessTokenAPI{
		DeleteFn: func(ctx context.Context, accessTokenIdentifier string) error {
			calledWith = accessTokenIdentifier
			return nil
		},
	}

	err := mock.Delete(t.Context(), "token-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith != "token-id-123" {
		t.Fatalf("expected identifier 'token-id-123', got %q", calledWith)
	}
}

func TestUnitAccessTokenAPI_ListError(t *testing.T) {
	mock := &mockAccessTokenAPI{
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.AccessToken, error) {
			return nil, fmt.Errorf("api error")
		},
	}

	r := &accessTokenResource{client: mock}
	_, err := r.GetTokenByName(t.Context(), "any-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "api error" {
		t.Fatalf("expected 'api error', got %q", err.Error())
	}
}
