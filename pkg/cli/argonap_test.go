package cli

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/ivanklee86/argonap/pkg/client"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCharacterCount(t *testing.T) {
	testString := "key=value"

	t.Run("Can count ='s correctly", func(t *testing.T) {
		assert.Equal(t, countCharacterOccurrences(testString, '='), 1)
	})
}

func TestArgonapHappyPath(t *testing.T) {
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
		Timeout:         240,
		Workers:         4,
	}

	b := bytes.NewBufferString("")

	argonap := NewWithConfig(config)
	argonap.Out = b
	argonap.Err = b

	t.Run("argonap configuration setup", func(t *testing.T) {
		expectedMap := make(map[string]string)
		expectedMap["purpose"] = "tests"
		assert.Equal(t, argonap.Config.Labels, expectedMap)
	})

	t.Run("argonap can clear all SyncWindows", func(t *testing.T) {
		testArgoCDClient := client.CreateTestClient("../../.env")
		appProjects := client.GenerateTestProjects("../../.env")
		defer client.DeleteTestProjects(appProjects, "../../.env")

		argonap.Connect()
		argonap.ClearSyncWindows()

		assert.Nil(t, err)
		for _, appProject := range appProjects {
			updatedAppProject, _ := testArgoCDClient.GetProject(context.Background(), appProject.Name)
			assert.Nil(t, updatedAppProject.Spec.SyncWindows)
		}
	})

	t.Run("argonap can set SyncWindows", func(t *testing.T) {
		testArgoCDClient := client.CreateTestClient("../../.env")
		appProjects := client.GenerateTestProjects("../../.env")
		defer client.DeleteTestProjects(appProjects, "../../.env")

		argonap.Connect()
		argonap.SetSyncWindows()

		assert.Nil(t, err)
		for index, appProject := range appProjects {
			updatedAppProject, _ := testArgoCDClient.GetProject(context.Background(), appProject.Name)
			if index == 1 { // SyncWindow already exists
				assert.Len(t, updatedAppProject.Spec.SyncWindows, 3)
			} else {
				assert.Len(t, updatedAppProject.Spec.SyncWindows, 2)
			}
		}
	})
}
