---
Title: Build a Go CLI to streamline safe git commits (ticket-aware docmgr integration)
Ticket: GITCOMMIT-XXXX
Status: complete
Topics:
    - go
    - cli
    - tooling
    - repo
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: .goreleaser.yaml
      Note: Release configuration
    - Path: Makefile
      Note: Build/lint/test targets
    - Path: README.md
      Note: Project overview
    - Path: cmd/gitcommit/cmds/commit/commit.go
      Note: Glazed commit command implementation
    - Path: cmd/gitcommit/cmds/docmgr/root.go
      Note: Glazed docmgr command group
    - Path: cmd/gitcommit/cmds/preflight/check.go
      Note: Glazed preflight checks
    - Path: cmd/gitcommit/cmds/root.go
      Note: Root cobra command + command tree registration
    - Path: cmd/gitcommit/cmds/ticket/show.go
      Note: Glazed ticket detection command
    - Path: cmd/gitcommit/main.go
      Note: |-
        CLI entry point
        Build info injection via ldflags
        CLI entry point (logging + help system init)
    - Path: go.mod
      Note: Go module root (github.com/go-go-golems/gitcommit)
    - Path: pkg/commitmsg/commitmsg.go
      Note: Ensure ticket prefix in commit summary
    - Path: pkg/doc/topics/01-how-to-use.md
      Note: End-user how-to guide (embedded help)
    - Path: pkg/docmgr/docmgr.go
      Note: |-
        Call docmgr to update ticket changelog
        Doctor wrapper and ticket create helper
    - Path: pkg/git/git.go
      Note: Wrap git commands (repo root
    - Path: pkg/layers/repository.go
      Note: Glazed layer for --repo flag
    - Path: pkg/ticket/ticket.go
      Note: Detect ticket ID from branch name
    - Path: pkg/validate/noise.go
      Note: Detect common noise artifacts in staged files
    - Path: test-scripts/README.md
      Note: How to run smoke tests
    - Path: test-scripts/setup-test-repo.sh
      Note: Create temp git+docmgr repo for smoke tests
    - Path: test-scripts/test-all.sh
      Note: End-to-end smoke test suite
    - Path: test-scripts/test-cli.sh
      Note: Minimal CLI smoke tests
ExternalSources: []
Summary: ""
LastUpdated: 2026-01-04T18:27:23.926420343-05:00
WhatFor: ""
WhenToUse: ""
---











# Build a Go CLI to streamline safe git commits (ticket-aware docmgr integration)

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- go
- cli
- tooling
- repo

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
