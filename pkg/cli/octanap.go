package cli

import (
	"io"
	"os"
)

type Config struct {
	NoExitCode bool
}

// Octanap is the logic/orchestrator.
type octanap struct {
	*Config

	// Allow swapping out stdout/stderr for testing.
	Out io.Writer
	Err io.Writer
}

// New returns a new instance of octanap.
func New() *octanap {
	config := Config{}

	return &octanap{
		Config: &config,
		Out:    os.Stdout,
		Err:    os.Stdin,
	}
}
