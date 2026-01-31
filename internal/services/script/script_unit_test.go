package script

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockScriptAPI implements ScriptAPI for unit testing.
type mockScriptAPI struct {
	CreateScriptFn func(ctx context.Context, createScriptRequest *govpsie.CreateScriptRequest) error
	GetScriptsFn   func(ctx context.Context) ([]govpsie.Script, error)
	GetScriptFn    func(ctx context.Context, scriptId string) (govpsie.ScriptDetail, error)
	UpdateScriptFn func(ctx context.Context, scriptUpdateRequest *govpsie.ScriptUpdateRequest) error
	DeleteScriptFn func(ctx context.Context, scriptId string) error
}

func (m *mockScriptAPI) CreateScript(ctx context.Context, createScriptRequest *govpsie.CreateScriptRequest) error {
	return m.CreateScriptFn(ctx, createScriptRequest)
}

func (m *mockScriptAPI) GetScripts(ctx context.Context) ([]govpsie.Script, error) {
	return m.GetScriptsFn(ctx)
}

func (m *mockScriptAPI) GetScript(ctx context.Context, scriptId string) (govpsie.ScriptDetail, error) {
	return m.GetScriptFn(ctx, scriptId)
}

func (m *mockScriptAPI) UpdateScript(ctx context.Context, scriptUpdateRequest *govpsie.ScriptUpdateRequest) error {
	return m.UpdateScriptFn(ctx, scriptUpdateRequest)
}

func (m *mockScriptAPI) DeleteScript(ctx context.Context, scriptId string) error {
	return m.DeleteScriptFn(ctx, scriptId)
}

// Compile-time check: mockScriptAPI satisfies ScriptAPI.
var _ ScriptAPI = &mockScriptAPI{}

func TestUnitScriptAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockScriptAPI{
		CreateScriptFn: func(ctx context.Context, req *govpsie.CreateScriptRequest) error {
			return nil
		},
		GetScriptsFn: func(ctx context.Context) ([]govpsie.Script, error) {
			return []govpsie.Script{}, nil
		},
		GetScriptFn: func(ctx context.Context, scriptId string) (govpsie.ScriptDetail, error) {
			return govpsie.ScriptDetail{}, nil
		},
		UpdateScriptFn: func(ctx context.Context, req *govpsie.ScriptUpdateRequest) error {
			return nil
		},
		DeleteScriptFn: func(ctx context.Context, scriptId string) error {
			return nil
		},
	}

	var api ScriptAPI = mock
	if api == nil {
		t.Fatal("expected mock to satisfy ScriptAPI interface")
	}
}

func TestUnitScriptAPI_GetScriptByName(t *testing.T) {
	tests := []struct {
		name       string
		scriptName string
		scripts    []govpsie.Script
		getScript  *govpsie.ScriptDetail
		expectErr  bool
	}{
		{
			name:       "script found by name",
			scriptName: "test-script",
			scripts: []govpsie.Script{
				{ScriptName: "other-script.sh", Identifier: "id-1"},
				{ScriptName: "test-script.sh", Identifier: "id-2"},
			},
			getScript: &govpsie.ScriptDetail{ScriptName: "test-script.sh", Identifier: "id-2", ID: 42},
			expectErr: false,
		},
		{
			name:       "script not found",
			scriptName: "missing-script",
			scripts: []govpsie.Script{
				{ScriptName: "other-script.sh", Identifier: "id-1"},
			},
			expectErr: true,
		},
		{
			name:       "empty script list",
			scriptName: "test-script",
			scripts:    []govpsie.Script{},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockScriptAPI{
				GetScriptsFn: func(ctx context.Context) ([]govpsie.Script, error) {
					return tt.scripts, nil
				},
				GetScriptFn: func(ctx context.Context, scriptId string) (govpsie.ScriptDetail, error) {
					if tt.getScript != nil {
						return *tt.getScript, nil
					}
					return govpsie.ScriptDetail{}, fmt.Errorf("not found")
				},
			}

			r := &scriptResource{client: mock}
			script, err := r.GetScriptByName(context.Background(), tt.scriptName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && script == nil {
				t.Fatal("expected script to be non-nil")
			}
		})
	}
}

func TestUnitScriptAPI_DeleteScript(t *testing.T) {
	var calledWith string

	mock := &mockScriptAPI{
		DeleteScriptFn: func(ctx context.Context, scriptId string) error {
			calledWith = scriptId
			return nil
		},
	}

	err := mock.DeleteScript(context.Background(), "script-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith != "script-id-123" {
		t.Fatalf("expected identifier 'script-id-123', got %q", calledWith)
	}
}
