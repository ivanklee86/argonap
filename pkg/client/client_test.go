package client

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestClinet(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	clientOptions := ArgoCDClientOptions{
		ServerAddr: "localhost:8080",
		Insecure:   true,
		AuthToken:  os.Getenv("ARGOCD_TOKEN"),
	}

	client, err := New(&clientOptions)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Can create projects", func(t *testing.T) {
		projectName := "octaproject"
		_, err = client.CreateProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}

		projects, err := client.ListProjects(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, projects.Items, 2)

		project, err := client.GetProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, project.ObjectMeta.Name, projectName)
	})
}
