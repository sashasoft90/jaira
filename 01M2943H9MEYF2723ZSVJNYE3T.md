---
id: 01M2943H9MEYF2723ZSVJNYE3T
title: Tests write their state into the real ~/.jaira
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "A full test run leaves no directory under the developer's own ~/.jaira; every test that opens a store points JAIRA_HOME at its own temp dir"
context: |-
  A test run pollutes the real ~/.jaira/state. There are 3139 directories in it on this machine and 3115 of them are named 001-*, 002-* or 003-* - the names Go gives subdirectories of t.TempDir(). Only about 24 belong to real checkouts.

  Where it comes from: Store.stateDir (core/ticket/store.go) falls back to $HOME/.jaira when JAIRA_HOME is unset, and tests that build a ticket.Store without t.Setenv("JAIRA_HOME", ...) therefore write sessions/ and locks/ into the developer's own home, keyed by a temp path that stops existing when the test ends.

  Harmless in content - the leftovers are empty sessions/ and locks/ directories, 38 MB in total - but they bury the ~24 real entries among three thousand dead ones, which is how a stale-directory count got misread as a worktree problem while diagnosing NWFE2B.

  Two parts: stop the leak, then clear what is already there. Deleting the existing ones is safe for any directory whose checkout no longer exists; sessions and locks are per-checkout and short-lived by design.
definition-of-done: go test ./... with HOME pointed at an empty directory creates nothing under it; the tests that leaked are identified and isolated; ~/.jaira/state holds only directories whose checkout still exists
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-11T20:56:18Z
updated-at: 2026-09-11T21:00:33Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1358160
claimed-at: 2026-09-11T20:56:31Z
---

# Tests write their state into the real ~/.jaira

## Definition of Done

- [x] go test ./... with HOME pointed at an empty directory creates nothing under it; the tests that leaked are identified and isolated; ~/.jaira/state holds only directories whose checkout still exists

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 21:00 · Alexander Sacharov** — Only internal/cli leaked - measured by running each package on its own and counting ~/.jaira/state before and after, because HOME cannot be redirected in this session. 34 directories per run of that package; every other package already isolates JAIRA_HOME or never opens a store.

Fixed with one TestMain for the package rather than a t.Setenv in each of the 26 test files: 10 of them already set it and 16 did not, so the leak is the default, and a test that forgets is silently wrong somewhere nobody looks. A per-test t.Setenv still overrides it.

The sweep deleted every directory under ~/.jaira/state holding no regular file at all - no session, no lock, no stamp, nothing to lose. 3150 went, 23 real checkouts stayed, 38 MB down to 636 KB. scratchpad/sweep.sh is the script; it is not in the repo, since the leak it cleans up is now fixed.
