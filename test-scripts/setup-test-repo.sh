#!/usr/bin/env bash
set -euo pipefail

TEST_DIR="${TEST_REPO_DIR:-/tmp/gitcommit-test-repo}"
TICKET_ID="${TICKET_ID:-GITCOMMIT-XXXX}"

echo "Creating test repository at $TEST_DIR..."

rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

git init -b main >/dev/null
git config user.email "test@example.com"
git config user.name "Test User"

cat > README.md << 'EOF'
# gitcommit smoke test repo
EOF

echo "hello" > hello.txt

git add .
git commit -m "Initial commit" >/dev/null

git checkout -b "feature/${TICKET_ID}-smoke" >/dev/null

echo "more" >> hello.txt
git add hello.txt

echo ""
echo "Test repository created successfully at: $TEST_DIR"
echo "Branch: $(git rev-parse --abbrev-ref HEAD)"
echo ""
