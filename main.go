package main

import (
	"os"

	"github.com/rarimo/issuer/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
