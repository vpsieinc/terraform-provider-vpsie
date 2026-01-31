package project

import (
	"context"
	"fmt"
	"testing"

	"github.com/vpsie/govpsie"
)

// mockProjectAPI implements ProjectAPI for unit testing.
type mockProjectAPI struct {
	CreateFn func(ctx context.Context, projectReq *govpsie.CreateProjectRequest) error
	GetFn    func(ctx context.Context, identifier string) (*govpsie.Project, error)
	ListFn   func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Project, error)
	DeleteFn func(ctx context.Context, id string) error
}

func (m *mockProjectAPI) Create(ctx context.Context, projectReq *govpsie.CreateProjectRequest) error {
	return m.CreateFn(ctx, projectReq)
}

func (m *mockProjectAPI) Get(ctx context.Context, identifier string) (*govpsie.Project, error) {
	return m.GetFn(ctx, identifier)
}

func (m *mockProjectAPI) List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Project, error) {
	return m.ListFn(ctx, options)
}

func (m *mockProjectAPI) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

// Compile-time check: mockProjectAPI satisfies ProjectAPI.
var _ ProjectAPI = &mockProjectAPI{}

func TestUnitProjectAPI_MockSatisfiesInterface(t *testing.T) {
	mock := &mockProjectAPI{
		CreateFn: func(ctx context.Context, projectReq *govpsie.CreateProjectRequest) error {
			return nil
		},
		GetFn: func(ctx context.Context, identifier string) (*govpsie.Project, error) {
			return &govpsie.Project{}, nil
		},
		ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Project, error) {
			return []govpsie.Project{}, nil
		},
		DeleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}

	var api ProjectAPI = mock
	_ = api // compile-time interface satisfaction verified by var _ above
}

func TestUnitProjectAPI_GetProjectByName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		projects    []govpsie.Project
		getProject  *govpsie.Project
		expectErr   bool
	}{
		{
			name:        "project found by name",
			projectName: "test-project",
			projects: []govpsie.Project{
				{Name: "other-project", Identifier: "id-1"},
				{Name: "test-project", Identifier: "id-2"},
			},
			getProject: &govpsie.Project{Name: "test-project", Identifier: "id-2", ID: 42},
			expectErr:  false,
		},
		{
			name:        "project not found",
			projectName: "missing-project",
			projects: []govpsie.Project{
				{Name: "other-project", Identifier: "id-1"},
			},
			expectErr: true,
		},
		{
			name:        "empty project list",
			projectName: "test-project",
			projects:    []govpsie.Project{},
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectAPI{
				ListFn: func(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Project, error) {
					return tt.projects, nil
				},
				GetFn: func(ctx context.Context, identifier string) (*govpsie.Project, error) {
					if tt.getProject != nil {
						return tt.getProject, nil
					}
					return nil, fmt.Errorf("not found")
				},
			}

			r := &projectResource{client: mock}
			project, err := r.GetProjectByName(t.Context(), tt.projectName)

			if tt.expectErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectErr && project == nil {
				t.Fatal("expected project to be non-nil")
			}
			if !tt.expectErr && project.Name != tt.projectName {
				t.Fatalf("expected name %q, got %q", tt.projectName, project.Name)
			}
		})
	}
}

func TestUnitProjectAPI_Delete(t *testing.T) {
	var calledWith string

	mock := &mockProjectAPI{
		DeleteFn: func(ctx context.Context, id string) error {
			calledWith = id
			return nil
		},
	}

	err := mock.Delete(t.Context(), "project-id-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWith != "project-id-123" {
		t.Fatalf("expected identifier 'project-id-123', got %q", calledWith)
	}
}
