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
echo "gitcommit - Complete Smoke Test Suite"
echo "=========================================="
echo ""

echo "Setting up test repository..."
bash "$SCRIPT_DIR/setup-test-repo.sh"
echo ""

echo "=== PHASE 1: Basic CLI ==="
gitcommit --help >/dev/null
gitcommit commit --help >/dev/null
gitcommit ticket --help >/dev/null
echo ""

echo "=== PHASE 2: Commit without docmgr ==="
hash1="$(gitcommit --repo "$TEST_REPO_DIR" commit --docmgr=false -m "No docmgr commit")"
subject1="$(git -C "$TEST_REPO_DIR" log -1 --pretty=%s)"
echo "$subject1" | rg -q "^${TICKET_ID}: "
echo "✓ Commit created: $hash1"
echo ""

echo "=== PHASE 3: Commit with docmgr changelog update ==="
echo "again" >> "$TEST_REPO_DIR/hello.txt"
git -C "$TEST_REPO_DIR" add hello.txt

if command -v docmgr >/dev/null; then
	hash2="$(gitcommit --repo "$TEST_REPO_DIR" commit -m "Docmgr updated commit")"
	# Find the ticket changelog path created by docmgr in this repo.
	changelog="$(find "$TEST_REPO_DIR/ttmp" -path "*${TICKET_ID}--*" -name changelog.md | head -1)"
	test -n "$changelog"
	rg -q "$hash2" "$changelog"
	echo "✓ docmgr changelog updated: $changelog"
else
	echo "SKIP: docmgr not found in PATH"
fi
echo ""

echo "=== PHASE 4: Noise-file rejection ==="
mkdir -p "$TEST_REPO_DIR/dist"
echo "noise" > "$TEST_REPO_DIR/dist/noise.bin"
git -C "$TEST_REPO_DIR" add "$TEST_REPO_DIR/dist/noise.bin"

set +e
gitcommit --repo "$TEST_REPO_DIR" commit --docmgr=false -m "Should fail due to noise" >/dev/null 2>&1
exit_code=$?
set -e

test $exit_code -ne 0
git -C "$TEST_REPO_DIR" reset -q HEAD -- "$TEST_REPO_DIR/dist/noise.bin"
rm -rf "$TEST_REPO_DIR/dist"
echo "✓ Noise-file check blocks commit by default"
echo ""

echo "=========================================="
echo "✓ Smoke tests passed"
echo "=========================================="
