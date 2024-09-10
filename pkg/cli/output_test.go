package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestOutputs(t *testing.T) {
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

	argonap := NewWithConfig(config)
	argonap.Out = b
	argonap.Err = b

	testPhrase := "I'm a little hamster."

	t.Run("outputs string", func(t *testing.T) {
		argonap.Output(testPhrase)

		out, err := io.ReadAll(b)
		assert.Nil(t, err)
		assert.Equal(t, testPhrase+"\n", string(out))
	})

	t.Run("outputs header", func(t *testing.T) {
		argonap.OutputHeading(testPhrase)

		out, err := io.ReadAll(b)
		assert.Nil(t, err)
		assert.Contains(t, stripansi.Strip(string(out)), testPhrase)
	})

	t.Run("outputs successful worker results", func(t *testing.T) {
		sampleResult := WorkerResult{
			Status:      StatusSuccess,
			ProjectName: "foobar",
		}

		argonap.OutputResult(sampleResult)

		out, err := io.ReadAll(b)
		assert.Nil(t, err)
		assert.Contains(t, stripansi.Strip(string(out)), "Success")
	})

	t.Run("outputs failed worker results", func(t *testing.T) {
		sampleError := errors.New("test")

		sampleResult := WorkerResult{
			Status:      StatusFailure,
			ProjectName: "foobar",
			Err:         &sampleError,
		}

		argonap.OutputResult(sampleResult)

		out, err := io.ReadAll(b)
		assert.Nil(t, err)
		assert.Contains(t, stripansi.Strip(string(out)), "Failure")
	})
}
