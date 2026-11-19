# Release Please

[`.github/workflows/release-please.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/release-please.yaml) runs on every push to `main` and keeps one open pull request - the release PR - with the next version number, the regenerated `CHANGELOG.md`, and the bumped version in `.release-please-manifest.json`. Merging that pull request is what cuts a release: release-please then pushes the `vN.N.N` tag, which is what [`release.yaml`](release.md) fires on.

Nothing here decides the version. Conventional commit types in the pull request titles do, by way of the `changelog-sections` in [`release-please-config.json`](https://github.com/OWNER/REPO/blob/main/release-please-config.json).

## Why the app token

The tag is the whole reason. A tag pushed with the default `GITHUB_TOKEN` does not trigger another workflow, by design, so `release.yaml` would never fire and a release would be cut with no artifacts attached to it. An app installation token has no such restriction.

That is also why the top-level `permissions:` is `contents: read` and nothing else: every write is done with the app token, so `GITHUB_TOKEN` is only along for the ride.

| Secret                             | What it is            |
| ---------------------------------- | --------------------- |
| `RELEASE_PLEASE_GITHUB_APP_ID`     | The app's client ID   |
| `RELEASE_PLEASE_GITHUB_APP_SECRET` | The app's private key |

The same credentials drive [`renovate.yaml`](renovate.md), so there is one installation to grant and one key to rotate.

## The two config files

| File                            | What it holds                                                                                  |
| ------------------------------- | ---------------------------------------------------------------------------------------------- |
| `release-please-config.json`    | How to version and what to write - release type, tag format, changelog sections, `extra-files` |
| `.release-please-manifest.json` | The current version. Release-please owns this file; do not edit it by hand                     |

`extra-files` is the one to fill in. Any file outside the manifest that pins this project's own version - an install snippet, a collection manifest, an image tag in a guide - has to be listed there *and* annotated in place, or release-please does not touch it. Both halves are required; see the "A pinned release version is release-please's to update" section of [`docs/AGENTS.md`](../AGENTS.md).

## `CHANGELOG.md` is not yours to edit

Release-please rewrites it wholesale on every run. Editing it by hand - including to resolve a merge conflict - is lost on the next release and shows up as a diff nobody authored. Let the bot win the conflict.
