# New repository checklist

Generating a repository from a template copies **files, not settings**. Everything in this template is a file, so it all arrives; everything in the Settings tab does not, and several of the files that arrive are inert until a matching setting exists. A `CODEOWNERS` that requires nothing and a workflow that goes green while the site 404s both look fine from the file list.

Two more things about generation, both easy to be surprised by: the new repository starts with a **single initial commit** - no history, tags, releases, or issues - and there is **no ongoing link to the template**. Nothing syncs a later template improvement forward; it has to be copied by hand.

## Do this first

The placeholder sweep is the first section of the root [`AGENTS.md`](../../AGENTS.md), and it comes before anything below. A badge or a `book.toml` URL left on the template's slug reports another repository's state as this one's.

## Settings that do not travel

| Setting                        | What stays broken without it                                                                                                                                  | Where                                      |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| Branch protection or a ruleset | `.github/CODEOWNERS` requires nothing - the file alone is advisory. `commitlint.yaml`'s `merge_group` pass-through is dead code until the check is *required* | Settings → Rules                           |
| Actions secrets                | `renovate.yaml` and `release-please.yaml` fail at the first step; `release.yaml` fails at signing                                                             | Settings → Secrets and variables → Actions |
| The GitHub App installation    | Separate from the secrets, and the one that fails partially: without `Workflows: write`, every action-pin pull request fails and nothing else does            | Org → Settings → GitHub Apps               |
| Labels                         | Every label an issue form or `renovate.json5` names. GitHub drops an unknown label **silently**, so the issue arrives unlabelled with no error                | `gh label create`, below                   |
| Issues enabled                 | Renovate's dependency dashboard logs an error and carries on without it                                                                                       | Settings → Features                        |
| Public visibility              | The Scorecard badge cannot resolve - there is no published score for a private repository                                                                     | Settings → Danger Zone                     |
| GitHub Pages                   | The book 404s while `deploy-book.yaml` goes green. Source: deploy from branch `gh-pages`, folder `/`                                                          | Settings → Pages                           |

**Two of these are worth fixing once rather than per repository.** An **org-level repository ruleset** targeting all repositories, or a name pattern, gives every new repository its required checks and review rules at creation with no per-repo step. **Org-level Actions secrets** with visibility set to all repositories do the same for the app credentials and the cosign key. If both exist, most of this table stops being a checklist.

## The parts that can be scripted

```sh
REPO=OWNER/REPO

# Labels the issue forms name.
for l in bug ci test-flake enhancement documentation stack-task; do
  gh label create "$l" --repo "$REPO"
done

# Labels renovate.json5 names.
for l in renovate/oci renovate/action renovate/gomod renovate/npm renovate/pre-commit security; do
  gh label create "$l" --repo "$REPO"
done

# Secrets. The app credentials are shared with release-please; the cosign pair signs the
# release archives. Omit any that are already org-level secrets visible to this repository.
gh secret set RELEASE_PLEASE_GITHUB_APP_ID --repo "$REPO"
gh secret set RELEASE_PLEASE_GITHUB_APP_SECRET --repo "$REPO" < path/to/app-private-key.pem
gh secret set COSIGN_KEY --repo "$REPO" < path/to/cosign.key
gh secret set COSIGN_PASS --repo "$REPO"
```

`COSIGN_KEY` and `COSIGN_PASS` are the **secret** names; `release.yaml` maps them onto `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD`, which are cosign's own variable names and are not renameable. Renaming one side only produces an error at the signing step rather than a missing-secret error.

## Verifying it rather than assuming it

Each of these answers one row of the table, and each has a failure mode that is silent otherwise:

```sh
gh label list --repo "$REPO"                    # every label the forms and renovate name
gh secret list --repo "$REPO"                   # repo-level only; org secrets do not appear here
gh api "repos/$REPO/rules/branches/main"        # [] means nothing is enforced on main
gh api "repos/$REPO/pages" --jq .html_url       # 404 until Pages is enabled
```

Then open `https://github.com/OWNER/REPO/issues/new/choose` on a branch after changing an issue form. A malformed form is reported as a banner on that page and nowhere else - parsing as YAML is not enough to know it renders.

## First run, in order

Each step fails in a way that looks like the one before it, so do them in this order:

1. **Push to `main`.** This is what first runs `scorecard.yaml` and `deploy-book.yaml`, and what gives the badges something to resolve against.
2. **Dispatch Renovate with `dryRun: true`.** The log lists every update it resolved without touching the repository - see [`renovate.md`](../workflows/renovate.md). A bad `customManagers` regex matches nothing and reports success, so this is the run that tells you.
3. **Open a throwaway pull request.** Five checks should appear: `Validate Lint`, `Validate pre-commit hooks`, `CodeQL`, `PR Title Check`, `Check for go mod changes`. A check that does not appear is a workflow whose trigger or filename is wrong, not a flake.
4. **Mark the checks required**, once they pass. Until then `commitlint.yaml`'s `merge_group` branch does nothing, and `CODEOWNERS` is a suggestion.
5. **Let release-please open its first pull request.** Merging it pushes the first tag, which is what `release.yaml` fires on, and what turns the Latest Release badge from `no releases` into a version.

## Still to write, per repository

[`adding-ci.md`](adding-ci.md) covers the two workflows this template deliberately does not ship - unit tests and end-to-end - and the three things that change in the same pull request when you add either.
