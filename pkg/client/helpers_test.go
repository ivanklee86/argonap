package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpers(t *testing.T) {
	testProjects := GenerateTestProjects()
	assert.Len(t, testProjects, 3)
	DeleteTestProjects(testProjects)
}
