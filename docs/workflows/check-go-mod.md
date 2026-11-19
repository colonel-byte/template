# Check for go mod changes

[`.github/workflows/check-go-mod.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/check-go-mod.yaml) is what makes the vendored, offline build a requirement rather than a convention. It runs on every pull request and on the merge queue, and it fails for one reason: `go.mod`, `go.sum` and `vendor/` do not agree with each other.

## The two steps

- **`go mod tidy`, `go mod vendor`, then `git diff --exit-code go.mod go.sum vendor/`.** Re-resolves the graph and re-vendors from scratch, then fails if the result differs from what was committed. This catches a dependency added without running `mage dev:vendor`, and a `vendor/` left behind by a hand-edited `go.mod`.
- **`go run ./magefiles/core dev:verifyVendor`.** Runs `go mod verify`, then builds with `GOFLAGS=-mod=vendor GOPROXY=off` - the invocation a release has to survive. This catches a vendor tree that is internally consistent but incomplete.

The `go run ./magefiles/core` form rather than `mage`: same targets, builds from `vendor/`, needs no mage binary on the runner.

## Why both

They fail on different mistakes, and the second is the one that costs an afternoon. A `vendor/` missing a file still builds locally, because the file is on disk - it is only a fresh clone that notices. The commonest cause is not a forgotten `go mod vendor` at all, but a `.gitignore` pattern that quietly excluded part of the vendor tree; see "A new `.gitignore` pattern can silently empty part of `vendor/`" in [`magefiles/AGENTS.md`](https://github.com/OWNER/REPO/blob/main/magefiles/AGENTS.md).

That check is one command, and it is worth running by hand when a pattern changes:

```sh
git ls-files --others --ignored --exclude-standard -- vendor/
```

It should print nothing.

## When it fails on a bot's pull request

A dependency update re-vendors as part of the change, so the first step passes. Where it does not - a module that moved import paths, say - the fix is a human running `mage dev:vendor` on the branch, not a rerun.
