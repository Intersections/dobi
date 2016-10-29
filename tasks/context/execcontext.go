package context

import (
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/execenv"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/common"
	docker "github.com/fsouza/go-dockerclient"
)

// ExecuteContext contains all the context for task execution
type ExecuteContext struct {
	modified    map[common.TaskName]bool
	Resources   *ResourceCollection
	Client      client.DockerClient
	authConfigs *docker.AuthConfigurations
	WorkingDir  string
	Env         *execenv.ExecEnv
	Quiet       bool
}

// IsModified returns true if any of the tasks named in names has been modified
// during this execution
func (ctx *ExecuteContext) IsModified(names ...common.TaskName) bool {
	for _, name := range names {
		if modified, _ := ctx.modified[name]; modified {
			return true
		}
	}
	return false
}

// SetModified sets the task name as modified
func (ctx *ExecuteContext) SetModified(name common.TaskName) {
	ctx.modified[name] = true
}

// GetAuthConfig returns the auth configuration for the repo
func (ctx *ExecuteContext) GetAuthConfig(repo string) docker.AuthConfiguration {
	auth, ok := ctx.authConfigs.Configs[repo]
	if !ok {
		logging.Log.Warnf("Missing auth config for %q", repo)
	}
	return auth
}

// NewExecuteContext craetes a new empty ExecuteContext
func NewExecuteContext(
	config *config.Config,
	client client.DockerClient,
	execEnv *execenv.ExecEnv,
	quiet bool,
) *ExecuteContext {

	authConfigs, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		logging.Log.Warnf("Failed to load auth config: %s", err)
	}

	return &ExecuteContext{
		modified:    make(map[common.TaskName]bool),
		Resources:   newResourceCollection(),
		WorkingDir:  config.WorkingDir,
		Client:      client,
		authConfigs: authConfigs,
		Env:         execEnv,
		Quiet:       quiet,
	}
}
