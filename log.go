package goutils

import (
	"log"
	"os"
)

var stdout = New(os.Stdout, "", log.LstdFlags)
var stderr = New(os.Stderr, "", log.LstdFlags)

func Printf(format string, v ...interface{}) {
	stdout.Output(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	stderr.Output(2, fmt.Sprintf(format, v...))
}
