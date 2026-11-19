# CodeQL

[`.github/workflows/codeql.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/codeql.yaml) runs GitHub's CodeQL analysis on pushes to `main`, on pull requests, on the merge queue, and weekly. Results land in the repository's Security tab. The weekly run is not redundant: it picks up queries and advisories published after a commit landed.

Two languages are analyzed, with different build modes:

| Language  | Build mode | What that means                                                    |
| --------- | ---------- | ------------------------------------------------------------------ |
| `go`      | `manual`   | The workflow's own `go build` decides which packages get extracted |
| `actions` | `none`     | The workflow files themselves are scanned, no build involved       |

## The Go build step is the package list

This is the part that needs maintaining. Because Go is analyzed by *building* it, CodeQL's `paths-ignore` setting does not apply - it has no effect on Go - so every exclusion has to be expressed as the set of packages the build step names.

Today that is `./magefiles/pkg/...` and nothing else, because `magefiles/*.go` is `package main` behind `//go:build mage` and mage supplies the main function at runtime, so `go build` cannot compile it at all. `./...` would pick it up and fail.

**Widen the list one package at a time as the repository grows** - `./cmd/...`, `./internal/...`, `./pkg/...`, and `.` for a root package. A pattern matching no package fails the build step, so adding a path before the package exists breaks the scan rather than skipping it.

## It builds offline, like everything else

`GOFLAGS=-mod=vendor`, `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`. Everything the build needs is vendored, so an accidental network access fails fast rather than quietly succeeding on a runner and failing in an air-gapped build. This matches the environment [`check-go-mod.yaml`](check-go-mod.md) enforces.

## Permissions

`security-events: write` to upload the results, `actions: read` for the `actions` language to read workflow files, and `packages: read` for any private module the build resolves. All three are on the job, not the file - see [`.github/AGENTS.md`](https://github.com/OWNER/REPO/blob/main/.github/AGENTS.md), where a top-level write zeroes Scorecard's Token-Permissions check.
