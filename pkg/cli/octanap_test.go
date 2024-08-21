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

func TestOctanapHappyPath(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		ServerAddr: "localhost:8080",
		Insecure:   true,
		AuthToken:  os.Getenv("ARGOCD_TOKEN"),
	}

	b := bytes.NewBufferString("")

	octanap := NewWithConfig(config)
	octanap.Out = b
	octanap.Err = b

	t.Run("Octonap can clear all SyncWindows", func(t *testing.T) {
		client.GenerateTestProjects()

		octanap.Connect()
		octanap.ClearSyncWindows()

		appProjects, err := octanap.ArgoCDClient.ListProjects(context.Background())
		assert.Nil(t, err)
		for _, appProject := range appProjects.Items {
			assert.Nil(t, appProject.Spec.SyncWindows)
		}
	})
}
