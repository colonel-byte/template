# PR Title Check

[`.github/workflows/commitlint.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/commitlint.yaml) lints the **pull request title** against [conventional commits](https://www.conventionalcommits.org/). It is a standard across the organization, and it is load-bearing rather than cosmetic: pull requests are squash-merged, so the title becomes the commit message on `main`, and release-please reads those messages to decide the next version and to write `CHANGELOG.md`. A title with the wrong type produces a release with the wrong version number.

This is why the job lints the title and not the commits. `wagoid/commitlint-github-action` lints commit messages, which are the ones squashing throws away, so it is not a drop-in replacement for this step.

## The two configs

| File                                        | Applies to                          | Difference                 |
| ------------------------------------------- | ----------------------------------- | -------------------------- |
| `.github/commitlint/.commitlintrc.yaml`     | Human authors                       | `header-max-length` of 100 |
| `.github/commitlint/.commitlintrc.bot.yaml` | Any author whose login ends `[bot]` | `header-max-length` off    |

A bot cannot shorten a module path, and `build(deps): bump <module> from X to Y` runs past 100 characters for any dependency with a long import path. Turning the rule off for everyone would lose the one thing it buys - a readable `git log` - so it is off for bots only, and every other rule still applies to them.

**The two files must stay in step on every rule but that one.** Commitlint cannot extend a YAML sibling, so there is no shared base and the duplication is real; a rule added to one and not the other is a rule humans or bots quietly escape.

## Why the bot match is a suffix, not a name

The author switch matches `*[bot]` rather than a list of names. Self-hosted Renovate opens pull requests as whichever GitHub App minted its token - see [`renovate.md`](renovate.md) - so the name is not knowable from inside this file, and a hard-coded one fails by silently linting a bot's title against the strict config. The suffix is assigned by GitHub and cannot be spoofed by a human account.

## It reports success on `merge_group`

The job is intended to be a required status check, which means the merge queue waits for a `title_check` result on the `merge_group` ref as well. There is no pull request under that event, so every step is skipped and the job passes - the title was already linted on the pull request that got queued. Skipping the whole job instead would leave the check with no result at all, which blocks the queue forever.

The same reasoning applies to bots: the job runs for them with a laxer config rather than being skipped, because a skipped job produces no result for a required check.

## Bumping commitlint

`package.json` and `package-lock.json` in `.github/commitlint/` pin `@commitlint/cli` and `@commitlint/config-conventional`, and the workflow runs `npm ci`, so the lockfile is what is installed. Renovate's `npm` manager owns both - there is no bump procedure to follow by hand.
