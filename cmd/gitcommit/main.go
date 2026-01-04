package main

import (
	"fmt"
	"os"

	"github.com/go-go-golems/gitcommit/pkg/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.SetBuildInfo(version, commit, date)

	if err := cli.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
