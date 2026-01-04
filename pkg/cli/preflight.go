package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	"github.com/go-go-golems/gitcommit/pkg/validate"
)

type preflightOptions struct {
	ticketOverride string
	requireStaged  bool
	useDocmgr      bool
	allowNoise     bool
}

type preflightResult struct {
	repoRoot     string
	ticketID     string
	ticketSource string
	stagedFiles  []string
	noise        []validate.NoiseFinding
}

func newPreflightCmd() *cobra.Command {
	var opts preflightOptions

	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Check repository state before committing",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			res, err := runPreflight(ctx, opts)
			printPreflight(cmd, res, opts)
			return err
		},
	}

	cmd.Flags().StringVar(&opts.ticketOverride, "ticket", "", "Ticket ID to use (overrides env/branch)")
	cmd.Flags().BoolVar(&opts.requireStaged, "require-staged", true, "Require staged files")
	cmd.Flags().BoolVar(&opts.useDocmgr, "docmgr", true, "Require docmgr checks")
	cmd.Flags().BoolVar(&opts.allowNoise, "allow-noise", false, "Allow common noise files (dist/, node_modules/, .env, etc.) to be staged")

	return cmd
}

func runPreflight(ctx context.Context, opts preflightOptions) (preflightResult, error) {
	var res preflightResult

	root, err := repoRoot(ctx)
	if err != nil {
		return res, err
	}
	res.repoRoot = root

	ticketID, source, err := resolveTicket(ctx, root, opts.ticketOverride)
	if err != nil {
		return res, err
	}
	res.ticketID = ticketID
	res.ticketSource = source

	files, err := git.StagedFiles(ctx, root)
	if err != nil {
		return res, err
	}
	res.stagedFiles = files

	if opts.requireStaged && len(files) == 0 {
		return res, errors.New("no staged files; stage changes first (git add ...)")
	}

	if !opts.allowNoise {
		noise := validate.FindNoise(files)
		res.noise = noise
		if len(noise) > 0 {
			return res, noiseError(noise)
		}
	}

	if opts.useDocmgr {
		if !docmgr.IsAvailable() {
			return res, errors.New("docmgr not found in PATH (use --docmgr=false to disable)")
		}
		if !hasDocmgrConfig(root) {
			return res, errors.New("docmgr not initialized (missing .ttmp.yaml); run: docmgr init --seed-vocabulary")
		}

		exists, err := docmgr.TicketExists(ctx, root, ticketID)
		if err != nil {
			return res, err
		}
		if !exists {
			return res, errors.Errorf("docmgr ticket %s not found; create it with: gitcommit docmgr ticket create --ticket %s --title \"...\" --topics ...", ticketID, ticketID)
		}
	}

	return res, nil
}

func noiseError(noise []validate.NoiseFinding) error {
	var b strings.Builder
	b.WriteString("refusing to commit common noise files (use --allow-noise to override):\n")
	for _, n := range noise {
		b.WriteString("- ")
		b.WriteString(n.Path)
		b.WriteString(" (")
		b.WriteString(n.Reason)
		b.WriteString(")\n")
	}
	return errors.New(b.String())
}

func printPreflight(cmd *cobra.Command, res preflightResult, opts preflightOptions) {
	out := cmd.OutOrStdout()

	if res.repoRoot != "" {
		_, _ = fmt.Fprintf(out, "repo:   %s\n", res.repoRoot)
	}
	if res.ticketID != "" {
		_, _ = fmt.Fprintf(out, "ticket: %s (%s)\n", res.ticketID, res.ticketSource)
	}
	_, _ = fmt.Fprintf(out, "staged: %d\n", len(res.stagedFiles))
	_, _ = fmt.Fprintf(out, "docmgr: %v\n", opts.useDocmgr)
	if len(res.noise) > 0 {
		_, _ = fmt.Fprintf(out, "noise:  %d\n", len(res.noise))
	}
}
