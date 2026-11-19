# `docs/guides/`

Long-form pages for the person **running** this project, not the person changing it. A guide explains a task end to end: what it is for, the commands in order, what each one needs beforehand, and what the output means.

## This is where help text overflows to

Flag descriptions are capped at 125 characters, and `--help` as a whole has a budget - see the "Help text has a budget" section of the root [`AGENTS.md`](../../AGENTS.md). When a flag needs more explanation than a sentence, the explanation lives here and the flag string names the page:

```go
CmdCreateFlagContentStore = "`PATH` of a content store to resolve images and files from instead of the network. See docs/guides/content-store.md."
```

That makes the relationship two-way, and both directions can go stale silently:

- **A guide named from a flag string has to exist, under that exact path.** Nothing in the build checks a path inside a Go string literal.
- **Content truncated out of help text has to already be in the guide.** Check the guide actually says the thing before deleting the sentence that said it - a cap is not a reason to lose information, and at least one truncation in this organization's history had to move a sentence into a guide first.

A walkthrough, a worked scenario, or an enumeration of how several flags interact is a guide rather than a `Long`. So is an example per flag combination, when `--help` should only carry the three to five common cases.

## These are the pages checked hardest

A command printed here is copied and pasted. The cost of a wrong one is not a moment of confusion, it is an operator assuming they mistyped before assuming the guide is wrong - so the "A doc's claims are checkable" rules in the root `AGENTS.md` apply most strictly to this directory. Before changing a guide, run what it prints.

Where a page pins this project's own version - an image tag, a release archive's filename - it needs release-please annotations *and* an entry in the release-please config, or it goes stale on the next release. [`docs/AGENTS.md`](../AGENTS.md) has the exact form, including why the per-line marker is the wrong one to use inside a fenced block.

## What does not belong here

Anything a contributor needs and an operator does not: build flags, test harness setup, release mechanics. That is [`docs/dev/`](../dev/README.md). Decisions and their tradeoffs are [`docs/agent/`](../agent/README.md).
