package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/ivanklee86/argonap/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	argoCDClient := client.CreateTestClient("./.env")

	t.Run("Root command", func(t *testing.T) {
		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"--server-address", "localhost:8080",
			"--plaintext",
			"--auth-token", os.Getenv("ARGOCD_TOKEN"),
		})
		err := command.Execute()
		if err != nil {
			t.Fatal(err)
		}

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, string(out), "argonap")
	})

	t.Run("Run clear command", func(t *testing.T) {
		appProjects := client.GenerateTestProjects("./.env")
		defer client.DeleteTestProjects(appProjects, "./.env")

		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"clear",
			"--server-address", "localhost:8080",
			"--name", appProjects[0].Name,
			"--name", appProjects[1].Name,
			"--name", appProjects[2].Name,
			"--plaintext",
			"--auth-token", os.Getenv("ARGOCD_TOKEN"),
		})
		err := command.Execute()
		assert.Nil(t, err)

		assert.Nil(t, err)
		for _, appProject := range appProjects {
			updatedAppProject, _ := argoCDClient.GetProject(context.TODO(), appProject.Name)
			assert.Nil(t, updatedAppProject.Spec.SyncWindows)
		}

		_, err = io.ReadAll(b)
		assert.Nil(t, err)
	})

	t.Run("Run set command", func(t *testing.T) {
		appProjects := client.GenerateTestProjects("./.env")
		defer client.DeleteTestProjects(appProjects, "./.env")

		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"set",
			"--server-address", "localhost:8080",
			"--plaintext",
			"--auth-token", os.Getenv("ARGOCD_TOKEN"),
			"--name", appProjects[0].Name,
			"--name", appProjects[1].Name,
			"--name", appProjects[2].Name,
			"--file", "./integration/exampleSyncWindows.json",
			"--label", "env=dev",
		})
		err := command.Execute()
		assert.Nil(t, err)

		assert.Nil(t, err)
		for index, appProject := range appProjects {
			updatedAppProject, _ := argoCDClient.GetProject(context.TODO(), appProject.Name)
			if index == 1 { // SyncWindow already exists
				assert.Len(t, updatedAppProject.Spec.SyncWindows, 3)
			} else {
				assert.Len(t, updatedAppProject.Spec.SyncWindows, 2)
			}
		}

	})
}
