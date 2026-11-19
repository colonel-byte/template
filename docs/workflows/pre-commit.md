# Validate pre-commit hooks

[`.github/workflows/pre-commit.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/pre-commit.yaml) runs the hooks from [`.pre-commit-config.yaml`](https://github.com/OWNER/REPO/blob/main/.pre-commit-config.yaml) over the whole tree with [prek](https://github.com/j178/prek), on every pull request and on the merge queue.

It exists because the root [`AGENTS.md`](../AGENTS.md) tells everyone to run the hooks before committing and nothing otherwise checks that they did. It also closes a gap the hooks have by design: a hook that reads the staged diff only sees uncommitted changes, so a clean local run says nothing about code that is already committed. `--all-files` here does.

## What it skips, and why

| Hook             | Why it is skipped                                                                                                      |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `docs-generator` | It regenerates files rather than checking them, so in CI it would report a diff it created itself                      |
| `golangci-lint`  | Already covered by [Validate Lint](scan-lint.md), at the same pinned version. Running it twice doubles the slowest job |

`markdown-table-padding` and `verify-vendor` are deliberately **not** skipped. Both are checks rather than rewrites, and the table check is not covered anywhere else.

## `setup-go` before prek is not optional

The `repo: local` hooks are `language: system` and shell out to the Go toolchain. Without Go on `PATH` they fail as "command not found", which reads like a broken hook rather than a missing toolchain - so the workflow installs Go from `go-version-file: go.mod` before prek runs.

## Reproducing a failure

The same command, locally:

```sh
prek run --all-files --skip docs-generator --skip golangci-lint
```

`prek` and `pre-commit` read the same config, so either works; `prek` is faster. Several of the hooks rewrite files, so re-stage after a run - and if the failure is only in CI, check the hook versions in the config are the ones installed locally.
