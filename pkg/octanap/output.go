package octanap

import (
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	headerPrefix = "octanap"
)

// printToStream prints a generic message to a stream (stdout/stderror) in color.
func printToStream(stream io.Writer, msg interface{}) {
	_, err := fmt.Fprintf(stream, "%v\n", msg)
	if err != nil {
		panic(err)
	}
}

// printToStreamWithColor prints a message after wrapping it in ANSI color codes.
func printToStreamWithColor(stream io.Writer, color text.Color, msg interface{}) {
	_, err := fmt.Fprint(stream, color.Sprintf("%v\n", msg))
	if err != nil {
		panic(err)
	}
}

// OutputHeading prints a header to stdout.
func (octanap octanap) OutputHeading(msg interface{}) {
	printToStreamWithColor(octanap.Out, text.FgHiCyan, fmt.Sprintf("%v: %v", headerPrefix, msg))
}

// Output prints a normal message to stdout.
func (octanap octanap) Output(msg interface{}) {
	printToStream(octanap.Out, msg)
}

// Error pritns an error to stderr and exits with error code 1.
func (octanap octanap) Error(msg interface{}) {
	printToStreamWithColor(octanap.Err, text.FgHiRed, fmt.Sprintf("Error: %v\n", msg))
	if !octanap.NoExitCode {
		os.Exit(1)
	}
}
