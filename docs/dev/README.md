# `docs/dev/`

Pages for the person **changing** this project: how it builds, how it is tested, how a release is cut, and what the task runner does. The audience is a contributor with the repository checked out, which is what separates this directory from [`docs/guides/`](../guides/README.md).

## What belongs here

- **Build flags, and what each one measured.** The flags a release is built with are not obvious from the code, and a size or performance number compared against a baseline built with different flags is worse than no number at all. A page here is the one safe reference for how to produce a comparable build, and the rule is to read it rather than to remember the flags.
- **The test harness.** What each suite covers, what it needs installed, and - the part that actually gets people - what each backend or environment **does not** exercise. A suite that skips fifteen steps because the environment cannot reach them is not a pass, and the page saying so is what stops it being read as one.
- **The task runner reference.** One entry per target, what it reads and what it writes. Several rules in the root [`AGENTS.md`](../../AGENTS.md) and in [`docs/AGENTS.md`](../AGENTS.md) defer to this page for which target regenerates what.
- **Release mechanics**, where they go beyond what [`docs/workflows/`](../workflows/) covers - the artifact set, the signing keys, how to run the release tooling locally without cutting a release.

## Two things to get right

- **Name the forbidden invocation, not just the right one.** A page that says how to run the tests and not which invocation pulls in the suites that need real infrastructure gets the wrong one run first. The root `AGENTS.md` carries the short version; this is where the reason lives.
- **A measured number needs the date and the method.** Every count, size, and pin in a tree defaults to having been measured once and never again. Say what was measured, with what flags, against what baseline - and re-measure rather than quoting an earlier run, since toolchain and dependency drift move absolute numbers with nobody editing a file.

## What does not belong here

Operator-facing walkthroughs are [`docs/guides/`](../guides/README.md). The reasoning behind a decision - as opposed to the mechanics of a process - is [`docs/agent/`](../agent/README.md); a `choice-*.md` explains why the build does an odd thing, while a page here explains how to do it.
