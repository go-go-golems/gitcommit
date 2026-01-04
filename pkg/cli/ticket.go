package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newTicketCmd() *cobra.Command {
	var ticketF string

	cmd := &cobra.Command{
		Use:   "ticket",
		Short: "Print the detected ticket ID (from --ticket/env/branch)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			repoRoot, err := repoRoot(ctx)
			if err != nil {
				return err
			}

			ticketID, source, err := resolveTicket(ctx, repoRoot, ticketF)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\t(%s)\n", ticketID, source)
			return nil
		},
	}

	cmd.Flags().StringVar(&ticketF, "ticket", "", "Ticket ID to use (overrides env/branch)")
	return cmd
}
