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

`go run ./magefiles/core dev:checkTables` enforces this, and prints the expected form of every line it rejects - so the fix is its own output rather than counting spaces by hand. **It only sees tracked files**, because it walks `git ls-files` rather than the filesystem: a brand-new page is invisible to it, and to the pre-commit hook that runs it, until the file has been `git add`ed. A clean run on a change that adds a document means nothing until then.

## A moved page keeps a redirect

Every page under `docs/` is published at `https://OWNER.github.io/REPO/<path>.html`, so renaming or moving one breaks whatever already links to it from outside this repository -- a bookmark, an older release's manifest, a third-party listing. mdBook does not notice: the old URL simply 404s.

So a move comes with an entry in `[output.html.redirect]` in [`book.toml`](../book.toml), keyed on the old path and valued with the new one. The value is resolved relative to the key's own directory, not the site root, so a page that moved one level down reads `"/guides/install.html" = "setup/install.html"`. mdBook writes a small redirecting HTML file at each key when the book is built.

Keep the entry for as long as anything outside the repository might still hold the old link. They cost one file each, and `mdbook build` is what proves one works -- the redirect file appears at the old path and names the new one.

A page that was never published does not need one.

## A pinned release version is release-please's to update

A page that pins this project's own version -- an image tag, a release archive's filename, a version in an install snippet -- goes stale on the next release, and nobody notices until a reader copies the command. Release-please will bump a Markdown page the same way it bumps a manifest, once the page tells it where to look and the configuration tells it the page exists. Both halves are required: an annotated page that is not listed is never read, and a listed page with no annotation is read and left alone.

Annotate with the block form, as HTML comments on their own lines around the fence:

````markdown
<!-- x-release-please-start-version -->
```sh
podman pull ghcr.io/OWNER/REPO:0.1.0
```
<!-- x-release-please-end-version -->
````

The per-line form, `x-release-please-version` on the same line as the version, is the one to reach for in YAML or a Dockerfile and the wrong one here: inside a fenced code block an HTML comment is content, so it renders to the reader as part of the command they are about to run.

Then add the page to `extra-files` in [`release-please-config.json`](../release-please-config.json). The next release PR rewrites every version inside every annotated region in it, and a page that nothing lists is a page release-please does not open.

**An annotated region must hold no other dotted-decimal string.** Release-please replaces each semver match it finds between the markers, and `0.0.0.0` matches one -- a region containing both an image tag and that address comes out of the next release with `https://0.1.1.0:41609` in it. `127.0.0.1`, a CIDR, and any third-party version pin are the same hazard. Where an example needs both, split the fence: the version-free half stays unannotated.

A version in a generated page is annotated where it is written, not where it lands: a pinned version in `README.md` gets the markers in `README.md`, since `docs/index.md` is rewritten from it on the next generate run.

Illustrative output is the one place to prefer a placeholder over an annotation. Where a command's own output carries a version with a suffix that reads as a semver prerelease - `rpm -q` printing `<project>-0.1.0-1.x86_64` - a rewrite of that region can swallow the suffix. Write `<project>-<version>-1.x86_64` instead, which matches the `:<tag>` placeholders elsewhere and cannot go stale.

Set the version to the current release when adding the annotation, rather than leaving whatever was there. The markers fix the page from the next release onward and do nothing for the one in front of a reader today.

## Do not hand-edit the generated pages

Every generated path below is also marked `linguist-generated` in [`.gitattributes`](../.gitattributes); see the "Keeping `.gitattributes` in sync with the generators" section of the root [`AGENTS.md`](../AGENTS.md) when a generator starts or stops writing one.

Parts of the `docs/` tree are generated and are overwritten wholesale on the next run, so an edit made there does not survive the next commit.

| Path               | Generated from                                                                     |
| ------------------ | ---------------------------------------------------------------------------------- |
| `docs/index.md`    | `README.md`, with its relative links rebased from the repository root onto `docs/` |
| `docs/security.md` | `.github/SECURITY.md`, with its relative links rebased the same way                |
| `docs/SUMMARY.md`  | The mdBook table of contents, compiled from the rest of the `docs/` tree           |

<!-- TEMPLATE: two more rows belong here once the generators exist - `docs/commands/` from the Cobra command tree rendered by `doc.GenMarkdownTreeCustom`, and `docs/schema/` from the struct reflection a schema generator walks. Both are deleted and recreated wholesale, unlike the three above, which are rewritten in place. Add the row and the `.gitattributes` pattern in the same change as the generator. -->

The three above are rewritten in place on every run, and `docs/index.md` and `docs/security.md` are copies of files that live outside `docs/` because GitHub reads them where they are.

To change one of those pages, change what it is generated from - a command's `Short`/`Long`/flag help, a godoc comment, or the struct and tags the schema generator reflects - and then regenerate:

```sh
go run ./magefiles/core generate:document
```

Use that `go run` form rather than `mage generate:document`. It runs the same entry point and needs only the Go toolchain, so it works on a host that has no `mage` binary installed.

Commit the regenerated files with the change that caused them, so the tree matches the code.
