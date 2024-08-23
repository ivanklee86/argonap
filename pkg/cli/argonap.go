package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/ivanklee86/argonap/pkg/client"
)

const TIMEOUT = 120

type Config struct {
	ServerAddr      string
	Insecure        bool
	AuthToken       string
	DryRun          bool
	ProjectName     string
	LabelsAsStrings []string
	Labels          map[string]string
	SyncWindowsFile string
}

// Argonap is the logic/orchestrator.
type Argonap struct {
	*Config

	// Client
	ArgoCDClient          client.IArgoCDClient
	ArgoCDClientConnected bool

	// Allow swapping out stdout/stderr for testing.
	Out io.Writer
	Err io.Writer
}

func countCharacterOccurrences(s string, c rune) int {
	count := 0
	for _, char := range s {
		if char == c {
			count++
		}
	}
	return count
}

func labelStringsToMap(labelsAsStrings []string) map[string]string {
	labels := make(map[string]string)

	for _, labelString := range labelsAsStrings {
		if countCharacterOccurrences(labelString, '=') == 1 {
			kv := strings.Split(labelString, "=")

			if len(kv) == 2 {
				labels[kv[0]] = kv[1]
			}
		}
	}

	return labels
}

// New returns a new instance of argonap.
func New() *Argonap {
	config := Config{}
	config.Labels = labelStringsToMap(config.LabelsAsStrings)

	return &Argonap{
		Config: &config,
		Out:    os.Stdout,
		Err:    os.Stdin,
	}
}

func NewWithConfig(config Config) *Argonap {
	config.Labels = labelStringsToMap(config.LabelsAsStrings)

	return &Argonap{
		Config: &config,
		Out:    os.Stdout,
		Err:    os.Stdin,
	}
}

func (a *Argonap) Connect() {
	if a.Config.ServerAddr == "" {
		a.Error("ArgoCD server address not set.")
	}

	if a.Config.AuthToken == "" {
		a.Error("ArgoCD JWT auth token is not set.")
	}

	clientConfig := client.ArgoCDClientOptions{
		ServerAddr: a.Config.ServerAddr,
		Insecure:   a.Config.Insecure,
		AuthToken:  a.Config.AuthToken,
	}
	argocdClient, err := client.New(&clientConfig)
	if err != nil {
		a.Error(fmt.Sprintf("Creating ArgoCD client: %s", err.Error()))
	}

	a.ArgoCDClient = argocdClient
	a.ArgoCDClientConnected = true
}

func (a *Argonap) ClearSyncWindows() {
	a.OutputHeading("Clearing SyncWindows on matching AppProjects.")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT)
	defer cancel()

	appProjects, err := a.ArgoCDClient.ListProjects(ctxTimeout)
	if err != nil {
		a.Error(fmt.Sprintf("Error fetching Projects: %s", err.Error()))
	}

	appProjectsToClear := filterProjects(appProjects, a.Config.ProjectName, a.Config.Labels, true)

	var selectedProjectNames []string
	for _, p := range appProjectsToClear {
		selectedProjectNames = append(selectedProjectNames, p.ObjectMeta.Name)
	}

	a.Output(fmt.Sprintf("%d projects found with SyncWindows: %s", len(appProjectsToClear), strings.Join(selectedProjectNames, ", ")))

	for _, appProjectToClear := range appProjectsToClear {
		appProjectToClear.Spec.SyncWindows = nil
		_, err := a.ArgoCDClient.UpdateProject(ctxTimeout, appProjectToClear)

		a.Output(fmt.Sprintf("Cleared SyncWindows from project %s.", appProjectToClear.ObjectMeta.Name))
		if err != nil {
			a.Error(fmt.Sprintf("Error updating %s project: %s", appProjectToClear.ObjectMeta.Name, err))
		}
	}
}

func (a *Argonap) SetSyncWindows() {
	a.OutputHeading("Setting SyncWindows on matching AppProjects.")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT)
	defer cancel()

	syncWindowsToSet, err := readSyncWindowsFromFile(a.Config.SyncWindowsFile)
	if err != nil {
		a.Error(fmt.Sprintf("Unable to read SyncWindows file. %s", err.Error()))
	}

	appProjects, err := a.ArgoCDClient.ListProjects(ctxTimeout)
	if err != nil {
		a.Error(fmt.Sprintf("Error fetching Projects. %s", err.Error()))
	}

	appProjectsToUpdate := filterProjects(appProjects, a.Config.ProjectName, a.Config.Labels, false)

	var selectedProjectNames []string
	for _, p := range appProjectsToUpdate {
		selectedProjectNames = append(selectedProjectNames, p.ObjectMeta.Name)
	}
	a.Output(fmt.Sprintf("%d projects found to update: %s", len(appProjectsToUpdate), strings.Join(selectedProjectNames, ", ")))

	for _, appProjectToUpdate := range appProjectsToUpdate {
		var mergedSyncWindows v1alpha1.SyncWindows

		mergedSyncWindows = appProjectToUpdate.Spec.SyncWindows
		for _, syncWindow := range syncWindowsToSet {
			mergedSyncWindows = append(mergedSyncWindows, &syncWindow)
		}

		appProjectToUpdate.Spec.SyncWindows = mergedSyncWindows

		_, err := a.ArgoCDClient.UpdateProject(ctxTimeout, appProjectToUpdate)

		a.Output(fmt.Sprintf("Added SyncWindows to project %s.", appProjectToUpdate.ObjectMeta.Name))

		if err != nil {
			a.Error(fmt.Sprintf("Error updating %s project: %s", appProjectToUpdate.ObjectMeta.Name, err))
		}
	}
}
