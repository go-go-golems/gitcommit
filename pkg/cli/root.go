package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Execute() error {
	return errors.Wrap(rootCmd.Execute(), "execute")
}

var rootCmd = &cobra.Command{
	Use:   "gitcommit",
	Short: "gitcommit helps craft consistent git commit messages",
}
