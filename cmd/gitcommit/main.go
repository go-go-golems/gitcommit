package main

import (
	"fmt"
	"os"

	"github.com/go-go-golems/gitcommit/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
