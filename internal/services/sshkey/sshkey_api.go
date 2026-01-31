package sshkey

import (
	"context"

	"github.com/vpsie/govpsie"
)

// SshkeyAPI defines the subset of govpsie.SshkeysService methods
// used by the sshkey resource and data source in this provider.
type SshkeyAPI interface {
	Create(ctx context.Context, privateKey, name string) error
	List(ctx context.Context) ([]govpsie.SShKey, error)
	Get(ctx context.Context, sshKeyIdentifier string) (*govpsie.SShKey, error)
	Delete(ctx context.Context, sshKeyIdentifier string) error
}
