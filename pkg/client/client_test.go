package client

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func randomProjectName() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	length := 5

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[random.Intn(len(charset))]
	}

	return fmt.Sprintf("project%s", string(randomString))
}

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
		projectName := randomProjectName()
		_, err = client.CreateProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteProject(context.Background(), projectName)

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
		projectName := randomProjectName()

		project, err := client.CreateProject(context.Background(), projectName)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteProject(context.Background(), projectName)

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
