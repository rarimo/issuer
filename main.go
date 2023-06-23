package main

import (
	"os"

	"gitlab.com/rarimo/identity/issuer/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
