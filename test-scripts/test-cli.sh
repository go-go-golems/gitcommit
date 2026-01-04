#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

TEST_REPO_DIR="${TEST_REPO_DIR:-/tmp/gitcommit-test-repo}"
TICKET_ID="${TICKET_ID:-GITCOMMIT-XXXX}"

gitcommit() {
	(
		cd "$REPO_ROOT" && go run ./cmd/gitcommit "$@"
	)
}

echo "=========================================="
echo "gitcommit - CLI Smoke Test"
echo "=========================================="
echo ""

echo "Setting up test repository..."
bash "$SCRIPT_DIR/setup-test-repo.sh"
echo ""

echo "Test 1: Help"
gitcommit --help >/dev/null
gitcommit commit --help >/dev/null
gitcommit ticket --help >/dev/null
echo "✓ Help works"
echo ""

echo "Test 2: Ticket detection (branch)"
gitcommit --repo "$TEST_REPO_DIR" ticket | rg -q "$TICKET_ID"
echo "✓ Ticket detection works"
echo ""

echo "Test 3: Ticket override flag"
gitcommit --repo "$TEST_REPO_DIR" ticket --ticket "ABC-123" | rg -q "ABC-123"
echo "✓ Ticket override works"
echo ""

echo "Test 4: Dry-run commit (does not create commit)"
before="$(git -C "$TEST_REPO_DIR" rev-parse HEAD)"
gitcommit --repo "$TEST_REPO_DIR" commit --dry-run -m "Test dry run" >/dev/null
after="$(git -C "$TEST_REPO_DIR" rev-parse HEAD)"
test "$before" = "$after"
echo "✓ Dry-run does not commit"
echo ""

echo "=========================================="
echo "All tests passed ✓"
echo "=========================================="

