package main

import (
	"github.com/marlinsearch/cli/marlin"
	"os"
)

func main() {
	marlin.Init(os.Args[1:])
}
