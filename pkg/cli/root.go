package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/go-go-golems/gitcommit/pkg/commitmsg"
	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	"github.com/go-go-golems/gitcommit/pkg/ticket"
)

func Execute() error {
	return errors.Wrap(rootCmd.Execute(), "execute")
}

var rootCmd = &cobra.Command{
	Use:   "gitcommit",
	Short: "gitcommit helps craft consistent git commit messages",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repoDir, "repo", ".", "Working directory inside the git repo")

	rootCmd.AddCommand(newTicketCmd())
	rootCmd.AddCommand(newCommitCmd())
}

var repoDir string

func newTicketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ticket",
		Short: "Print the detected ticket ID (from --ticket/env/branch)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			repoRoot, err := git.RepoRoot(ctx, repoDir)
			if err != nil {
				return err
			}

			ticketID, source, err := resolveTicket(ctx, repoRoot, "")
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\t(%s)\n", ticketID, source)
			return nil
		},
	}
}

func newCommitCmd() *cobra.Command {
	var (
		message string
		body    []string
		ticketF string
		dryRun  bool
		useDoc  bool
	)

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit with ticket prefix and docmgr changelog update",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if strings.TrimSpace(message) == "" {
				return errors.New("missing --message")
			}

			repoRoot, err := git.RepoRoot(ctx, repoDir)
			if err != nil {
				return err
			}

			stagedFiles, err := git.StagedFiles(ctx, repoRoot)
			if err != nil {
				return err
			}
			if len(stagedFiles) == 0 {
				return errors.New("no staged files; stage changes first (git add ...)")
			}

			ticketID, source, err := resolveTicket(ctx, repoRoot, ticketF)
			if err != nil {
				return err
			}

			finalMessage := commitmsg.EnsureTicketPrefix(ticketID, message)

			if dryRun {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "repo:   %s\n", repoRoot)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "ticket: %s (%s)\n", ticketID, source)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "msg:    %s\n", finalMessage)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "files:  %d staged\n", len(stagedFiles))
				return nil
			}

			if err := git.Commit(ctx, repoRoot, finalMessage, body); err != nil {
				return err
			}

			hash, err := git.HeadHash(ctx, repoRoot)
			if err != nil {
				return err
			}

			if useDoc {
				if !docmgr.IsAvailable() {
					return errors.New("docmgr not found in PATH (use --docmgr=false to disable)")
				}

				if _, err := os.Stat(filepath.Join(repoRoot, ".ttmp.yaml")); err != nil {
					return errors.Wrap(err, "docmgr not initialized (missing .ttmp.yaml); run: docmgr init --seed-vocabulary")
				}

				exists, err := docmgr.TicketExists(ctx, repoRoot, ticketID)
				if err != nil {
					return err
				}
				if !exists {
					return errors.Errorf("docmgr ticket %s not found; create it with: docmgr ticket create-ticket --ticket %s --title \"...\" --topics ...", ticketID, ticketID)
				}

				fileNotes := make([]docmgr.FileNote, 0, len(stagedFiles))
				for _, f := range stagedFiles {
					fileNotes = append(fileNotes, docmgr.FileNote{
						Path: filepath.Join(repoRoot, f),
						Note: "Changed in commit " + hash,
					})
				}

				entry := fmt.Sprintf("Commit %s: %s", hash, finalMessage)
				if err := docmgr.ChangelogUpdate(ctx, repoRoot, ticketID, entry, fileNotes); err != nil {
					return err
				}
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "docmgr updated; commit docs changes separately (e.g. git add ttmp/... && git commit -m \"Docs: update changelog\")")
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", hash)
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit summary (ticket prefix will be ensured)")
	cmd.Flags().StringArrayVar(&body, "body", nil, "Commit body paragraph (repeatable)")
	cmd.Flags().StringVar(&ticketF, "ticket", "", "Ticket ID to use (overrides env/branch)")
	cmd.Flags().BoolVar(&useDoc, "docmgr", true, "Update docmgr changelog for the ticket")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print what would happen without committing")

	return cmd
}

func resolveTicket(ctx context.Context, repoRoot string, ticketFlag string) (ticketID string, source string, _ error) {
	if strings.TrimSpace(ticketFlag) != "" {
		return strings.ToUpper(strings.TrimSpace(ticketFlag)), "--ticket", nil
	}

	if env := strings.TrimSpace(os.Getenv("GITCOMMIT_TICKET")); env != "" {
		return strings.ToUpper(env), "env:GITCOMMIT_TICKET", nil
	}

	branch, err := git.CurrentBranch(ctx, repoRoot)
	if err != nil {
		return "", "", err
	}
	if t, ok := ticket.FromBranch(branch); ok {
		return t, "branch:" + branch, nil
	}

	return "", "", errors.Errorf("could not detect ticket ID from branch %q (provide --ticket or set GITCOMMIT_TICKET)", branch)
}
