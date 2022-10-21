package logger

import (
	"os"

	"github.com/charlesbases/colors"
)

// Fatal .
func Fatal(v ...interface{}) {
	stderr(colors.RedSprint(v...))
}

// Fatalf .
func Fatalf(format string, v ...interface{}) {
	stderr(colors.RedSprintf(format, v...))
}

// stderr .
func stderr(err string) {
	os.Stderr.WriteString(colors.RedSprint("--apidoc_out: "))
	os.Stderr.WriteString(err)
	os.Stderr.WriteString("\n")
	os.Exit(1)
}
