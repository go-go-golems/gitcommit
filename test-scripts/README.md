# test-scripts

These are manual smoke-test scripts for `gitcommit`. They create a temporary git repo under `/tmp` and run the CLI against it.

## Usage

```bash
cd gitcommit

# full suite
bash test-scripts/test-all.sh

# minimal checks
bash test-scripts/test-cli.sh
```

## Environment variables

- `TEST_REPO_DIR`: path for the temporary test repo (default: `/tmp/gitcommit-test-repo`)

