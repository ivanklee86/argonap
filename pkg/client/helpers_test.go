package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpers(t *testing.T) {
	assert.Len(t, GenerateTestProjects(), 3)
}
