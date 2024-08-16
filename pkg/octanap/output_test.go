package octanap

import (
	"bytes"
	"io"
	"testing"

	"github.com/acarl005/stripansi"
	"github.com/stretchr/testify/assert"
)

func TestOutputs(t *testing.T) {
	b := bytes.NewBufferString("")

	octanap := New()
	octanap.Out = b
	octanap.Err = b

	testPhrase := "I'm a little hamster."

	t.Run("outputs string", func(t *testing.T) {
		octanap.Output(testPhrase)

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testPhrase+"\n", string(out))
	})

	t.Run("outputs header", func(t *testing.T) {
		octanap.OutputHeading(testPhrase)

		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, stripansi.Strip(string(out)), testPhrase)
	})
}