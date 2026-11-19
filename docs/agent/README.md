# `docs/agent/`

The record of decisions that are not recoverable from the code. Everything here answers a question of the form "why is this odd thing the way it is", asked by someone - a person or an agent - who is about to change it back.

## `choice-*.md` are constraints, not history

This is the convention the root [`AGENTS.md`](../../AGENTS.md) leans on, and the reason the directory exists. A `choice-*.md` file records a tradeoff that was made deliberately and that looks, from the code alone, like an oversight waiting to be cleaned up: a suppression in a scanner config, a dependency pinned below its latest version, a backend selected by one property rather than the obvious one, a test fixture that duplicates something upstream.

Three rules follow from that:

- **Write one when the code cannot carry the reasoning.** A comment is better than a document whenever it fits. What goes here is the decision too long for a comment and too load-bearing to lose - typically one where the alternative was tried, or where the cost of the choice is paid somewhere far from where it was made.
- **Read before reversing.** The root `AGENTS.md` names three or four things that look reversible and are not. If a change here would undo something a `choice-*.md` argues for, the document is the first thing to read and the argument is the thing to answer - not the code.
- **A superseded choice is marked, not deleted.** When the thing a document justified goes away, the reasoning usually still explains something that is still here: a Dependabot ignore entry, a Scorecard annotation, the shape of a package. Mark it superseded, name what superseded it, and leave the argument in place. Deleting it loses the answer to "why is this still configured".

## What else belongs here

- **AI usage disclosure.** How agents are used in this repository and what review they are subject to. Where a release publishes a disclosure manifest, this is the page it is a snapshot of.
- **Design notes an agent wrote for itself** - an investigation, a measurement, a survey of three approaches - where the conclusion is worth keeping but is not a decision. Name these something other than `choice-*`, so the constraint files stay findable.

## What does not

An internal note is not a guide. If an operator would ever run the commands in it, it belongs in [`docs/guides/`](../guides/README.md); if a contributor would, it belongs in [`docs/dev/`](../dev/README.md). The distinction matters because the checking standard is different - a wrong count in a note here costs a reader a minute, and a wrong command in a guide costs an operator an afternoon.
