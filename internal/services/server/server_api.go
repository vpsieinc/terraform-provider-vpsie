package server

import (
	"context"

	"github.com/vpsie/govpsie"
)

// ServerAPI defines the subset of govpsie.ServerService methods
// used by the server resource and data source in this provider.
type ServerAPI interface {
	CreateServer(ctx context.Context, req *govpsie.CreateServerRequest) error
	List(ctx context.Context, options *govpsie.ListOptions) ([]govpsie.VmData, error)
	GetServerByIdentifier(ctx context.Context, identifierId string) (*govpsie.VmData, error)
	DeleteServer(ctx context.Context, identifierId, password, reason, note string) error
	ChangeHostName(ctx context.Context, identifierId string, newHostname string) error
	StartServer(ctx context.Context, identifierId string) error
	StopServer(ctx context.Context, identifierId string) error
	Lock(ctx context.Context, identifierId string) error
	UnLock(ctx context.Context, identifierId string) error
	AddSsh(ctx context.Context, identifierId, sshKeyIdentifier string) error
	AddScript(ctx context.Context, identifierId, scriptIdentifier string) error
	ResizeServer(ctx context.Context, identifierId, cpu, ram string) error
}
