package project

import (
	"context"

	"github.com/vpsie/govpsie"
)

// ProjectAPI defines the subset of govpsie.ProjectsService methods
// used by the project resource and data source in this provider.
type ProjectAPI interface {
	Create(ctx context.Context, projectReq *govpsie.CreateProjectRequest) error
	Get(ctx context.Context, identifier string) (*govpsie.Project, error)
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Project, error)
	Delete(ctx context.Context, id string) error
}
