package main

import (
	"os"

	"github.com/deiwin/test-converter/convert"
)

func main() {
	convert.Test(os.Stdin, os.Stdout)
}
