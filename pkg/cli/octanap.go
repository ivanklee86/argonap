package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ivanklee86/octanap/pkg/client"
)

const TIMEOUT = 120

type Config struct {
	ServerAddr      string
	Insecure        bool
	AuthToken       string
	DryRun          bool
	Labels          map[string]string
	SyncWindowsFile string
}

// Octanap is the logic/orchestrator.
type octanap struct {
	*Config

	// Client
	ArgoCDClient          client.IArgoCDClient
	ArgoCDClientConnected bool

	// Allow swapping out stdout/stderr for testing.
	Out io.Writer
	Err io.Writer
}

// New returns a new instance of octanap.
func New() *octanap {
	config := Config{}

	return &octanap{
		Config: &config,
		Out:    os.Stdout,
		Err:    os.Stdin,
	}
}

func NewWithConfig(config Config) *octanap {
	return &octanap{
		Config: &config,
		Out:    os.Stdout,
		Err:    os.Stdin,
	}
}

func (o *octanap) Connect() {
	clientConfig := client.ArgoCDClientOptions{
		ServerAddr: o.Config.ServerAddr,
		Insecure:   o.Config.Insecure,
		AuthToken:  o.Config.AuthToken,
	}
	argocdClient, err := client.New(&clientConfig)
	if err != nil {
		o.Error(fmt.Sprintf("Error creating ArgoCD client: %s", err.Error()))
	}

	o.ArgoCDClient = argocdClient
	o.ArgoCDClientConnected = true
}

func (o *octanap) ClearSyncWindows() {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT)
	defer cancel()

	appProjects, err := o.ArgoCDClient.ListProjects(ctxTimeout)
	if err != nil {
		o.Error(fmt.Sprintf("Error fetching Projects: %s", err.Error()))
	}

	appProjectsToClear := filterProjects(appProjects, o.Config.Labels, true)

	for _, appProjectToClear := range appProjectsToClear {
		appProjectToClear.Spec.SyncWindows = nil
		_, err := o.ArgoCDClient.UpdateProject(ctxTimeout, appProjectToClear)
		if err != nil {
			o.Error(fmt.Sprintf("Error updating %s project: %s", appProjectToClear.ObjectMeta.Name, err))
		}
	}
}
