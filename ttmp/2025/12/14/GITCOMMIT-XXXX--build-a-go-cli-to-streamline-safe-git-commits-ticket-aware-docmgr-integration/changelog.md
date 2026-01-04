# Changelog

## 2026-01-04

- Initial workspace created


## 2026-01-04

Step 1: Bootstrap repo + Go module + CLI skeleton (commit fe657a0)

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/cmd/gitcommit/main.go — Add CLI entry point
- /home/manuel/workspaces/2026-01-04/gitcommit/go.mod — Initialize module github.com/go-go-golems/gitcommit
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/root.go — Add cobra root command


## 2026-01-04

Step 2: Initialize docmgr + create ticket workspace + diary

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/.ttmp.yaml — Initialize docmgr configuration
- /home/manuel/workspaces/2026-01-04/gitcommit/ttmp/2025/12/14/GITCOMMIT-XXXX--build-a-go-cli-to-streamline-safe-git-commits-ticket-aware-docmgr-integration/reference/01-diary.md — Create and start implementation diary

## 2026-01-04

Commit b6a781c7cef712999246591c0b187f001349777b: GITCOMMIT-XXXX: Add ticket-aware commit command

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/README.md — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/root.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/commitmsg/commitmsg.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/commitmsg/commitmsg_test.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/docmgr/docmgr.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/git/git.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/ticket/ticket.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/ticket/ticket_test.go — Changed in commit b6a781c7cef712999246591c0b187f001349777b

## 2026-01-04

Commit eccb6a8566c692eaba19c68aa725e526c941baeb: GITCOMMIT-XXXX: Allow --ticket override for ticket command

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/root.go — Changed in commit eccb6a8566c692eaba19c68aa725e526c941baeb


## 2026-01-04

Step 6: Add and run smoke test scripts

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-all.sh — End-to-end smoke test suite
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-cli.sh — Minimal CLI smoke test suite


## 2026-01-04

Commit 0798a66ca8e823270b93083323f7a4aff764b158: GITCOMMIT-XXXX: Add smoke test scripts

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/README.md — Changed in commit 0798a66ca8e823270b93083323f7a4aff764b158
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/setup-test-repo.sh — Changed in commit 0798a66ca8e823270b93083323f7a4aff764b158
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-all.sh — Changed in commit 0798a66ca8e823270b93083323f7a4aff764b158
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-cli.sh — Changed in commit 0798a66ca8e823270b93083323f7a4aff764b158


## 2026-01-04

Step 7: Block noise files by default (safe commit)

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/root.go — Enforce noise check and docmgr preflight
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/validate/noise.go — Noise detection


## 2026-01-04

Commit 9e93d831e726cd926f7687fb5ef9c77af37e37ee: GITCOMMIT-XXXX: Block common noise files by default

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/root.go — Changed in commit 9e93d831e726cd926f7687fb5ef9c77af37e37ee
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/validate/noise.go — Changed in commit 9e93d831e726cd926f7687fb5ef9c77af37e37ee
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-all.sh — Changed in commit 9e93d831e726cd926f7687fb5ef9c77af37e37ee


## 2026-01-04

Step 8: Add preflight + docmgr helper commands; refactor CLI

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/docmgr.go — Docmgr wrapper commands
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/preflight.go — New preflight command
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-all.sh — Smoke suite updated to use wrappers


## 2026-01-04

Commit 07141cb993f1f5284942f117f2d978952fe0d9bf: GITCOMMIT-XXXX: Add preflight and docmgr helper commands

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/README.md — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/commit.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/docmgr.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/helpers.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/preflight.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/root.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/ticket.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/docmgr/docmgr.go — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/setup-test-repo.sh — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-all.sh — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-cli.sh — Changed in commit 07141cb993f1f5284942f117f2d978952fe0d9bf


## 2026-01-04

Step 9: Add version/build info and docmgr doctor/exists wrappers

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/buildinfo.go — Version/build info
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/docmgr.go — Add doctor + ticket exists
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-cli.sh — Smoke test version flag


## 2026-01-04

Commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406: GITCOMMIT-XXXX: Add version and docmgr doctor/exists commands

### Related Files

- /home/manuel/workspaces/2026-01-04/gitcommit/.goreleaser.yaml — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/README.md — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/cmd/gitcommit/main.go — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/buildinfo.go — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/cli/docmgr.go — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/pkg/docmgr/docmgr.go — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-all.sh — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406
- /home/manuel/workspaces/2026-01-04/gitcommit/test-scripts/test-cli.sh — Changed in commit 82f6e2c6df18ffb74a50fe3b3e6f2bad28c37406


## 2026-01-04

Ticket complete: initial gitcommit MVP

