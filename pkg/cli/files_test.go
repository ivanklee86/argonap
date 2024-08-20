package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadingConfigFromFile(t *testing.T) {

	t.Run("Loads config file", func(t *testing.T) {
		syncWindows, err := readSyncWindowsFromFile("../../integration/exampleSyncWindows.json")

		assert.Nil(t, err)
		assert.Len(t, syncWindows, 2)
	})
}
