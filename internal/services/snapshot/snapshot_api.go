package snapshot

import (
	"context"

	"github.com/vpsie/govpsie"
)

// SnapshotAPI defines the subset of govpsie.SnapshotService methods
// used by the server_snapshot, snapshot_policy resources and their data sources.
type SnapshotAPI interface {
	Create(ctx context.Context, name, vmIdentifier, note string) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Snapshot, error)
	Get(ctx context.Context, buckupIdentifier string) (*govpsie.Snapshot, error)
	Update(ctx context.Context, snapshotIdentifier, newNote string) error
	Delete(ctx context.Context, snapshotIdentifier, reason, note string) error
	CreateSnapShotPolicy(ctx context.Context, createReq *govpsie.CreateSnapShotPolicyReq) error
	GetSnapShotPolicy(ctx context.Context, identifier string) (*govpsie.SnapShotPolicy, error)
	ListSnapShotPolicies(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.SnapShotPolicyListDetail, error)
	AttachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error
	DetachSnapShotPolicy(ctx context.Context, policyId string, vms []string) error
	DeleteSnapShotPolicy(ctx context.Context, policyId, identifier string) error
}
