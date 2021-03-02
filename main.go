package main

import (
	"os"

	"github.com/eliostvs/nudolar/nudolar"
)

func main() {
	os.Exit(nudolar.CLI(os.Args[1:]))
}
