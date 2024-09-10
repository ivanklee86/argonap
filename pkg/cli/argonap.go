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
	"github.com/jedib0t/go-pretty/v6/list"
)

const TIMEOUT = 120

type Config struct {
	ServerAddr      string
	Insecure        bool
	AuthToken       string
	DryRun          bool
	ProjectNames    []string
	LabelsAsStrings []string
	Labels          map[string]string
	SyncWindowsFile string
	Timeout         int
	Workers         int
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

func displayFilteredProjects(projects *[]v1alpha1.AppProject) string {
	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)
	for _, p := range *projects {
		l.AppendItem(p.ObjectMeta.Name)
	}

	return l.Render()
}

// New returns a new instance of argonap.
func New() *Argonap {
	config := Config{}
	config.Labels = labelStringsToMap(config.LabelsAsStrings)

	if config.Timeout == 0 {
		config.Timeout = 240 // Set default for tests, etc.
	}

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
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(a.Config.Timeout)*time.Second)
	defer cancel()

	a.OutputHeading("🔍  Searching for Projects.")
	appProjects, err := a.ArgoCDClient.ListProjects(ctxTimeout)
	if err != nil {
		a.Error(fmt.Sprintf("Error fetching Projects: %s", err.Error()))
	}

	appProjectsToClear := filterProjects(appProjects, a.Config.ProjectNames, a.Config.Labels, true)

	a.Output(fmt.Sprintf("%d projects found with SyncWindows.", len(appProjectsToClear)))
	if len(appProjectsToClear) > 0 {
		a.Output(displayFilteredProjects(&appProjectsToClear))
	}

	projects := make(chan *v1alpha1.AppProject, len(appProjectsToClear))
	results := make(chan WorkerResult, len(appProjectsToClear))

	a.OutputHeading("🛠️  Setting SyncWindows on Projects.")
	// Start workers
	for i := 1; i <= a.Config.Workers; i++ {
		go ClearWorker(i, a.ArgoCDClient, ctxTimeout, projects, results)
	}

	// Load projects
	for _, project := range appProjectsToClear {
		projects <- &project
	}
	close(projects)

	for r := 0; r < len(appProjectsToClear); r++ {
		result := <-results
		a.OutputResult(result)
	}

	a.OutputHeading("🎉  Complete!")
}

func (a *Argonap) SetSyncWindows() {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(a.Config.Timeout)*time.Second)
	defer cancel()

	syncWindowsToSet, err := readSyncWindowsFromFile(a.Config.SyncWindowsFile)
	if err != nil {
		a.Error(fmt.Sprintf("Unable to read SyncWindows file. %s", err.Error()))
	}

	a.OutputHeading("🔍  Searching for Projects.")
	appProjects, err := a.ArgoCDClient.ListProjects(ctxTimeout)
	if err != nil {
		a.Error(fmt.Sprintf("Error fetching Projects. %s", err.Error()))
	}

	appProjectsToUpdate := filterProjects(appProjects, a.Config.ProjectNames, a.Config.Labels, false)

	a.Output(fmt.Sprintf("%d projects found to update:", len(appProjectsToUpdate)))
	if len(appProjectsToUpdate) > 0 {
		a.Output(displayFilteredProjects(&appProjectsToUpdate))
	}

	projects := make(chan *v1alpha1.AppProject, len(appProjectsToUpdate))
	results := make(chan WorkerResult, len(appProjectsToUpdate))

	a.OutputHeading("🛠️  Setting SyncWindows on Projects.")
	// Start workers
	for i := 1; i <= a.Config.Workers; i++ {
		go SetWorker(i, a.ArgoCDClient, ctxTimeout, syncWindowsToSet, projects, results)
	}

	// Load projects
	for _, project := range appProjectsToUpdate {
		projects <- &project
	}
	close(projects)

	for r := 0; r < len(appProjectsToUpdate); r++ {
		result := <-results
		a.OutputResult(result)
	}
	a.OutputHeading("🎉  Complete!")
}
