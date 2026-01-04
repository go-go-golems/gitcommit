---
Title: Build a Go CLI to streamline safe git commits (ticket-aware docmgr integration)
Ticket: GITCOMMIT-XXXX
Status: active
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
    - Path: cmd/gitcommit/main.go
      Note: CLI entry point
    - Path: go.mod
      Note: Go module root (github.com/go-go-golems/gitcommit)
    - Path: pkg/cli/root.go
      Note: |-
        Cobra root command + Execute
        Implements gitcommit subcommands
    - Path: pkg/commitmsg/commitmsg.go
      Note: Ensure ticket prefix in commit summary
    - Path: pkg/docmgr/docmgr.go
      Note: Call docmgr to update ticket changelog
    - Path: pkg/git/git.go
      Note: Wrap git commands (repo root
    - Path: pkg/ticket/ticket.go
      Note: Detect ticket ID from branch name
ExternalSources: []
Summary: ""
LastUpdated: 2026-01-04T17:14:42.037821946-05:00
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
