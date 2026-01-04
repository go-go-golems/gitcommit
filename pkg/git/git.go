package git

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func RepoRoot(ctx context.Context, dir string) (string, error) {
	out, err := run(ctx, dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", errors.Wrap(err, "get repo root")
	}
	return strings.TrimSpace(out), nil
}

func CurrentBranch(ctx context.Context, dir string) (string, error) {
	out, err := run(ctx, dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.Wrap(err, "get current branch")
	}
	return strings.TrimSpace(out), nil
}

func StagedFiles(ctx context.Context, dir string) ([]string, error) {
	out, err := run(ctx, dir, "diff", "--cached", "--name-only", "-z")
	if err != nil {
		return nil, errors.Wrap(err, "list staged files")
	}

	var files []string
	for _, part := range strings.Split(out, "\x00") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		files = append(files, filepath.Clean(part))
	}
	return files, nil
}

func Commit(ctx context.Context, dir string, message string, body []string) error {
	args := []string{"commit", "-m", message}
	for _, b := range body {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		args = append(args, "-m", b)
	}
	_, err := run(ctx, dir, args...)
	return errors.Wrap(err, "git commit")
}

func HeadHash(ctx context.Context, dir string) (string, error) {
	out, err := run(ctx, dir, "rev-parse", "HEAD")
	if err != nil {
		return "", errors.Wrap(err, "get HEAD hash")
	}
	return strings.TrimSpace(out), nil
}

func run(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "git %s\n%s", strings.Join(args, " "), buf.String())
	}

	return buf.String(), nil
}
