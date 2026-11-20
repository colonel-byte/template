# Working in this repository

<!-- One paragraph: what this project is, who runs it, and what it is not. Link README.md for the full pitch, and the two or three docs that explain the model an agent has to hold to make a correct change. -->

<!-- One paragraph on the operating assumptions that are not obvious from the code - where it runs, what it may and may not reach, what is staged ahead of time. -->

## First, if this repository was just created from the template

This file came from [`colonel-byte/template`](https://github.com/colonel-byte/template) and still carries its placeholders. Replace them before doing anything else, and treat it as one change rather than a thing to come back to - a placeholder that survives the first commit is a wrong URL nothing will ever fail on.

```sh
grep -rn 'OWNER/REPO\|OWNER\.github\.io\|<project>' . --exclude-dir=.git --exclude-dir=vendor
```

- **`OWNER/REPO`** - the new repository's slug. It appears in the issue-tracker section below, in [`.github/AGENTS.md`](.github/AGENTS.md) (the Scorecard viewer link, the advisory URL, the `gh label list` snippet), in [`.github/ISSUE_TEMPLATE/config.yml`](.github/ISSUE_TEMPLATE/config.yml), in [`.github/ISSUE_TEMPLATE/01-failing-build.yml`](.github/ISSUE_TEMPLATE/01-failing-build.yml), in the `docs/workflows/` pages, in [`.goreleaser.yaml`](.goreleaser.yaml) (the `release.github` block, the `-ldflags -X` module path, and the `cosign.pub` link in the footer), in **the three badges at the top of [`README.md`](README.md)**, and in **[`book.toml`](book.toml)**.
- **`OWNER.github.io`** - the published book's host, in `docs/AGENTS.md` and the issue chooser's documentation link.
- **`<project>`** - the binary and module name, in the layout table below, the help-text examples, and the reproduction steps in [`.github/ISSUE_TEMPLATE/03-bug.yml`](.github/ISSUE_TEMPLATE/03-bug.yml).

**The badges and [`book.toml`](book.toml) are the two worth checking by eye**, because they are where a wrong value produces something that builds, renders, and looks right. A badge left on the template's slug reports another repository's release and security posture as this one's, at the top of the front page; `README.md` says what each badge needs before it resolves at all.

As for `book.toml`: `git-repository-url` is the repository icon in the book's header and `edit-url-template` is the pencil beside every page; pointed at the template's own slug, both silently send a reader to the wrong repository. Set them both to the new slug, and change `main` in each if the default branch is named something else:

```toml
git-repository-url   = "https://github.com/OWNER/REPO/tree/main"
edit-url-template    = "https://github.com/OWNER/REPO/edit/main/{path}"
```

Also in `book.toml`: `[book] authors` and the copyright header, which are the template author's, not necessarily this project's.

**The placeholders are only half of it.** Generating a repository from a template copies files, not settings, so several of the files that just arrived are inert: `.github/CODEOWNERS` requires nothing until a ruleset says so, the workflows have no secrets, the labels the issue forms name do not exist, and the book 404s until Pages is enabled. [`docs/dev/new-repo-checklist.md`](docs/dev/new-repo-checklist.md) is the list, with the `gh` commands for the parts that can be scripted and the order to do the first run in.

Then work through the `<!-- ... -->` comments in this file. They mark the five things that cannot be templated - what the project is, the layout table, the build and test commands, which test invocations are forbidden, and the quirks worth knowing - and a rule left as a comment is a rule nobody follows. Delete the sections whose machinery this repository does not have, rather than leaving them pointing at files that are not there; [`README.md`](README.md) lists which sections depend on what.

Delete this whole section once it is done.

## Layout

<!-- One row per directory an agent will touch. Say what lives there, not what it is named. -->

| Path                     | What lives there                                                                                           |
| ------------------------ | ---------------------------------------------------------------------------------------------------------- |
| `cmd/`                   | Cobra commands, and `cmd/<project>`, the released binary's `main`                                          |
| `pkg/`                   | Packages usable from anywhere in the module                                                                |
| `internal/`              | Go packages private to this module                                                                         |
| `test/`                  | The shared e2e harness and its suites                                                                      |
| `magefiles/`             | The mage task runner's source - build, test, and generate targets                                          |
| `docs/`                  | The mdBook source. Generated subpaths must not be hand-edited                                              |
| `docs/agent/choice-*.md` | Records of non-obvious decisions and the tradeoffs behind them - read before reversing one                 |

Several directories carry their own `AGENTS.md` with rules specific to that directory. Check the closest one before editing. Keep this table in sync when adding or removing one -- add a row for a new `AGENTS.md`, and remove its row when the file goes away.

| `AGENTS.md`                                  | Covers                                                        |
| -------------------------------------------- | ------------------------------------------------------------- |
| [`docs/AGENTS.md`](docs/AGENTS.md)           | Markdown formatting, generated pages                          |
| [`.github/AGENTS.md`](.github/AGENTS.md)     | Scorecard-safe workflow changes, PR description format        |
| [`magefiles/AGENTS.md`](magefiles/AGENTS.md) | Target layout, the vendored offline build, pin ownership      |

`CLAUDE.md` in the repository root is a symlink to this file, so a Claude Code session started there loads these rules without being asked. That covers the root and nothing else: the directory-specific files in the table above are not symlinked, so a session whose working directory is one of them starts with none of their rules in context. The root is where a session normally starts, so this is usually moot - but do not treat "nothing was loaded for me" as "there are no rules here", and do not grep a rules file for the one thing you thought to ask about. Open the closest `AGENTS.md` and read it before writing anything in that directory.

## Building and testing

Build and test targets run through mage, not a Makefile. `magefiles/` itself does not build with a plain `go build` - it needs mage's special build tag - so drive it through the `mage` CLI or `go run ./magefiles/core <namespace>:<target>`. The latter form needs only the Go toolchain and works on a host with no `mage` binary installed; it's what CI and the pre-commit hooks use.

```sh
go run ./magefiles/core -l                  # list every target
go run ./magefiles/core generate:document   # regenerate docs/index.md, docs/security.md, docs/SUMMARY.md
go run ./magefiles/core dev:checkTables     # check Markdown table padding
go run ./magefiles/core dev:vendor          # go mod tidy, then go mod vendor
go run ./magefiles/core dev:verifyVendor    # go mod verify, then a GOPROXY=off build
```

Plain Go commands work for anything mage doesn't wrap - `go build ./...`, `go vet ./...` - but exclude `./magefiles/...` from those, since it fails a bare build for the reason above.

<!-- TEMPLATE: name the test invocations that are forbidden and why. The usual one: never run `go test ./...` without `-short`, because it pulls in the suites that need real infrastructure. -->

### The build is vendored and offline, and that is a requirement

`vendor/` is committed, and the build has to work on a machine that cannot reach a module proxy. The Go toolchain picks `-mod=vendor` automatically when `vendor/` is present, which makes this invisible until it isn't: **a dependency added without re-vendoring builds fine locally and fails only where it matters.**

So a change to the dependency graph is not finished until `vendor/` is rewritten and committed in the same change. Three things enforce it, and each catches the same mistake at a different distance from the author:

- **`go run ./magefiles/core dev:vendor`** rewrites it - `go mod tidy` then `go mod vendor`.
- **A `verify-vendor` pre-commit hook** runs `dev:verifyVendor` whenever `go.mod`, `go.sum`, or `vendor/` changes.
- **[`.github/workflows/check-go-mod.yaml`](.github/workflows/check-go-mod.yaml)** re-runs tidy and vendor and fails on any diff, then runs `dev:verifyVendor` - which builds with `GOFLAGS=-mod=vendor GOPROXY=off`, the invocation a release has to survive.

Two consequences. **A new dependency is a bigger decision than it looks**, because it is carried by every clone from then on; the cost of a dependency is what it links, not what you call, so check `go list -deps ./...` before adding an import. And **`vendor/` is generated** - marked `linguist-generated` in [`.gitattributes`](.gitattributes), and excluded from every content-reading pre-commit hook, because a hook that reformats a vendored file makes the next vendor run look like a change nobody made.

## Keeping `.gitattributes` in sync with the generators

Everything a generator writes is marked `linguist-generated` in [`.gitattributes`](.gitattributes), which is what collapses it in a GitHub diff and keeps it out of the repository's language statistics. A generated file that is not marked shows up as hand-written work in every review, and a thousand-line regenerated page buries the twenty lines that were actually written.

So a change to what a generator produces is a change to `.gitattributes` as well. That means all three of:

- **a new generator, or a new output path** - add the path. A generator that writes three files beside each input needs all three marked; a pattern covering only the first leaves the others looking hand-written.
- **a generator that writes into a new subdirectory** - widen the pattern rather than adding a sibling entry, where the directory depth is the only difference.
- **a generator removed, or its output moved** - delete the entry. A pattern matching nothing is not an error anybody sees.

Patterns are grouped by the mage target that writes them, with the target named in a comment, so the file reads as a list of what each generator owns.

Hand-written files inside a generated tree stay unmarked. Mark the pattern the generator writes, not the directory it writes into.

Check a specific file rather than reasoning about the patterns:

```sh
git check-attr linguist-generated -- docs/commands/example.md
```

## Keeping `README.md` current

[`README.md`](README.md) is the repository's front page and the pitch: what the project does, what it supports, how to install it, and one example per entry point. It is also the source for the book's landing page, so a stale README is stale in two places.

Treat it as part of the change, not a follow-up:

- **a feature lands** - add it where a reader would look for it. A new command belongs in `Usage Examples` with a runnable snippet; a new install channel in `Installation`; a new config file or schema in `Configuration and Schemas`.
- **a feature is removed or renamed** - delete or rename every mention, including the example that exercised it. A documented flag that no longer exists is worse than no documentation.
- **behaviour changes under an existing heading** - fix the prose and the snippet together. Commands shown in `README.md` are copied and pasted; they should work as written against the current build.

The generated reference under `docs/` comes from the code, so it corrects itself on the generate target. `README.md` does not - it is hand-written, and nothing fails if it drifts. After editing it, run that target so the landing page matches.

## Say so when you notice user-facing drift, even out of scope

Documentation and user-facing strings drift silently. Nothing fails, no test goes red, and the generated pages under `docs/` regenerate the stale text faithfully - so a string that describes a capability the code removed can survive several releases, looking authoritative the whole time.

So drift noticed is drift reported. When work in one corner of the tree turns up a stale user-facing surface somewhere else - a flag description, a `Long` or example string, a `jsonschema` example, a module doc block, a README snippet, a key missing from the sample config - name it in the change description, with the file and what is wrong. Do not quietly widen the change to fix it, and do not stay silent because it is out of scope: the operator deciding whether it is worth a follow-up is the point, and an unfixed drift they know about is a fine state to leave the tree in.

Two things make this worth a rule rather than a habit:

- **the stale copy is usually upstream of several generated files**, so what looks like one wrong line in `docs/` is one wrong line in a Go string plus four regenerated pages, and a reader who corrects the page by hand loses the fix on the next generate run. Chase the wording back to the string that produced it before reporting where it lives.
- **a vendored string can be the source.** Where a flag description comes from a vendored upstream package, `vendor/` cannot be edited - the fix is a constant in this repository's own strings package and a changed reference.

## A doc's claims are checkable, so check them

The section above is about wording that has gone stale. This one is about the harder case: a document that is confidently, verifiably wrong, because the thing it describes no longer exists. A dead reference reads exactly like a live one. `mdbook build` does not fail on a dead relative link, `go vet` never opens a Markdown file, and no test asserts that a path, symbol, flag, or module named in prose is still in the tree - so these survive indefinitely, and the longer they survive the more authoritative they look.

The cost is not spread evenly. A wrong count in an internal design note costs a reader a minute of confusion. A command printed in `docs/guides/` that cannot run costs an operator an afternoon, because they will assume they typed it wrong before they assume the guide is wrong. Weight the checking accordingly.

So when a change touches something a document describes, check what that document says about it. Each of these is one command:

| Failure mode                                   | How to check it                                                                                                                    |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| A path, file, or directory that moved or went  | `ls`/`find` the path as written                                                                                                    |
| A Go symbol that was renamed or unexported     | `grep -rn --include='*.go' '<name>' .`, or `go doc`. An unexported-to-exported rename is the easy one to miss                      |
| A module asserted to be present or absent      | `grep -n '<module>' go.mod vendor/modules.txt`, and `go list -deps ./cmd/<project>` for whether it is linked                       |
| A mage target that does not exist              | `go run ./magefiles/core -l`                                                                                                       |
| A command or flag that does not exist          | `grep -rn '<flag>' cmd/`. Check the command's `Use:` line, not your memory of the command tree                                     |
| A config key with no backing field             | It must be in the config struct, the viper defaults, and the sample config file                                                    |
| A quoted error or log string                   | `grep -rn '<fragment>'` for the literal. Quoting these from memory is how a guide ends up printing an error the code does not emit |
| A relative link between pages                  | `ls` the target. Nothing in the build catches a dead one                                                                           |
| A count, version, or pin given as current fact | Read the thing being counted. A number that was measured once and never re-measured is the default state of every count in a tree  |

### Flag it, and say whether it should be removed rather than fixed

Report what you find the way the section above says: name the file, the line, and what is actually true, and leave the decision with the operator rather than quietly widening the change. The one thing to add is a recommendation on **removal**, because that is the call a reader cannot make for themselves and the one that is easy to get wrong in both directions.

- **A dead reference inside a live document** is a fix: correct it, or delete the sentence if what it pointed at is gone and nothing replaced it. An orphaned heading with no body, left behind when its content was removed, is the same case.
- **A document whose entire subject was removed** is not a deletion. `docs/agent/choice-*.md` are constraints rather than history, and the reasoning in one usually still explains something that is still here - a Dependabot ignore entry, a Scorecard annotation, the shape of a package. Mark it superseded, name the change that superseded it, and leave the argument in place. Deleting it loses the answer to "why is this odd thing still configured".
- **A configuration entry a removed document justified** is a separate question with its own blast radius, and it is not yours to fold into a documentation change. A stale Scorecard annotation or Dependabot ignore is doing nothing, but removing it changes what CI measures. Flag it and let the operator decide.

When a change removes a module, moves a package out of the tree, or renames a type, the documentation sweep is part of it - and `grep -rn '<old-name>' --include='*.md' .` is the whole technique.

## A new flag belongs in the config file too

A flag on a cobra command is one half of how this project is configured. The other half is the config struct that `v.Unmarshal(&resolvedConfig)` decodes the config file, the environment and the viper defaults into - and `resolvedConfig` is what a flag registration reads for its default value. A flag wired only into cobra works, which is why it is easy to stop there, and it leaves an operator restating on every invocation something their whole fleet has one answer to.

So adding a flag that carries a setting rather than a one-off instruction means all of:

- **a field in the right section of the config struct**, with the `json` and `mapstructure` tags that name its key. The doc comment is the generated schema's description, so write it for an operator authoring YAML, and say which commands honour it when not all of them do.
- **a `v.SetDefault(...)` beside the other defaults**, even when the default is the Go zero value. `Unmarshal` builds its result from `v.AllKeys()`, so a key with no default registered is invisible to it and silently loses environment-variable support.
- **the flag's default read from `resolvedConfig`**, not a literal. This is the line that makes the config file take effect, and the ordering falls out of it for free: flag beats config beats default.
- **home-path expansion, for a path-valued field.** A shell expands a leading `~` in a flag value before the binary sees it; nothing expands one in a config file or an environment variable, so `~/.cache/...` reaches `os.Open` as a literal and fails naming a path a shell would have found. Add the field to the expansion helper and a case to its test, which compares the whole struct and so fails until the expectation is written out.
- **a line in the sample config at the repository root.** It is what an operator copies and what the `README.md` examples point `--config` at, so a key missing from it is a key nobody discovers. Add it commented out, with its one-line explanation, where the setting should not be on by default - a sample that writes files or changes host state on a first run is worse than one that is incomplete. A config under `test/` is a fixture rather than a sample: change it only when a suite needs the key.
- **a `README.md` snippet showing the config-file form**, per the section above, and a regenerated schema.

Not every flag qualifies. `--confirm` and `--dry-run` are decisions about one invocation, and a config file that answered them standingly would be a hazard rather than a convenience; they stay cobra-only. The test is whether an operator would give the same answer every time they run the command.

**Where the key goes is a question for the person, not a guess.** The section a field lands in is the key an operator types, it ships in a generated schema, and moving it later breaks every config file that set it - while from inside the code two or three placements usually look equally defensible. Say which one you think is right and why, show the YAML it produces, and ask before writing it.

## Help text has a budget, and flag help has a hard cap

`--help` is the whole interface. This tool runs from a terminal, often on a host with no UI and no hosted docs to tab over to - so what `--help` prints is what an operator has. It is also the one surface nothing tests: a 600-character flag description compiles, generates clean documentation, and passes every hook.

### Flag help: 125 characters, hard cap

One sentence. It has to fit the hanging-indent column of an 80-column terminal, which is about a line and a half. Over the cap, the detail goes into a guide under `docs/guides/` and the flag string names it:

```go
// Good: one sentence, then where the rest lives.
CmdCreateFlagContentStore = "`PATH` of a content store to resolve images and files from instead of the network. See docs/guides/content-store.md."
```

**The trap is that the cap is invisible from where the string is written.** The same constant feeds `--help` and the generated page under `docs/commands/`, and the page renders 600 characters as a tidy paragraph. So check the terminal, not the page:

```sh
COLUMNS=80 go run . create --help
```

Two further things about the measurement:

- **Cobra appends `(default ...)` to the rendered line**, so a string at 120 characters can still produce a 130-character line. Where the default is long or machine-specific, leave room for it.
- **The constant name is not a reliable filter.** Flag help with no `Flag` in the name is missed entirely by a grep for `Flag\w* = "`. Measure the rendered `--help`, which cannot be fooled.

### `Long` and `Example`: a budget, not a number

These are deliberately not capped. A `Long` has real work to do, and trimming good prose to fit an arbitrary count loses information. What matters is whether the whole of `--help` can be read without scrolling, because what gets scrolled past does not get read.

- **Aim for the whole `--help` under about 45 lines.** Check it after changing a `Long`, an `Example`, or a flag string - the number to watch is the total, not any one part.
- **Three to five examples.** Pick the common cases; an example per flag combination belongs in a guide. The house form is two lines each, a `#` comment then a runnable `$ <project> ...` line.
- **`Long` says what the command does and what an operator must know before running it.** A walkthrough, a worked scenario, or an enumeration of how several flags interact is a guide.

### None of this is enforced in CI, on purpose

A gate would be raised the first time it blocked something urgent, and the useful behaviour is a person noticing. Measure it when you touch it:

```sh
go run . <command> --help | sed -n 's/^  *-\{1,2\}[^ ]*.*  \{2,\}\(.*\)$/\1/p' | awk 'length > 125'
```

## Exit codes are a compatibility promise

`cmd/exit.go` maps an error to a process exit code, and `README.md` publishes the table. This tool is driven from scripts, and a code that means "nothing was touched and a retry is safe" is a promise a wrapper relies on - so the numbers are an interface, not an implementation detail.

Adding a classification means adding a code. Do not repurpose or renumber an existing one, and update the mapping, the `README.md` table and the test together. A new error that nobody classified falls through to `1`, which is the right default: a caller sees "something failed" rather than a wrong promise.

## Pinned versions, digests, and checksums

Pins live all over the tree - container image digests in Dockerfiles, package versions, upstream tags in a pins file, checksums in Go constants - but the machinery that moves them is mage and Renovate. Every pin needs one of those owning it, or a written bump procedure where neither fits. Before adding a pinned value anywhere, read the rule and the worked examples in [`magefiles/AGENTS.md`](magefiles/AGENTS.md#every-pin-needs-an-owner).

## Committing changes

Run the hooks over the whole tree before every commit and fix anything they flag. Prefer [`prek`](https://github.com/j178/prek) - `prek run --all-files` - if it is on `PATH`; it reads the same [`.pre-commit-config.yaml`](.pre-commit-config.yaml) and runs the same hooks, just faster. Fall back to `pre-commit run --all-files` when it isn't. Several of the hooks rewrite files (`end-of-file-fixer`, `trailing-whitespace`, `keep-sorted`, `addlicense`, doc/schema generators) - re-stage after the run. Don't skip hooks with `--no-verify` or `SKIP=`.

Two gaps worth knowing about, because a green hook run is weaker evidence than it looks:

- **A lint hook that reads the staged diff only sees uncommitted changes**, so a clean `prek run --all-files` does not mean CI's lint job passes on code that is already committed.
- **`--all-files` skips untracked files.** A new file is invisible to every hook until it is `git add`ed.

## Go test formatting

Do not condense struct literals in table-driven test cases onto single lines (such as `{name: "...", input: "...", want: "..."}`). Format each field on its own line within the struct literal to keep test cases legible and easy to review, diff, and edit.

```go
// Do not condense:
{name: "unconstrained release", release: "", want: true},

// Prefer multi-line:
{
	name:    "unconstrained release",
	release: "",
	want:    true,
},
```

## Writing style

Use `-` (hyphen), not `—` (em dash), wherever a hyphen reads fine.

## Commit messages

Keep commit messages to a short title only - no body paragraph explaining the change. End with a `Co-Authored-By` trailer.

```
test(fuzz): add fuzz targets for cfg, clustercfg, and identify source

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
```

## Pull requests

Before creating or updating a pull request description, read [`.github/pull_request_template.md`](.github/pull_request_template.md) and follow the format rules in [`.github/AGENTS.md`](.github/AGENTS.md#writing-a-pull-request-description).

### Generated output goes in a pull request of its own

A change that adds or re-renders generated artifacts does not share a pull request with the code change that prompted it. Where both are needed, they are two layers of a stacked PR: the generator, template, or format change first, then the rendered output on top of it.

This is about what a reviewer can see. A generator run that touches an existing template re-renders every directory already on disk - so the twenty lines that were actually written end up buried in a diff nobody can read to the end. GitHub collapses `linguist-generated` files, which helps with the file list and not at all with the review: the reasoning that needs checking is in the template and the generator, and it is the thing that disappears.

So:

- **A generator or template change, and the output it renders, are separate layers.** The lower layer carries whatever was reasoned about. The upper layer carries the rendered output and nothing else.
- **A render that is only a re-render is its own pull request.** Bumping a pin, backfilling a release line, refreshing a cached asset: there is no code to review, so there is nothing for it to be mixed with.
- **Check the split before pushing, not after.** `git show --stat <branch> -- <generated-path>` on the lower layer should print nothing. It is easy to leave rendered output in the wrong layer by running the generator while the wrong branch is checked out, and the fix once a reviewer is reading is a rebase of everything above it.

The one thing that does belong with the code is a *fixture* a test reads - those live under `test/`, and the rule does not reach them.

### Naming the branches in a stacked PR

Every branch in a stack is named `stack/<n>-<2-3 word header>`:

```
stack/1-color-tty
stack/2-strict-config
stack/3-verify-flags
```

**`<n>` counts within the stack, from 1.** It is the layer's position, bottom to top, so a nine-layer stack runs `stack/1-...` to `stack/9-...` and the next stack starts again at `stack/1-...`. A stack is self-contained; the number says where a branch sits in its own chain and nothing about any other work. Do not carry on from the highest number seen in the repository's history - that couples unrelated stacks to each other, and it means a branch name cannot be read without knowing the whole branch history.

**`<2-3 word header>`** is the layer's concern in kebab case, not a sentence: `exit-codes`, `flag-shorthands`, `phase-progress`. Two or three words, matching the one-sentence description the layer should already have - a layer whose concern will not fit in three words is usually two layers. Do not restate the stack's topic in every branch; the chain supplies that.

Because the number restarts, two live stacks can both hold a `stack/1-`, and uniqueness rests entirely on the header - `stack/1-role-parameter-gating` and `stack/1-color-tty` coexist happily. Pick a header that names this layer's concern and the clash does not arise; `gh stack add` fails outright on a branch that already exists, so a collision is loud rather than silent.

Names are used verbatim by `gh stack add`, so there is nothing to correct afterwards. Renaming a branch once its PR is open is not worth it: `gh stack` tracks the chain by branch, so the fix is to delete the stack and re-open it, which discards any review already left on it.

## Tracking work in issues

Work here is tracked in issues on the main repository, [`OWNER/REPO`](https://github.com/OWNER/REPO). Before starting a task, check for a related issue - someone may already have filed it, or filed something that overlaps enough to change the approach:

```sh
gh issue list --repo OWNER/REPO --search "<keywords>" --state all
```

If there is no issue covering the work, ask the person whether they would like one submitted before continuing. Don't open one unprompted, and don't treat a missing issue as a reason to stop - they may well say no.

When work is being broken up into a stacked PR, the issue structure should mirror the stack: the parent issue describes the whole change, and each layer of the stack gets a subissue created off that parent. Link each PR in the stack to its own subissue, not to the parent.

Only the top-most PR in the stack closes the parent issue - each lower PR closes its own subissue and nothing more, so the parent stays open until the whole change has landed. Don't move that closing reference to an earlier PR unless the operator asks for it.

## Quirks worth knowing

<!-- The things that have cost somebody an afternoon and are invisible from the code. Delete the ones that do not apply here; the three below apply to most Go repositories in this org. -->

### Relative-path fixture constants don't move themselves

Test files and YAML fixtures hand-encode `../` depth to reach the repository root or a shared `testdata/` directory. Moving either endpoint of one of these paths changes the required `../` count by exactly the number of levels moved, and nothing catches a wrong count except actually running the test - `go vet` only compiles, it doesn't execute. Recompute depth from the final location of both the file and its target, not from an intermediate state, if the move happens in more than one step.

### A `test/` path is not hidden from Scorecard

OpenSSF Scorecard's `isTestdataFile` logic excludes anything under a `src/test/` or `testdata/` prefix from every check, by path convention. A plain `test/` prefix gets no such exclusion - so don't rely on a `test/`-adjacent path to keep something out of a check. Read `isTestdataFile` in `checks/fileparser/listing.go` upstream if it matters.

### `docs/agent/choice-*.md` are constraints, not history

When a decision looks reversible from the code alone, check `docs/agent/` first. These are usually the record of a tradeoff already made on purpose, not an oversight waiting to be cleaned up.

# Markdown line formatting

## Do not hard-wrap lines in Markdown documents

Markdown files should NOT have hard-wrapped or truncated lines. Instead, write full paragraphs and list items on single, continuous lines without manual line breaks, letting editors and renderers handle line wrapping dynamically.

- Write entire paragraphs or list items on a single line without manual line breaks / hard wraps.
- Let editors and rendering tools handle line wrapping dynamically.

## Pad table cells so the columns line up

Tables are written in their aligned form, so that the source reads as a table and a reviewer can see the columns without rendering the file. Pad every cell with spaces to the width of the widest cell in its column, keep one space of padding inside each pipe, and run the dashes of the delimiter row out to that same width. Every line of the table then has the same length.

Write this:

```markdown
| Path                                   | What it holds     |
| -------------------------------------- | ----------------- |
| `.spec.config.registries[N].auth.user` | Registry username |
```

Not this:

```markdown
| Path | What it holds |
| --- | --- |
| `.spec.config.registries[N].auth.user` | Registry username |
```

Adding a row that is wider than the column repads the whole table, which is expected - the alignment is part of the content, not a one-time cleanup.
