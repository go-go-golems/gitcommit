# Tasks

## TODO

- [x] Add tasks here

- [x] Replace template placeholders (XXX) with gitcommit + create compilable CLI skeleton
- [x] Initialize docmgr, create ticket workspace + diary, and relate key files
- [x] Push bootstrap commits to origin
- [x] Define gitcommit CLI UX: commands, flags, commit message format
- [x] Implement ticket detection (flag/env/branch) + message prefixing
- [x] Implement 'gitcommit commit' with staged-file safety checks
- [x] Integrate docmgr: update ticket changelog with commit hash + file notes
- [x] Add unit tests for ticket detection + prefix logic
- [x] Update README with usage examples
- [x] Add smoke test scripts (create temp repo, exercise ticket/commit commands)
- [x] Run smoke test scripts and record results
- [x] Add noise-file safety check for staged files (block dist/, node_modules/, .env, etc.)
- [x] Extend smoke tests to cover noise-file rejection
- [x] Add preflight command to validate staged state/ticket/docmgr before commit
- [x] Add gitcommit docmgr wrappers (init, ticket create)
- [x] Refactor CLI into separate files (commit/ticket/preflight/docmgr)
- [x] Add version/build info support and document it
- [x] Add docmgr doctor/ticket-exists wrappers
- [x] Extend smoke scripts to cover version + docmgr doctor
