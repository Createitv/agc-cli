package main

import (
	"fmt"
	"os"

	"github.com/your-org/agc-cli/cmd/agc/command"
)

func main() {
	if err := command.NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
