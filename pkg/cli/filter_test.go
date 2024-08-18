package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSubset(t *testing.T) {
	bigMap := make(map[string]string)
	bigMap["a"] = "x"
	bigMap["b"] = "y"
	bigMap["c"] = "z"

	matchingMap := make(map[string]string)
	matchingMap["a"] = "x"
	matchingMap["b"] = "y"

	notMatchingMap := make(map[string]string)
	notMatchingMap["d"] = "u"

	differentValueMap := make(map[string]string)
	differentValueMap["a"] = "c"

	t.Run("Map contains a submap", func(t *testing.T) {
		assert.Equal(t, true, isMapSubset(bigMap, matchingMap))
	})

	t.Run("Map does not contains a submap", func(t *testing.T) {
		assert.Equal(t, false, isMapSubset(bigMap, notMatchingMap))
		assert.Equal(t, false, isMapSubset(bigMap, differentValueMap))
	})
}
