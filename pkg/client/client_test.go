package client

import (
	"context"
	"os"
	"testing"

	"github.com/ivanklee86/octanap/pkg/testhelpers"
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
		projectName := testhelpers.RandomProjectName()
		_, err = client.CreateProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteProject(context.Background(), projectName) //nolint:all

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

	t.Run("Can update projects", func(t *testing.T) {
		projectName := testhelpers.RandomProjectName()

		project, err := client.CreateProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteProject(context.Background(), projectName) //nolint:all

		project.Annotations = make(map[string]string)
		project.Annotations["test"] = "value"

		_, err = client.UpdateProject(context.Background(), *project)
		if err != nil {
			t.Fatal(err)
		}

		project, err = client.GetProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, project.Annotations["test"], "value")
	})
}
