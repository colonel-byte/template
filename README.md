# colonel-byte repository template
[![Latest Release](https://img.shields.io/github/v/release/OWNER/REPO)](https://github.com/OWNER/REPO/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/OWNER/REPO/release.yaml)](https://github.com/OWNER/REPO/actions/workflows/release.yaml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/OWNER/REPO/badge)](https://securityscorecards.dev/viewer/?uri=github.com/OWNER/REPO)

The agent rules and GitHub templates every repository in this organization starts from. Nothing here is project-specific: the cargoship-only sections have been dropped and the names that remain are placeholders.

## What is here

| Path                                 | What it is                                                                                    |
| ------------------------------------ | --------------------------------------------------------------------------------------------- |
| `AGENTS.md`                          | The root rules an agent loads at the start of a session. `CLAUDE.md` is a symlink to it       |
| `docs/AGENTS.md`                     | Markdown formatting, mdBook redirects, release-please annotations, generated pages            |
| `.github/AGENTS.md`                  | Scorecard-safe workflow changes, issue-form upkeep, the pull request description format       |
| `.github/pull_request_template.md`   | The three headings every pull request body keeps                                              |
| `.github/ISSUE_TEMPLATE/`            | Six YAML issue forms plus the chooser config                                                  |
| `docs/css/`                          | The mdBook stylesheets - `root.css` for the site, `print.css` for the single-page print view  |
| `.github/renovate.json5`             | Renovate config, in the house shape - digest pinning, grouped updates, custom regex pins      |
| `.github/workflows/renovate.yaml`    | Self-hosted Renovate, with the release-please app's credentials                               |
| `release-please-config.json`         | Versioning and changelog sections. `.release-please-manifest.json` holds the version itself   |
| `.goreleaser.yaml`                   | One build, archives, checksums, SBOMs, cosign signing. No packages or images                  |
| `.github/workflows/`                 | `release-please.yaml` opens the release PR; `release.yaml` builds on the tag it pushes        |
| `.github/workflows/scorecard.yaml`   | OpenSSF Scorecard on `main`, weekly, and on branch-protection changes. Feeds the badge        |
| `.github/workflows/commitlint.yaml`  | PR title linting, with the pinned commitlint in `.github/commitlint/`                         |
| `.github/workflows/pre-commit.yaml`  | The hooks, over the whole tree, so a local skip cannot hide                                   |
| `.github/workflows/scan-lint.yaml`   | golangci-lint at the version `.pre-commit-config.yaml` pins, read out of it                   |
| `.github/workflows/codeql.yaml`      | CodeQL for `go` and `actions`, weekly and per pull request                                    |
| `.golangci.yaml`                     | The house linter set. `magefiles/*.go` excluded - the build tag makes it uncompilable         |
| `.github/workflows/deploy-book.yaml` | mdBook to `gh-pages` on every push to `main`. Needs Pages enabled, which is a setting         |
| `.github/SECURITY.md`                | Security policy. The two `TEMPLATE:` blocks in it are the scope lists, which must be written  |
| `.github/CODEOWNERS`                 | Review from `@colonel-byte/maintainers` on every change                                       |
| `.gitattributes`                     | What each generator owns, grouped by the target that writes it                                |
| `.pre-commit-config.yaml`            | The org-standard hooks, plus three `local` ones that call mage. Three more are named in it    |
| `magefiles/`                         | A thin task runner - `generate:document`, `dev:checkTables`, `dev:vendor`, `dev:verifyVendor` |
| `go.mod`, `go.sum`, `vendor/`        | The vendored, offline build. `mage` and `testify` only                                        |
| `.gitignore`                         | Build output and `node_modules/`. `vendor/` is deliberately not in it                         |
| `docs/{guides,dev,agent}/`           | A `README.md` each, saying what belongs in that directory and what does not                   |
| `docs/workflows/`                    | One page per workflow, which `.github/AGENTS.md` requires                                     |

## The task runner, and the offline build

`magefiles/` is a thin version of cargoship's: four targets, `mage` and `testify` as the only dependencies, no distro or packaging generators. It exists because about ten rules across `AGENTS.md`, `docs/AGENTS.md` and `.github/AGENTS.md` describe maintaining generated output and defer to a task runner, and the "Pinned versions" section defers entirely to `magefiles/AGENTS.md`.

```sh
go run ./magefiles/core -l                  # list every target
go run ./magefiles/core generate:document   # docs/index.md, docs/security.md, docs/SUMMARY.md
go run ./magefiles/core dev:checkTables     # the Markdown table padding rule, enforced
go run ./magefiles/core dev:vendor          # go mod tidy, then go mod vendor
go run ./magefiles/core dev:verifyVendor    # go mod verify, then a GOPROXY=off build
```

Three of those are wired into `.pre-commit-config.yaml` as `local` hooks, so the generated pages, the table padding, and `vendor/` are checked at commit time rather than discovered in CI.

**The build is vendored and offline by requirement, not by habit.** `vendor/` is committed, the Go toolchain picks `-mod=vendor` automatically when it is present, and `dev:verifyVendor` builds with `GOPROXY=off` to prove a release can be built with no network. A dependency added without re-vendoring builds fine locally and fails only where it matters, so three things catch it: the `verify-vendor` hook, `.github/workflows/check-go-mod.yaml` re-running tidy and vendor and failing on any diff, and `dev:verifyVendor` in that same job.

Two files are now generator output rather than hand-written: `docs/SUMMARY.md` (compiled from the `docs/` tree) and `docs/index.md` / `docs/security.md` (copied from `README.md` and `.github/SECURITY.md` with links rebased). All are marked `linguist-generated`. Edit the source, then regenerate.

## CI

Ten workflows, each with a page under `docs/workflows/`. Five run on every pull request - `Validate Lint`, `Validate pre-commit hooks`, `CodeQL`, `PR Title Check`, `Check for go mod changes` - and the rest fire on a push to `main`, a tag, or a schedule.

**No test workflow ships.** A green test check over a repository with no tests means nothing, and the decisions one encodes - which packages, `-short` or not, which backend provisions the fixture - are not templatable. [`docs/dev/adding-ci.md`](docs/dev/adding-ci.md) is the note for writing the unit-test and end-to-end workflows, including the three things that change in the same pull request when you do.

One coupling worth knowing before editing either file: **`scan-lint.yaml` reads the golangci-lint version out of `.pre-commit-config.yaml`** rather than pinning its own. Renovate's `pre-commit` manager bumps that `rev:`; its `github-actions` manager updates `uses:` references, not inputs inside a `with:` block. A second pin in the workflow would be one nothing keeps in step with the first.

## Badges

The three at the top of this file are the house set, and each one needs something to exist before it resolves:

| Badge             | Resolves once                                                             | Fails as      |
| ----------------- | ------------------------------------------------------------------------- | ------------- |
| Latest Release    | A release has been published, which means `release.yaml` has run on a tag | `no releases` |
| Build Status      | `.github/workflows/release.yaml` exists under that exact filename         | `no status`   |
| OpenSSF Scorecard | A Scorecard run has published a score for the repository                  | `unknown`     |

They carry `OWNER/REPO` like everything else in the template, so they render as "repo not found" until the bootstrap grep replaces them - which is deliberate, since a badge is the most visible thing on the page and a wrong slug in one silently reports another repository's health as this one's.

Two things to know:

- **The Build Status badge is keyed on the workflow filename**, not its `name:`. Rename `release.yaml` and the badge goes to `no status` without anything else breaking.
- **The Scorecard badge needs `publish_results: true` and a public repository.** `workflows/scorecard.yaml` sets it, but the workflow can otherwise run happily and leave the badge on `unknown` - see [`docs/workflows/scorecard.md`](docs/workflows/scorecard.md). The first run on `main` is what fills it in.

Cargoship carries a fourth, keyed on `go-mod/go-version`, which this template now has a `go.mod` for. Add it if you want it:

```markdown
[![Go version](https://img.shields.io/github/go-mod/go-version/OWNER/REPO?filename=go.mod)](https://go.dev/)
```

## Releases

Release-please opens one pull request holding the next version and the regenerated changelog; merging it pushes a `vN.N.N` tag, and that tag is what goreleaser fires on. Both are deliberately basic - the cargoship versions also build `.deb`/`.rpm`/`.apk` packages, four container images, and an Ansible collection tarball, none of which is templatable. Read [`docs/workflows/release-please.md`](docs/workflows/release-please.md) and [`docs/workflows/release.md`](docs/workflows/release.md), then:

- **Set `package-name` and `extra-files` in `release-please-config.json`.** `extra-files` starts empty, and a page pinning this project's own version has to be listed there *and* annotated in place or it goes stale silently.
- **Check the `-ldflags -X` symbol paths in `.goreleaser.yaml` against the code.** A `-X` naming a variable that does not exist is ignored without complaint, and the binary reports no version.
- **The `COSIGN_KEY` and `COSIGN_PASS` secrets are required**, since `cosign.pub` at the root is the published half of that keypair. There is no unsigned fallback. The workflow maps them onto `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD`, which are cosign's own variable names.
- **Leave `.release-please-manifest.json` and `CHANGELOG.md` alone.** Release-please owns both and rewrites them wholesale.

## Renovate

The workflow runs Renovate itself rather than installing the hosted `renovate[bot]` app, so the pull requests are opened with the release-please app's installation token. Read [`docs/workflows/renovate.md`](docs/workflows/renovate.md) before enabling it - three things bite:

- **The app needs `Workflows: write`**, or every GitHub Actions pin update fails and nothing else does.
- **The author is the app's bot, not `renovate[bot]`.** Any `if: github.actor == 'renovate[bot]'` gate, author-selected commitlint config, or bot-pinned commit email stops matching, silently.
- **`automerge` is off**, explicitly rather than by default, so a preset added to `extends` cannot turn it back on. The grouped rules mean one pull request per family, each reviewed and merged by a person.

## Wiring up `docs/css/`

mdBook loads neither file unless `book.toml` names it, and a stylesheet nothing references fails silently - the book builds and looks untouched. Add both to `[output.html]`:

```toml
[output.html]
additional-css = [ "docs/css/root.css", "docs/css/print.css" ]
```

`print.css` carries two blocks keyed on mdBook heading slugs, marked `TEMPLATE:` in the file: one shrinks the generated reference chapters, the other drops the developer-facing tail from the print page. Both match on an ID, so a wrong slug matches nothing and reports nothing. Build the book, open the print page, and read the real IDs off the headings before editing them - or delete both blocks, which leaves the rest of the print styling working.

## Filling it in

Generating a repository from this template copies **files, not settings** — so `.github/CODEOWNERS` requires nothing until a ruleset says so, the workflows have no secrets, the labels the issue forms name do not exist, and the book 404s until Pages is enabled. [`docs/dev/new-repo-checklist.md`](docs/dev/new-repo-checklist.md) is the list, with the `gh` commands for the scriptable parts and the order to do the first run in. The placeholder sweep below comes first.

Two placeholders appear throughout. Replace both before the first commit in a new repository:

```sh
grep -rn 'OWNER/REPO\|OWNER\.github\.io\|<project>' . --exclude-dir=.git
```

- `OWNER/REPO` - the repository slug, in the issue-tracker section, the Scorecard links, the advisory URL, the `gh` snippets, and **`book.toml`**.
- `OWNER.github.io` - the published book's host, in `docs/AGENTS.md` and the issue chooser.
- `<project>` - the binary and module name, in the layout table, the `go list -deps` example, and the help-text examples.

`book.toml` is the one to check by eye: `git-repository-url` and `edit-url-template` are the repository icon and the per-page edit pencil, and pointed at the wrong slug they build cleanly and send every reader somewhere else. Its `[book] authors` and copyright header are the template author's, too.

The same instructions are repeated as the first section of `AGENTS.md`, since that is the file an agent loads and this one is not. It ends with "delete this whole section once it is done" - keep the two in step if you edit either.

Then work through the `<!-- ... -->` comments in `AGENTS.md`. They mark the five places that cannot be templated: what the project is, the layout table, the build and test commands, which test invocations are forbidden, and the quirks worth knowing.

## Sections to delete when they do not apply

Several rules in `AGENTS.md` assume a Cobra CLI with a viper config file, a mage task runner, and an mdBook under `docs/`. Delete the ones a given repository does not have rather than leaving a rule pointing at machinery that is not there - a rules file whose claims do not check out is the failure mode its own "A doc's claims are checkable" section exists to prevent.

- **A new flag belongs in the config file too** - needs cobra plus viper plus a sample config at the root.
- **Help text has a budget** - needs a CLI with `--help` as a real surface.
- **Exit codes are a compatibility promise** - needs `cmd/exit.go` and a published table.
- **Pinned versions, digests, and checksums** - needs `magefiles/AGENTS.md` to exist, since it defers to it.
- **Generated output goes in a pull request of its own** - needs a generator that renders committed artifacts.
- **`docs/AGENTS.md`** in full - needs mdBook and release-please.

## Labels the issue forms need

GitHub drops an unknown label silently, so create these before the forms are used:

```sh
gh label create bug --repo OWNER/REPO
gh label create ci --repo OWNER/REPO
gh label create test-flake --repo OWNER/REPO
gh label create enhancement --repo OWNER/REPO
gh label create documentation --repo OWNER/REPO
gh label create stack-task --repo OWNER/REPO
```

Renovate names five more in `renovate.json5`, and drops them just as silently:

```sh
gh label create renovate/oci --repo OWNER/REPO
gh label create renovate/action --repo OWNER/REPO
gh label create renovate/gomod --repo OWNER/REPO
gh label create renovate/pre-commit --repo OWNER/REPO
gh label create security --repo OWNER/REPO
```
