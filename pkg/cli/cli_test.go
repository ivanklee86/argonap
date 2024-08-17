package cli

import (
	"bytes"
	"testing"
)

func TestOctanapHappyPath(t *testing.T) {
	b := bytes.NewBufferString("")

	octanap := New()
	octanap.Out = b
	octanap.Err = b
}
