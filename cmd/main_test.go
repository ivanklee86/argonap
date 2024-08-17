package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	b := bytes.NewBufferString("")

	command := NewRootCommand()
	command.SetOut(b)
	command.SetArgs([]string{})
	err := command.Execute()
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
