package cli

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/ivanklee86/octanap/pkg/client"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCharacterCount(t *testing.T) {
	testString := "key=value"

	t.Run("Can count ='s correctly", func(t *testing.T) {
		assert.Equal(t, countCharacterOccurrences(testString, '='), 1)
	})
}

func TestOctanapHappyPath(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		ServerAddr:      "localhost:8080",
		Insecure:        true,
		AuthToken:       os.Getenv("ARGOCD_TOKEN"),
		SyncWindowsFile: "../../integration/exampleSyncWindows.json",
		LabelsAsStrings: []string{"purpose=tests"},
	}

	b := bytes.NewBufferString("")

	octanap := NewWithConfig(config)
	octanap.Out = b
	octanap.Err = b

	t.Run("Octanap configuration setup", func(t *testing.T) {
		expectedMap := make(map[string]string)
		expectedMap["purpose"] = "tests"
		assert.Equal(t, octanap.Config.Labels, expectedMap)
	})

	t.Run("Octonap can clear all SyncWindows", func(t *testing.T) {
		testArgoCDClient := client.CreateTestClient()
		appProjects := client.GenerateTestProjects()
		defer client.DeleteTestProjects(appProjects)

		octanap.Connect()
		octanap.ClearSyncWindows()

		assert.Nil(t, err)
		for _, appProject := range appProjects {
			updatedAppProject, _ := testArgoCDClient.GetProject(context.TODO(), appProject.Name)
			assert.Nil(t, updatedAppProject.Spec.SyncWindows)
		}
	})

	t.Run("Octonap can clear set SyncWindows", func(t *testing.T) {
		testArgoCDClient := client.CreateTestClient()
		appProjects := client.GenerateTestProjects()
		defer client.DeleteTestProjects(appProjects)

		octanap.Connect()
		octanap.SetSyncWindows()

		assert.Nil(t, err)
		for index, appProject := range appProjects {
			updatedAppProject, _ := testArgoCDClient.GetProject(context.TODO(), appProject.Name)
			if index == 1 { // SyncWindow already exists
				assert.Len(t, updatedAppProject.Spec.SyncWindows, 3)
			} else {
				assert.Len(t, updatedAppProject.Spec.SyncWindows, 2)
			}
		}
	})
}
