package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ivanklee86/octanap/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	argoCDClient := client.CreateTestClient()

	t.Run("Root command", func(t *testing.T) {
		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"--server-address", "localhost:8080",
			"--insecure",
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

		assert.Contains(t, string(out), "octanap")
	})

	t.Run("Run clear command", func(t *testing.T) {
		appProjects := client.GenerateTestProjects()

		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"clear",
			"--server-address", "localhost:8080",
			"--insecure",
			"--auth-token", os.Getenv("ARGOCD_TOKEN"),
		})
		err := command.Execute()
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		for _, appProject := range appProjects {
			updatedAppProject, _ := argoCDClient.GetProject(context.TODO(), appProject.Name)
			assert.Nil(t, updatedAppProject.Spec.SyncWindows)
		}

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Sprintln(string(out))
	})
}
