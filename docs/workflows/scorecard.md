# Scorecards supply-chain security

[`.github/workflows/scorecard.yaml`](https://github.com/OWNER/REPO/blob/main/.github/workflows/scorecard.yaml) runs [OpenSSF Scorecard](https://securityscorecards.dev/) over the repository and publishes the result three ways: to the public API behind the `README.md` badge, to GitHub's code-scanning dashboard as SARIF, and as a run artifact kept for five days.

It fires on pushes to `main`, weekly on Saturday at 01:30 UTC, and whenever a branch protection rule changes - that last trigger exists because the Branch-Protection check reads settings rather than files, so a protection change is the only thing that would otherwise move the score with no commit behind it.

## It only works on the default branch

Scorecard scores `main` and nothing else. A pull request cannot move the published score, and the score a reviewer sees is always the last merge's. To check a change before it lands, score the working tree by hand instead - [`.github/AGENTS.md`](https://github.com/OWNER/REPO/blob/main/.github/AGENTS.md) has the `--local` invocation and the handful of checks worth running.

## `publish_results: true` is what the badge needs

The workflow can run, fill in the code-scanning dashboard, and leave the badge reading `unknown` indefinitely - that is what `publish_results: false` produces. Publishing is what writes the score to `api.scorecard.dev`, which is the endpoint the badge image comes from.

Two further conditions, neither of which the workflow can satisfy on its own:

- **The repository has to be public.** The public API holds no score for a private repository, so the badge cannot resolve.
- **The first score takes a run.** Until the workflow has completed once on `main`, the badge reads `unknown`, and that is not a misconfiguration.

## Permissions

`permissions: read-all` at the top, with two writes on the job:

| Permission               | Why                                                                      |
| ------------------------ | ------------------------------------------------------------------------ |
| `security-events: write` | Upload the SARIF to the code-scanning dashboard                          |
| `id-token: write`        | The OIDC token the published result is signed with - no badge without it |

`persist-credentials: false` on the checkout is deliberate and Scorecard's own Token-Permissions check looks for it.

## Annotating a check that is wrong rather than fixing it

Some findings are correct about the file and wrong about the project - a vendored dependency's own CI script tripping Pinned-Dependencies, say, which no build here ever invokes. Those go in `.github/scorecard.yml`, which this template deliberately does **not** ship, because an annotation with no reason behind it is worse than a lower score.

The shape, when one is needed:

```yaml
annotations:
  - checks:
      - pinned-dependencies
    reasons:
      # Why this finding does not apply here. Name the file, say what it is, and say why
      # nothing in this repository's build or CI path runs it.
      - reason: not-applicable
```

Two rules about that file, both learned the hard way:

- **Every annotation carries its reason as a comment**, not just the enum. `not-applicable` on its own tells the next reader nothing, and the finding will look like an oversight to be cleaned up.
- **An annotation outlives what justified it.** When the dependency goes or the script stops being vendored, the entry is still there, still suppressing a real finding. The root [`AGENTS.md`](../AGENTS.md) treats removing one as the operator's call rather than a documentation change - flag it, do not quietly delete it.
