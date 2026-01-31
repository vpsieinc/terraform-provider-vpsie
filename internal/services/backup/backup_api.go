package backup

import (
	"context"

	"github.com/vpsie/govpsie"
)

// BackupAPI defines the subset of govpsie.BackupsService methods
// used by the backup, backup_policy resources and their data sources.
type BackupAPI interface {
	CreateBackups(ctx context.Context, vmIdentifier, name, notes string) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.Backup, error)
	Get(ctx context.Context, identifier string) (*govpsie.Backup, error)
	Rename(ctx context.Context, backupIdentifier, newName string) error
	DeleteBackup(ctx context.Context, backupIdentifier, deleteReason, deleteNote string) error
	CreateBackupPolicy(ctx context.Context, createReq *govpsie.CreateBackupPolicyReq) error
	GetBackupPolicy(ctx context.Context, identifier string) (*govpsie.BackupPolicy, error)
	ListBackupPolicies(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.BackupPolicyListDetail, error)
	AttachBackupPolicy(ctx context.Context, policyId string, vms []string) error
	DetachBackupPolicy(ctx context.Context, policyId string, vms []string) error
	DeleteBackupPolicy(ctx context.Context, policyId, identifier string) error
}
