package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Execute() error {
	return errors.Wrap(rootCmd.Execute(), "execute")
}

var (
	repoDir string

	rootCmd = &cobra.Command{
		Use:   "gitcommit",
		Short: "gitcommit streamlines safe git commits and docmgr updates",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&repoDir, "repo", ".", "Working directory inside the git repo")

	rootCmd.AddCommand(newTicketCmd())
	rootCmd.AddCommand(newPreflightCmd())
	rootCmd.AddCommand(newCommitCmd())
	rootCmd.AddCommand(newDocmgrCmd())
}
