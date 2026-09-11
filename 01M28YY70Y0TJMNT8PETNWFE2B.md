---
id: 01M28YY70Y0TJMNT8PETNWFE2B
title: "Snapshot and fetch stamps are per worktree, so every new worktree snapshots at once"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Key the snapshot and fetch schedule stamps by repository instead of by working tree, so the 72h/interval clock is shared across all worktrees of one clone"
context: |-
  The jaira/board snapshot branch gets a commit every 20-30 minutes instead of every 72 hours.

  Cause: the stamp file that records when the background job last ran is stored per working tree. core/ticket/store.go:487 builds the state dir as ~/.jaira/state/<basename>-<sha256(Root)[:4]>, keyed on the working tree path. core/snapshot/schedule.go:23 puts snapshot.json there; core/refsync/schedule.go:34 puts fetch.json there.

  A fresh git worktree is a new path, so it gets a new state dir with no stamp. bgrun.Due (core/bgrun/bgrun.go:106) treats a missing stamp as 'never ran' and fires immediately. Every worktree-agent-* worktree therefore takes a full snapshot on its first jaira command.

  Measured: ~/.jaira/state holds 2694 directories; the agent-* ones each carry their own recent ran_at (17:22, 17:23, 17:24, 18:05, 18:06, 18:34, 18:51, 18:52, 19:17) - exactly the commit spacing on jaira/board.

  The tree-identity guard in core/snapshot/snapshot.go:169 does not help, because tickets really do change between runs, so each run writes a real commit and pushes.

  Known: sessions and locks under the same state dir are correctly per working tree - only the two background-job stamps are repo-wide concerns. git rev-parse --git-common-dir is identical for every worktree of one clone and is the obvious key.
definition-of-done: snapshot.json and fetch.json are keyed by repository (shared by all worktrees of one clone); a second worktree of the same repo does not trigger a snapshot or fetch when the first one ran within the interval; sessions and locks stay per working tree; tests cover the shared-clock behaviour across two worktrees
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-11T19:26:01Z
updated-at: 2026-09-11T19:26:16Z
updated-by: Alexander Sacharov
---

# Snapshot and fetch stamps are per worktree, so every new worktree snapshots at once

## Definition of Done

- [ ] snapshot.json and fetch.json are keyed by repository (shared by all worktrees of one clone); a second worktree of the same repo does not trigger a snapshot or fetch when the first one ran within the interval; sessions and locks stay per working tree; tests cover the shared-clock behaviour across two worktrees

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

