package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/go-go-golems/gitcommit/pkg/docmgr"
)

func newDocmgrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docmgr",
		Short: "Helpers around docmgr (init, ticket create)",
	}

	cmd.AddCommand(newDocmgrInitCmd())
	cmd.AddCommand(newDocmgrTicketCmd())

	return cmd
}

func newDocmgrInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Run docmgr init --seed-vocabulary in the repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			root, err := repoRoot(ctx)
			if err != nil {
				return err
			}
			if !docmgr.IsAvailable() {
				return errors.New("docmgr not found in PATH")
			}
			return docmgr.Init(ctx, root)
		},
	}
}

func newDocmgrTicketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ticket",
		Short: "Create/check docmgr tickets",
	}

	cmd.AddCommand(newDocmgrTicketCreateCmd())
	return cmd
}

func newDocmgrTicketCreateCmd() *cobra.Command {
	var (
		ticketF string
		title   string
		topics  []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a docmgr ticket workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			root, err := repoRoot(ctx)
			if err != nil {
				return err
			}

			ticketID, _, err := resolveTicket(ctx, root, ticketF)
			if err != nil {
				return err
			}

			if title == "" {
				return errors.New("missing --title")
			}

			if !docmgr.IsAvailable() {
				return errors.New("docmgr not found in PATH")
			}

			return docmgr.CreateTicket(ctx, root, ticketID, title, topics)
		},
	}

	cmd.Flags().StringVar(&ticketF, "ticket", "", "Ticket ID to create (defaults to env/branch detection)")
	cmd.Flags().StringVar(&title, "title", "", "Ticket title")
	cmd.Flags().StringSliceVar(&topics, "topics", []string{"chat"}, "Topics for the ticket (default: chat)")

	return cmd
}
