package script

import (
	"context"

	"github.com/vpsie/govpsie"
)

// ScriptAPI defines the subset of govpsie.ScriptsService methods
// used by the script resource and data source in this provider.
type ScriptAPI interface {
	CreateScript(ctx context.Context, createScriptRequest *govpsie.CreateScriptRequest) error
	GetScripts(ctx context.Context) ([]govpsie.Script, error)
	GetScript(ctx context.Context, scriptId string) (govpsie.ScriptDetail, error)
	UpdateScript(ctx context.Context, scriptUpdateRequest *govpsie.ScriptUpdateRequest) error
	DeleteScript(ctx context.Context, scriptId string) error
}
