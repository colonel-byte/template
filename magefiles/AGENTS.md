# Working in `magefiles/`

This directory houses the repository's task automation layer using Mage.

## Layout and organization

The top level of `magefiles/` is strictly for Mage namespace declarations, targets, and flag parsers in `package main`. Do not place heavy domain logic, parsers, HTTP scrapers, code generators, or template walkers directly in `package main` files at the root of `magefiles/`.

All implementation logic belongs in modular Go packages under `magefiles/pkg/`:

| Path                      | Purpose                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| `magefiles/dev.go`        | Mage entrypoint for the `Dev` namespace (`dev:checkTables`, `dev:vendor`) |
| `magefiles/generate.go`   | Mage entrypoint for the `Generate` namespace (`generate:document`)        |
| `magefiles/core/core.go`  | Bootstrapping main package for `go run ./magefiles/core`                  |
| `magefiles/pkg/gen/docs/` | The book's generators - the copied-in pages, and the table of contents    |
| `magefiles/pkg/util/`     | Shared helpers - currently the aligned-table check                        |

## Target structure rules

1. **Keep root namespace targets thin.** Root Mage methods should only perform argument defaults, environmental setup, or high-level status output, delegating immediately to the appropriate package under `magefiles/pkg/`.
2. **Encapsulate and test.** Logic under `magefiles/pkg/` should be written as standard Go packages with unit test files (`*_test.go`) adhering to the table-driven test formatting guidelines in the root [`AGENTS.md`](../AGENTS.md).

## The build tag, and why `go run ./magefiles/core` is the spelling to use

Every file at the root of `magefiles/` carries `//go:build mage`, which is how Mage finds targets and why a plain `go build ./...` does not compile this directory. `magefiles/core/core.go` carries no tag: it is an ordinary `main` that calls `mage.Main()`, so `go run ./magefiles/core <namespace>:<target>` runs the same targets as the `mage` CLI without needing it installed.

Use that form in CI, in hooks, and in documentation. It needs only the Go toolchain, and - because this repository vendors - it needs no network either.

## Everything builds from `vendor/`, with no network

`vendor/` is committed, and the build is expected to work on a machine that cannot reach a proxy. The Go toolchain selects `-mod=vendor` automatically when `vendor/` is present, so this is usually invisible - which is the problem: a dependency added without re-vendoring builds fine locally and fails only where it matters.

So a change to the dependency graph is not finished until `vendor/` is rewritten and committed in the same change:

```sh
go run ./magefiles/core dev:vendor        # go mod tidy, then go mod vendor
go run ./magefiles/core dev:verifyVendor  # go mod verify, then a GOPROXY=off build
```

`dev:verifyVendor` is the one that proves it: it builds with `GOFLAGS=-mod=vendor GOPROXY=off`, which is the invocation a release has to survive. [`.github/workflows/check-go-mod.yaml`](../.github/workflows/check-go-mod.yaml) runs the same target, and a `verify-vendor` pre-commit hook runs it whenever `go.mod`, `go.sum`, or `vendor/` changes.

Two consequences worth stating:

- **A new dependency is a bigger decision than it looks**, because it lands in `vendor/` and is carried by every clone of the repository from then on. The cost of a dependency is what it links, not what you call - check with `go list -deps ./...` before adding an import.
- **`vendor/` is generated, and marked so.** It is `linguist-generated` in [`.gitattributes`](../.gitattributes), and every content-reading pre-commit hook excludes it. A hook that reformats a vendored file makes the next `go mod vendor` look like a change nobody made.

### A new `.gitignore` pattern can silently empty part of `vendor/`

This is the failure worth knowing before writing one. `go mod vendor` copies whatever the dependencies ship, which includes files matching patterns written for this repository's own build output. An **unanchored** `dist/`, `bin/`, `*.out` or `*.tar.gz` matches at every depth, so it swallows a vendored package's files as readily as a local build artifact - and nothing reports it. The tree builds locally, because the files are on disk. It fails after a fresh clone, because they were never committed.

So a new ignore pattern is anchored with a leading `/` unless it genuinely has to match at every depth, and [`.gitignore`](../.gitignore) ends with `!/vendor/**` as the backstop. Two details about that line, both easy to get wrong:

- **It has to be last.** Git takes the last matching pattern, so a negation placed above the patterns it overrides does nothing.
- **`!/vendor/` is not the same thing.** The bare form un-ignores the directory entry and leaves every file inside it to be matched against the patterns above, which is no protection at all.

Check a change rather than reasoning about it - this lists every vendored file git is dropping, and should print nothing:

```sh
git ls-files --others --ignored --exclude-standard -- vendor/
```

## Every pin needs an owner

A pinned value is a version, tag, digest, or checksum written into the tree that names something living upstream: a container image digest, a package version, an upstream tag, a cloud image's `sha256`. Pinning is usually the right default - a build has to be reproducible - but a pin with nothing watching it goes stale silently. Nothing fails; the tree just keeps installing a version that stopped being current, and often stopped being patched, months ago.

So a new pin is not done when it is written. It is done when something can move it. Pick one of two owners and wire it up in the same change:

- **A mage target that rewrites the pin in place** - when the current value has to be *discovered* (queried from repodata, resolved from a registry, read from an upstream checksum file) and then written back into one or more files. The target does the query, the rewrite, and the reporting; the human runs it and reviews the diff. Where the pinned value is a whole list, the target needs a second half that *tells you to run it* - a verify target that re-fetches and exits non-zero on drift - because nobody spots a changed digest by reading a thousand-line file.
- **A [`.github/renovate.json5`](../.github/renovate.json5) entry** - when the pin sits in a file format Renovate already parses, so it costs a `packageRules` entry rather than a generator. `gomod`, `github-actions`, `pre-commit` and `npm` cover most of it; a pin in a Go constant, a Dockerfile, or a bespoke JSON manifest is reachable through the `customManagers` regexes, by annotating the line above it:

```go
// renovate: datasource=github-releases depName=owner/repo
const toolVersion = "1.2.3"
```

Prefer the Renovate entry when the format is already supported or can be annotated: it is a few lines, it opens the pull request unprompted, and it needs nobody to remember. Reach for a mage target when the value has to be computed rather than looked up.

The rule holds for values pinned outside `magefiles/` too, which is most of them. This section lives here because the automation does.
