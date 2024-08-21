package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ivanklee86/octanap/pkg/client"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	clientOptions := client.ArgoCDClientOptions{
		ServerAddr: "localhost:8080",
		Insecure:   true,
		AuthToken:  os.Getenv("ARGOCD_TOKEN"),
	}

	argoCDClient, err := client.New(&clientOptions)

	t.Run("Root command", func(t *testing.T) {
		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"--server-address", "localhost:8080",
			"--insecure",
			"--auth-token", os.Getenv("ARGOCD_TOKEN"),
		})
		err = command.Execute()
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
		client.GenerateTestProjects()

		b := bytes.NewBufferString("")

		command := NewRootCommand()
		command.SetOut(b)
		command.SetArgs([]string{
			"clear",
			"--server-address", "localhost:8080",
			"--insecure",
			"--auth-token", os.Getenv("ARGOCD_TOKEN"),
		})
		err = command.Execute()
		if err != nil {
			t.Fatal(err)
		}

		appProjects, err := argoCDClient.ListProjects(context.Background())
		assert.Nil(t, err)
		for _, appProject := range appProjects.Items {
			assert.Nil(t, appProject.Spec.SyncWindows)
		}

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Sprintln(string(out))
	})
}
