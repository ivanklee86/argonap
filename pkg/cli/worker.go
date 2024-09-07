package cli

import (
	"context"
	"fmt"
	"time"

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
	ProjectName string
	Error       *error
}

func SetWorker(id int, client client.IArgoCDClient, timeout int, syncWindowsToSet []v1alpha1.SyncWindow, projectChannel <-chan *v1alpha1.AppProject, resultChannel chan<- WorkerResult) {
	fmt.Printf("Worker %d: Starting\n", id)
	for project := range projectChannel {
		fmt.Printf("Worker %d: Processing project %s\n", id, project.ObjectMeta.Name)
		result := WorkerResult{
			Status:      StatusIncomplete,
			ProjectName: project.ObjectMeta.Name,
		}
	
		ctxTimeout, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	
		appProjectToUpdate, err := client.GetProject(ctxTimeout, project.ObjectMeta.Name)
		if err != nil {
			result.Status = StatusFailure
			result.Error = &err
			resultChannel <- result
			fmt.Printf("Worker %d: Error %s\n", id, *result.Error)
		}
	
		var mergedSyncWindows v1alpha1.SyncWindows
	
		mergedSyncWindows = appProjectToUpdate.Spec.SyncWindows
		for _, syncWindow := range syncWindowsToSet {
			mergedSyncWindows = append(mergedSyncWindows, &syncWindow)
		}
	
		appProjectToUpdate.Spec.SyncWindows = mergedSyncWindows
	
		_, err = client.UpdateProject(ctxTimeout, *appProjectToUpdate)
		if err != nil {
			result.Status = StatusFailure
			result.Error = &err
			resultChannel <- result
			fmt.Printf("Worker %d: Error %s\n", id, *result.Error)
		}

		fmt.Printf("Worker %d: hmm...\n", id)
		result.Status = StatusSuccess
		resultChannel <- result
		fmt.Printf("Worker %d: Completed for %s\n", id, project.ObjectMeta.Name)
	}
}
