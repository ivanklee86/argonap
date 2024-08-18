package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/joho/godotenv"
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
}
