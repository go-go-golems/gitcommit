package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/go-go-golems/gitcommit/pkg/commitmsg"
	"github.com/go-go-golems/gitcommit/pkg/docmgr"
	"github.com/go-go-golems/gitcommit/pkg/git"
	"github.com/go-go-golems/gitcommit/pkg/validate"
)

func newCommitCmd() *cobra.Command {
	var (
		message string
		body    []string
		ticketF string
		dryRun  bool
		useDoc  bool
		allow   bool
	)

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit with ticket prefix and docmgr changelog update",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if strings.TrimSpace(message) == "" {
				return errors.New("missing --message")
			}

			root, err := repoRoot(ctx)
			if err != nil {
				return err
			}

			ticketID, source, err := resolveTicket(ctx, root, ticketF)
			if err != nil {
				return err
			}

			stagedFiles, err := git.StagedFiles(ctx, root)
			if err != nil {
				return err
			}
			if len(stagedFiles) == 0 {
				return errors.New("no staged files; stage changes first (git add ...)")
			}

			if !allow {
				noise := validate.FindNoise(stagedFiles)
				if len(noise) > 0 {
					return noiseError(noise)
				}
			}

			finalMessage := commitmsg.EnsureTicketPrefix(ticketID, message)

			if dryRun {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "repo:   %s\n", root)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "ticket: %s (%s)\n", ticketID, source)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "msg:    %s\n", finalMessage)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "files:  %d staged\n", len(stagedFiles))
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "docmgr: %v\n", useDoc)
				return nil
			}

			if useDoc {
				if !docmgr.IsAvailable() {
					return errors.New("docmgr not found in PATH (use --docmgr=false to disable)")
				}

				if _, err := os.Stat(filepath.Join(root, ".ttmp.yaml")); err != nil {
					return errors.Wrap(err, "docmgr not initialized (missing .ttmp.yaml); run: docmgr init --seed-vocabulary")
				}

				exists, err := docmgr.TicketExists(ctx, root, ticketID)
				if err != nil {
					return err
				}
				if !exists {
					return errors.Errorf("docmgr ticket %s not found; create it with: gitcommit docmgr ticket create --ticket %s --title \"...\" --topics ...", ticketID, ticketID)
				}
			}

			if err := git.Commit(ctx, root, finalMessage, body); err != nil {
				return err
			}

			hash, err := git.HeadHash(ctx, root)
			if err != nil {
				return err
			}

			if useDoc {
				fileNotes := make([]docmgr.FileNote, 0, len(stagedFiles))
				for _, f := range stagedFiles {
					fileNotes = append(fileNotes, docmgr.FileNote{
						Path: filepath.Join(root, f),
						Note: "Changed in commit " + hash,
					})
				}

				entry := fmt.Sprintf("Commit %s: %s", hash, finalMessage)
				if err := docmgr.ChangelogUpdate(ctx, root, ticketID, entry, fileNotes); err != nil {
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
	cmd.Flags().BoolVar(&allow, "allow-noise", false, "Allow common noise files (dist/, node_modules/, .env, etc.) to be committed")

	return cmd
}
