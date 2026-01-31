package kubernetes

import (
	"context"

	"github.com/vpsie/govpsie"
)

// KubernetesAPI defines the subset of govpsie.K8sService methods
// used by the kubernetes resource, group resource, and data sources in this provider.
type KubernetesAPI interface {
	Create(ctx context.Context, createReq *govpsie.CreateK8sReq) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.ListK8s, error)
	Get(ctx context.Context, identifier string) (*govpsie.K8s, error)
	Delete(ctx context.Context, identifier, reason, note string) error
	AddSlave(ctx context.Context, identifier string) error
	RemoveSlave(ctx context.Context, identifier string) error
	ListK8sGroups(ctx context.Context, identifier string) ([]govpsie.K8sGroup, error)
	CreateK8sGroup(ctx context.Context, createReq *govpsie.CreateK8sGroupReq) error
	DeleteK8sGroup(ctx context.Context, groupId string, reason, note string) error
	AddNode(ctx context.Context, identifier, nodeType string, groupId int) error
	RemoveNode(ctx context.Context, identifier, nodeType string, groupId int) error
}
