# GitHub Pages

[`.github/workflows/deploy-book.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/deploy-book.yaml) builds the mdBook from `docs/` on every push to `main` and pushes the result to the `gh-pages` branch, which is what serves `https://OWNER.github.io/REPO/`.

It needs no secrets: `GITHUB_TOKEN` is enough to push to `gh-pages`, and `contents: write` is granted on the job rather than at the top of the file.

## It publishes, it does not regenerate

`mdbook build` renders whatever `docs/` holds at that commit. Three of those pages are generator output - `docs/index.md`, `docs/security.md` and `docs/SUMMARY.md` - and nothing here runs the generator, so a page that was not regenerated before merging is published stale rather than rebuilt.

That is deliberate: a deploy that regenerates would publish content that is not in the repository, so the committed tree and the published site would disagree. The `docs-generator` pre-commit hook is what keeps them in step, and [`check-go-mod.yaml`](check-go-mod.md) is the model - catch it before the merge, not after.

## Pages has to be enabled, and that is not in the repository

This is the one part the template cannot carry. The workflow pushes to `gh-pages` whether or not Pages is configured, so it goes green while the site 404s. In the repository's Settings, under Pages, set the source to **Deploy from a branch**, branch `gh-pages`, folder `/`. See [`new-repo-checklist.md`](../dev/new-repo-checklist.md) for the rest of the settings that do not travel from a template.

Until that is done, every `https://OWNER.github.io/REPO/` link - in `README.md`, the issue chooser, and `docs/AGENTS.md`'s redirect rules - points at nothing.

## The redirect rule matters here

Because the site is published, renaming or moving a page breaks whatever already links to it from outside the repository, and mdBook does not notice - the old URL simply 404s. A move comes with an `[output.html.redirect]` entry in [`book.toml`](https://github.com/OWNER/REPO/blob/main/book.toml). [`docs/AGENTS.md`](../AGENTS.md) has the exact form, including why the value is resolved relative to the key's own directory rather than the site root.

## Building it locally

```sh
mdbook build && mdbook serve
```

`mdbook serve` watches `docs/` and reloads, which is the quick way to check a redirect actually resolves - the redirect file appears at the old path and names the new one.
