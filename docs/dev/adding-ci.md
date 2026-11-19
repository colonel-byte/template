# Adding the test workflows

Six workflows ship with this template: [`check-go-mod`](../workflows/check-go-mod.md), [`codeql`](../workflows/codeql.md), [`commitlint`](../workflows/commitlint.md), [`pre-commit`](../workflows/pre-commit.md), [`scan-lint`](../workflows/scan-lint.md), and [`scorecard`](../workflows/scorecard.md), plus the two release workflows. **Nothing runs the tests**, because a test workflow for a repository with no tests is a green check that means nothing.

The two that are deliberately left to write are below. Both are small; what matters is the handful of decisions each one encodes, which is why they are notes rather than files.

When you add either, three things change in the same pull request:

1. The workflow, under `.github/workflows/`.
2. A page under `docs/workflows/` named after it - required by [`.github/AGENTS.md`](https://github.com/OWNER/REPO/blob/main/.github/AGENTS.md), then `go run ./magefiles/core generate:document` to put it in the summary.
3. The `Which CI job` dropdown in [`.github/ISSUE_TEMPLATE/01-failing-build.yml`](https://github.com/OWNER/REPO/blob/main/.github/ISSUE_TEMPLATE/01-failing-build.yml), which lists the `name:` a workflow shows in the checks list. That one goes stale silently.

## Unit tests

Triggers on `pull_request` and `merge_group`, checks out, `setup-go` with `go-version-file: go.mod`, and runs the suite. Four decisions to make deliberately:

- **Name the packages, or name the exclusion.** `go test ./...` picks up `magefiles/` at the root, which cannot build outside mage. Either list the package paths or pass `-short` and let the suites gate themselves.
- **`-short` is not optional once there are e2e suites in the module.** Without it, a unit-test job pulls in suites that want Docker or a cluster and either hangs or fails for reasons that have nothing to do with the change. This is the rule the `TEMPLATE:` stub in the root [`AGENTS.md`](../../AGENTS.md) is waiting for - write it there too, naming the forbidden invocation.
- **Build offline, like every other job.** `GOFLAGS=-mod=vendor`, `GOPROXY=off`. A test run that silently reaches a proxy hides a vendor tree that is missing a file.
- **Replay the fuzz seed corpus if there are fuzz targets.** `go test -run=Fuzz` replays the committed corpus without fuzzing, which is cheap and catches a regression on an input that already found a bug once. It is not the same as a fuzzing run, and should not be confused for one in the job name.

```yaml
# The shape, not the whole file.
- name: Run unit tests
  run: go test -short -mod=vendor ./cmd/... ./internal/... ./pkg/...
  env:
    GOFLAGS: -mod=vendor
    GOPROXY: "off"
```

## End-to-end tests

Worth writing only once there is something to drive end to end, and worth treating as the most expensive workflow in the repository - cargoship's is 154 lines and provisions a fleet.

- **Decide what provisions the fixture**, and make the job's name say which. A suite that brings up containers, a suite that brings up VMs, and a suite that runs against hosts that already exist are three different workflows' worth of setup, and a result is only as good as the backend that produced it.
- **Check what skipped, not just what failed.** A skip on an e2e suite usually means a capability the environment does not have, and a run with fifteen of them read as a pass is the commonest way this kind of suite misleads. Make the job surface the skip count.
- **Do not make it a required status check until it is reliable.** An e2e job that flakes becomes a job people merge past, which is worse than not having it.
- **Upload diagnostics on failure.** `actions/upload-artifact` with the suite's logs is the difference between a reproducible report and "it failed in CI". The flaky-test issue form asks for exactly that artifact.

If the suite needs credentials or a registry, give it its own workflow rather than adding secrets to the unit-test job - a fork pull request cannot see secrets, and a unit-test job that silently degrades on forks is a check that means different things depending on who opened the pull request.
