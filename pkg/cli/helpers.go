package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-go-golems/gitcommit/pkg/git"
	"github.com/go-go-golems/gitcommit/pkg/ticket"
)

func repoRoot(ctx context.Context) (string, error) {
	root, err := git.RepoRoot(ctx, repoDir)
	if err != nil {
		return "", err
	}
	return root, nil
}

func hasDocmgrConfig(repoRoot string) bool {
	_, err := os.Stat(filepath.Join(repoRoot, ".ttmp.yaml"))
	return err == nil
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
