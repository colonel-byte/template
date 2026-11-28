# Renovate

[`.github/workflows/renovate.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/renovate.yaml) runs Renovate against this repository and nothing else, every four hours and on demand. It replaces the hosted `renovate[bot]` GitHub App: the bot is not installed, and the config in [`.github/renovate.json5`](https://github.com/OWNER/REPO/blob/main/.github/renovate.json5) is read by this run rather than by Mend's servers.

## Why it runs self-hosted

The hosted app is simpler and needs no workflow. What it cannot do is use this organization's own GitHub App credentials, and that matters for two reasons:

- **Pull requests opened with an app installation token trigger other workflows.** A pull request opened with the default `GITHUB_TOKEN` does not, by design, so a dependency branch would sit with no checks on it - and an update nothing has tested is not something a reviewer can act on.
- **One app, one audit trail.** The same app that release-please uses opens these, so there is one installation to grant, one private key to rotate, and one bot in the repository's audit log.

The cost is that the workflow is now a thing to maintain, and that the Renovate version is pinned to the action's SHA rather than tracking the hosted service.

## The run is cached

Renovate resolves every dependency from the registry on each run, which once a repository has real dependencies is most of the wall clock. The workflow restores Renovate's package cache - resolved versions and changelogs - before the Renovate step, so a scheduled run only re-asks the registries about what moved. The cache key carries `github.run_id`, so every run writes a new entry and `restore-keys` picks up the newest previous one; a fixed key would be written once and then never updated.

The other cache Renovate keeps, the repository cache of branch SHAs and pull request state, is deliberately disabled (`RENOVATE_REPOSITORY_CACHE: disabled` in the workflow). Restored across runs it goes stale the moment a Renovate pull request merges: the next run sees the merge happen "during" the run, aborts with "Repository has changed during renovation", and - because an aborted run still saves its cache entry - every following run inherits the same stale state and aborts the same way. Re-reading branch and pull request state from the API each run is a handful of list calls, not the minutes the package cache saves.

## What it needs

| Secret                             | What it is                                        |
| ---------------------------------- | ------------------------------------------------- |
| `RELEASE_PLEASE_GITHUB_APP_ID`     | The app's client ID, shared with release-please   |
| `RELEASE_PLEASE_GITHUB_APP_SECRET` | The app's private key, shared with release-please |

The app's installation needs six permissions. Four are obvious and two are not:

| Permission            | Why                                                                                                                                                                       |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Contents write        | Push the branch                                                                                                                                                           |
| Pull requests write   | Open and update the pull request                                                                                                                                          |
| Issues write          | Maintain the dependency dashboard. Renovate logs an error and carries on without it                                                                                       |
| Workflows write       | Change a file under `.github/workflows/`. Without it, every action-pin update fails and nothing else does                                                                 |
| Commit statuses write | Report a branch's artifact error as a commit status. Without it, the first branch whose artifacts fail aborts the whole run as "Repository has changed during renovation" |
| Checks read           | Read CI results when gating a rebase, which GitHub Actions reports as check runs rather than commit statuses                                                              |

The last two are the failures worth knowing about in advance, because each surfaces as Renovate misbehaving - ignoring GitHub Actions, or claiming the repository moved mid-run - rather than as a permission problem.

## The author is not `renovate[bot]`

Self-hosted under this app, pull requests are authored by the app's own bot rather than by `renovate[bot]`. Anything keyed on the author has to be told the new name, and each of these fails by quietly doing nothing:

- a post-processing workflow gated on `if: github.actor == 'renovate[bot]'`
- a commitlint config selected by pull request author
- a `CODEOWNERS` or branch-protection exemption naming the bot
- a commit-author email pinned to `29139614+renovate[bot]@users.noreply.github.com`

The workflow exposes the right name as `steps.app-token.outputs.app-slug`, so a downstream gate can read `${{ ... }}[bot]` rather than hard-coding it.

## Running it by hand

`workflow_dispatch` takes two inputs: `logLevel`, for a `debug` run when a dependency is not being picked up, and `dryRun`, which resolves every update and logs what it would open without touching the repository. Start with a dry run after changing `renovate.json5` - a bad `customManagers` regex matches nothing and reports success.

## Bringing it up the first time

In this order, because each step fails in a way that looks like the one before it:

1. **The two secrets exist on the repository**, or at the organization with this repository in scope. Without them `create-github-app-token` fails on the first step with an authentication error, which is the one unambiguous failure in the list.
2. **The app is installed on this repository**, with the six permissions in the table above. A missing install looks like Renovate reporting the repository as not found; a missing permission looks like Renovate ignoring one manager.
3. **The labels in `renovate.json5` exist.** GitHub drops an unknown label from a pull request without reporting it, so the branch arrives unlabelled and nothing says why - the same trap the issue forms have. `README.md` has the `gh label create` lines.
4. **Dispatch with `dryRun: true` first.** The log lists every update Renovate resolved and what it would have opened. This is the run that tells you whether the custom-manager regexes match anything, and it cannot leave a half-made branch behind.
5. **Then dispatch for real and read the dependency dashboard**, not the pull request list. The dashboard is the proof the run completed and had write access; an empty pull request list on its own is ambiguous between "working" and "found nothing".

Two things to expect on a repository that has only this workflow in it:

- **`gomod`, `pre-commit` and `docker` find nothing**, because there is no `go.mod`, no `.pre-commit-config.yaml` and no Dockerfile yet. Those managers reporting zero dependencies is correct, not a misconfiguration. Only `github-actions` has anything to read.
- **The first real run should open exactly one pull request**: `actions/create-github-app-token` is pinned at `v3.3.0` here and `v3.4.0` is out. That is a deliberately useful starting state - it means the first run has one unambiguous thing to do, and a run that opens nothing is a signal rather than a shrug. Do not bump that pin by hand before testing.
