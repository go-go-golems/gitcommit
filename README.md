# gitcommit

gitcommit streamlines safe git commits by:
- ensuring a ticket prefix in the commit summary, and
- optionally updating `docmgr` changelog entries for that ticket.

## Usage

Detect the ticket ID (from `--ticket`, `GITCOMMIT_TICKET`, or current git branch):

`gitcommit ticket`

Commit (requires staged files):

`gitcommit commit --ticket ABC-123 -m "Fix widget ordering"`

This will produce a commit summary like:

`ABC-123: Fix widget ordering`

### docmgr integration

By default, `gitcommit commit` updates `docmgr` via `docmgr changelog update --ticket <TICKET> ...`.

- Disable with `--docmgr=false`
- `docmgr` must be installed and the repo must be initialized (`.ttmp.yaml` present)

## Development

- Test: `go test ./...`
- Lint: `make lint`
- Run: `go run ./cmd/gitcommit`
