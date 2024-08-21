package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	b := bytes.NewBufferString("")

	command := NewRootCommand()
	command.SetOut(b)
	command.SetArgs([]string{
		"--server-address", "localhost:8080",
		"--insecure", "true",
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

	fmt.Print(string(out))

	assert.Contains(t, string(out), "octanap")
}
