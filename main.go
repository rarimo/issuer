package main

import (
	"os"

	"gitlab.com/q-dev/q-id/qid-issuer/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
