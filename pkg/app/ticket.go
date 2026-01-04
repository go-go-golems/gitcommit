package app

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-go-golems/gitcommit/pkg/git"
	"github.com/go-go-golems/gitcommit/pkg/ticket"
)

const TicketEnvVar = "GITCOMMIT_TICKET"

type TicketResolution struct {
	TicketID string
	Source   string
	Branch   string
}

func ResolveTicket(ctx context.Context, repoRoot string, ticketOverride string) (TicketResolution, error) {
	if strings.TrimSpace(ticketOverride) != "" {
		return TicketResolution{
			TicketID: strings.ToUpper(strings.TrimSpace(ticketOverride)),
			Source:   "--ticket",
		}, nil
	}

	if env := strings.TrimSpace(os.Getenv(TicketEnvVar)); env != "" {
		return TicketResolution{
			TicketID: strings.ToUpper(env),
			Source:   "env:" + TicketEnvVar,
		}, nil
	}

	branch, err := git.CurrentBranch(ctx, repoRoot)
	if err != nil {
		return TicketResolution{}, errors.Wrap(err, "get current branch")
	}

	if t, ok := ticket.FromBranch(branch); ok {
		return TicketResolution{
			TicketID: t,
			Source:   "branch:" + branch,
			Branch:   branch,
		}, nil
	}

	return TicketResolution{}, errors.Errorf("could not detect ticket ID from branch %q (provide --ticket or set %s)", branch, TicketEnvVar)
}
