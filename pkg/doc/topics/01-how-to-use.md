---
Title: How to Use gitcommit
Slug: how-to-use
Short: Practical guide for using gitcommit to make safe, ticket-prefixed commits with optional docmgr changelog updates.
Topics:
- gitcommit
- git
- commit
- docmgr
- workflow
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: GeneralTopic
---

# How to Use gitcommit

## Overview

`gitcommit` is a small CLI that makes “safe commits” the default:

- It detects a ticket ID (`--ticket`, `GITCOMMIT_TICKET`, or branch name like `feature/ABC-123-...`)
- It ensures the commit summary starts with `TICKET: ...`
- It refuses to commit if no files are staged
- It blocks common noise artifacts (like `dist/`, `node_modules/`, `.env`) unless you explicitly override
- Optionally, it updates a `docmgr` ticket changelog for the commit

## Quick start

### 1) Stage your changes

```bash
git status --porcelain
git add path/to/files
```

### 2) Ensure a ticket is detectable

Option A: pass it explicitly:

```bash
gitcommit ticket --ticket ABC-123
```

Option B: set it for the session:

```bash
export GITCOMMIT_TICKET=ABC-123
gitcommit ticket
```

Option C: use a ticketed branch name:

```bash
git checkout -b feature/ABC-123-add-thing
gitcommit ticket
```

### 3) Preflight (recommended)

```bash
gitcommit preflight --docmgr=false
```

This checks:
- ticket resolution
- staged files
- noise policy
- (optionally) docmgr initialization + ticket existence

### 4) Commit

```bash
gitcommit commit -m "Fix widget ordering" --docmgr=false
```

If your ticket is `ABC-123`, the commit summary becomes:

```text
ABC-123: Fix widget ordering
```

## Noise policy (safety rail)

By default, `gitcommit commit` refuses to commit common noise files (examples):
- `dist/`, `build/`, `out/`
- `node_modules/`, `vendor/`
- `.env*`
- `*.log`, binaries, coverage output

Override only when you really mean it:

```bash
gitcommit commit -m "..." --allow-noise --docmgr=false
```

## docmgr integration

If you want automatic changelog entries in a `docmgr` ticket:

### 1) Initialize docmgr in your repo

```bash
gitcommit docmgr init
```

### 2) Create the ticket workspace (once)

```bash
gitcommit docmgr ticket create --ticket ABC-123 --title "..." --topics chat
```

### 3) Commit with docmgr enabled (default)

```bash
gitcommit commit -m "Fix widget ordering"
```

This will:
- create a git commit
- append an entry to the ticket changelog with the commit hash and file notes

Important: `docmgr` updates are *working tree changes*; commit them separately:

```bash
git add ttmp/...
git commit -m "Docs: update changelog (ABC-123)"
```

### 4) Run docmgr doctor (optional)

```bash
gitcommit docmgr doctor --ticket ABC-123
```

## Smoke tests (contributors)

`gitcommit` includes throwaway integration tests that create a temporary repo under `/tmp`:

```bash
bash test-scripts/test-cli.sh
bash test-scripts/test-all.sh
```

