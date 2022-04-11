package goutils

import (
	"fmt"
	"log"
	"os"
)

var stdout = log.New(os.Stdout, "", log.LstdFlags)
var stderr = log.New(os.Stderr, "", log.LstdFlags)

func Printf(format string, v ...interface{}) {
	stdout.Output(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	stderr.Output(2, fmt.Sprintf(format, v...))
}
