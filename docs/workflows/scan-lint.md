# Validate Lint

[`.github/workflows/scan-lint.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/scan-lint.yaml) runs golangci-lint over the module on every pull request and on the merge queue, against [`.golangci.yaml`](https://github.com/OWNER/REPO/blob/main/.golangci.yaml).

## The version is read, not written

This is the part worth understanding before changing the workflow. `golangci-lint-action` installs the **latest release** when given no version, so CI can lint with a newer golangci-lint than the hooks pin - and a new release that adds a check produces a failure that reproduces nowhere locally.

So the workflow reads the version out of `.pre-commit-config.yaml` with `awk` and passes it to the action. One pin, in one place:

- Renovate's `pre-commit` manager bumps the `rev:` in that file.
- Renovate's `github-actions` manager updates `uses:` references - **not** inputs inside a `with:` block.

A second pin written into the workflow is therefore one that nothing keeps in step with the first, which is exactly the drift this avoids. The `awk` fails the job if it cannot find the rev, rather than silently falling through to the next `repo:` entry and reporting some other tool's version back.

## What `.golangci.yaml` excludes, and why it has to

`magefiles/*.go` and `magefiles/core/` are excluded. Those files are `package main` behind `//go:build mage`, and mage supplies the main function at runtime, so the linter cannot build them - with or without `-tags mage`. Everything under `magefiles/pkg/` is ordinary Go and is linted like the rest of the module.

Two settings worth knowing before adding a `//nolint`:

- **`nolintlint` is on with `require-specific: true`**, so a bare `//nolint` fails. Name the linter.
- **`errcheck` has `check-blank: true`**, so `_ = f()` is an error too. The house answer is an entry in `errcheck.exclude-functions` for a function whose error is never actionable, not a suppression at each call site.

## Dockerfiles

Cargoship's version of this workflow carries a second `containers` job, one matrix entry per Dockerfile, running hadolint. It is left out here as a `TEMPLATE:` comment because there are no Dockerfiles yet. Note when adding it that hadolint's action takes a single path and its recursive mode globs from the repository root, which drags in the vendored Dockerfiles under `vendor/` - so the paths are listed explicitly.
