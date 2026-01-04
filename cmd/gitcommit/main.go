package main

import (
	"fmt"
	"os"

	"github.com/go-go-golems/gitcommit/cmd/gitcommit/cmds"
	gitcommit_doc "github.com/go-go-golems/gitcommit/pkg/doc"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	help "github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmds.SetBuildInfo(version, commit, date)
	rootCmd := cmds.NewRootCmd()

	if err := logging.AddLoggingLayerToRootCommand(rootCmd, "gitcommit"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	helpSystem := help.NewHelpSystem()
	help_cmd.SetupCobraRootCommand(helpSystem, rootCmd)
	_ = helpSystem.LoadSectionsFromFS(gitcommit_doc.FS, "topics")

	if err := cmds.InitRootCmd(rootCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmds.Execute(rootCmd)
}
