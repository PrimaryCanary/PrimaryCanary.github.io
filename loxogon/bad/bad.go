package bad

import (
	"fmt"
	"os"
)

var hadError = false

func Raise(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}

func Ok() bool {
	return hadError
}

func Reset() {
	hadError = false
}
