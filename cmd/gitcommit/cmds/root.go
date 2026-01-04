package cmds

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-go-golems/gitcommit/cmd/gitcommit/cmds/commit"
	"github.com/go-go-golems/gitcommit/cmd/gitcommit/cmds/docmgr"
	"github.com/go-go-golems/gitcommit/cmd/gitcommit/cmds/preflight"
	"github.com/go-go-golems/gitcommit/cmd/gitcommit/cmds/ticket"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gitcommit",
		Short:   "Streamline safe git commits (ticket-aware) with docmgr integration",
		Long:    "gitcommit streamlines safe commits by enforcing ticket prefixes, running preflight checks, and optionally updating docmgr changelogs.",
		Version: buildInfo.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return logging.InitLoggerFromCobra(cmd)
		},
	}

	// Global flags (parsed by glazed layers via ExistingCobraFlagsLayer).
	rootCmd.PersistentFlags().StringP("repo", "r", ".", "Path to git repository")

	return rootCmd
}

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

var buildInfo = BuildInfo{
	Version: "dev",
	Commit:  "none",
	Date:    "unknown",
}

func SetBuildInfo(version, commitHash, date string) {
	if version != "" {
		buildInfo.Version = version
	}
	if commitHash != "" {
		buildInfo.Commit = commitHash
	}
	if date != "" {
		buildInfo.Date = date
	}
}

var initOnce sync.Once
var initErr error

func InitRootCmd(rootCmd *cobra.Command) error {
	initOnce.Do(func() {
		if rootCmd == nil {
			initErr = errors.New("rootCmd is nil")
			return
		}

		// Explicit initialization of subcommand trees (no init() ordering reliance).
		if err := ticket.Init(); err != nil {
			initErr = errors.Wrap(err, "failed to init ticket commands")
			return
		}
		if err := preflight.Init(); err != nil {
			initErr = errors.Wrap(err, "failed to init preflight commands")
			return
		}
		if err := commit.Init(); err != nil {
			initErr = errors.Wrap(err, "failed to init commit commands")
			return
		}
		if err := docmgr.Init(); err != nil {
			initErr = errors.Wrap(err, "failed to init docmgr commands")
			return
		}

		rootCmd.AddCommand(
			ticket.TicketCmd,
			preflight.PreflightCmd,
			commit.CommitCmd,
			docmgr.DocmgrCmd,
		)
	})
	return initErr
}

func Execute(rootCmd *cobra.Command) {
	rootCmd.SetVersionTemplate(fmt.Sprintf("{{with .Name}}{{printf \"%%s \" .}}{{end}}{{printf \"%%s\\n\" .Version}}commit: %s\\ndate: %s\\n", buildInfo.Commit, buildInfo.Date))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
