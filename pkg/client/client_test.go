package client //nolint:all

import (
	"context"
	"os"
	"testing"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/ivanklee86/argonap/pkg/testhelpers"
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
		assert.GreaterOrEqual(t, len(projects.Items), 2)

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
		project.Spec.SyncWindows = append(project.Spec.SyncWindows, &v1alpha1.SyncWindow{
			Kind:       "allow",
			Schedule:   "10 1 * * *",
			Duration:   "1h",
			Namespaces: []string{"*"},
		})

		_, err = client.UpdateProject(context.Background(), *project)
		if err != nil {
			t.Fatal(err)
		}

		project, err = client.GetProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, project.Annotations["test"], "value")
		assert.Len(t, project.Spec.SyncWindows, 1)

		project, err = client.GetProject(context.Background(), projectName)
		assert.Nil(t, err)

		project.Spec.SyncWindows = nil
		_, err = client.UpdateProject(context.Background(), *project)
		assert.Nil(t, err)

		project, err = client.GetProject(context.Background(), projectName)
		assert.Nil(t, err)
		assert.Len(t, project.Spec.SyncWindows, 0)
	})
}
