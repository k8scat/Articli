package utils

import (
	"fmt"
	"os"
)

func Err(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, f, a...)
	os.Exit(1)
}
