package cli

import (
	"fmt"
	"testing"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/ivanklee86/argonap/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestUpdateWorker(t *testing.T) {

	t.Run("update", func(t *testing.T) {
		testArgoCDClient := client.CreateTestClient("../../.env")
		appProjects := client.GenerateTestProjects("../../.env")
		// defer client.DeleteTestProjects(appProjects, "../../.env")

		syncWindows, err := readSyncWindowsFromFile("../../integration/exampleSyncWindows.json")
		assert.Nil(t, err)
		projectChannel := make(chan *v1alpha1.AppProject, 3)
		resultChannel := make(chan WorkerResult, 3)
		timeout := 30

		for i := 1; i <= 2; i++ {
			go SetWorker(i, testArgoCDClient, timeout, syncWindows, projectChannel, resultChannel)
		}

		for _, project := range appProjects {
			projectChannel <- project
		}
		close(projectChannel)

		for a := 1; a <= len(projectChannel); a++ {
			result := <-resultChannel
			assert.Equal(t, result.Status, StatusSuccess)
			fmt.Print(result)
		}
	})
}
