package docmgr

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type FileNote struct {
	Path string
	Note string
}

func IsAvailable() bool {
	_, err := exec.LookPath("docmgr")
	return err == nil
}

func Init(ctx context.Context, repoRoot string) error {
	_, err := run(ctx, repoRoot, "init", "--seed-vocabulary")
	return err
}

func CreateTicket(ctx context.Context, repoRoot string, ticketID string, title string, topics []string) error {
	args := []string{"ticket", "create-ticket", "--ticket", ticketID, "--title", title}

	var cleaned []string
	for _, t := range topics {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		cleaned = append(cleaned, t)
	}
	if len(cleaned) > 0 {
		args = append(args, "--topics", strings.Join(cleaned, ","))
	}

	_, err := run(ctx, repoRoot, args...)
	return err
}

func TicketExists(ctx context.Context, repoRoot string, ticketID string) (bool, error) {
	out, err := run(ctx, repoRoot, "ticket", "list", "--ticket", ticketID)
	if err != nil {
		return false, err
	}

	if strings.Contains(out, "No tickets found.") {
		return false, nil
	}

	re := regexp.MustCompile("(?m)^###\\s+" + regexp.QuoteMeta(ticketID) + "\\b")
	return re.FindStringIndex(out) != nil, nil
}

func ChangelogUpdate(ctx context.Context, repoRoot string, ticketID string, entry string, fileNotes []FileNote) error {
	args := []string{"changelog", "update", "--ticket", ticketID, "--entry", entry}
	for _, fn := range fileNotes {
		note := strings.TrimSpace(fn.Note)
		if note == "" {
			note = "Changed"
		}

		absPath := fn.Path
		if !filepath.IsAbs(absPath) {
			absPath = filepath.Join(repoRoot, absPath)
		}

		args = append(args, "--file-note", absPath+":"+note)
	}

	_, err := run(ctx, repoRoot, args...)
	return err
}

func Doctor(ctx context.Context, repoRoot string, ticketID string, staleAfterDays int) (string, error) {
	args := []string{"doctor", "--ticket", ticketID}
	if staleAfterDays > 0 {
		args = append(args, "--stale-after", strconv.Itoa(staleAfterDays))
	}
	return run(ctx, repoRoot, args...)
}

func run(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "docmgr", args...)
	cmd.Dir = dir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "docmgr %s\n%s", strings.Join(args, " "), buf.String())
	}

	return buf.String(), nil
}
