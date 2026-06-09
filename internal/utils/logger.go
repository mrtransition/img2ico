package utils

import (
	"fmt"
	"os"
)

var verboseEnabled = false

// EnableVerbose turns on verbose logging.
func EnableVerbose() {
	verboseEnabled = true
}

// Log prints a message only if verbose mode is on.
func Log(format string, args ...interface{}) {
	if verboseEnabled {
		fmt.Printf(format+"\n", args...)
	}
}

// FatalError prints an error message and exits with code 1.
func FatalError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
