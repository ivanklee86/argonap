package cli

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/ivanklee86/argonap/pkg/client"
)

type StatusType string

const (
	StatusSuccess    StatusType = "Success"
	StatusFailure    StatusType = "Failure"
	StatusIncomplete StatusType = "Incomplete"
)

type WorkerResult struct {
	Status      StatusType
	SyncWindows int
	ProjectName string
	Err         *error
}

func SetWorker(id int, client client.IArgoCDClient, context context.Context, syncWindowsToSet []v1alpha1.SyncWindow, projectChannel <-chan *v1alpha1.AppProject, resultChannel chan<- WorkerResult) {
	for project := range projectChannel {
		result := WorkerResult{
			Status:      StatusIncomplete,
			ProjectName: project.ObjectMeta.Name,
		}

		appProjectToUpdate, err := client.GetProject(context, project.ObjectMeta.Name)
		if err != nil {
			result.Status = StatusFailure
			result.Err = &err
			resultChannel <- result
		}

		var mergedSyncWindows v1alpha1.SyncWindows

		mergedSyncWindows = appProjectToUpdate.Spec.SyncWindows
		for _, syncWindow := range syncWindowsToSet {
			mergedSyncWindows = append(mergedSyncWindows, &syncWindow)
		}

		appProjectToUpdate.Spec.SyncWindows = mergedSyncWindows

		_, err = client.UpdateProject(context, *appProjectToUpdate)
		if err != nil {
			result.Status = StatusFailure
			result.Err = &err
			resultChannel <- result
		}

		result.Status = StatusSuccess
		result.SyncWindows = len(mergedSyncWindows)
		resultChannel <- result
	}
}
