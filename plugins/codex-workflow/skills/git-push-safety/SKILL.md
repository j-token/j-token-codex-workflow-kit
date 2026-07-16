---
name: git-push-safety
description: Apply when Codex is asked to execute, prepare, review, explain, force-push, set an upstream for, publish a branch with, or prepare a PR for `git push`. Prevent accidental pushes to the wrong or protected branch. Push only after explicit user authorization for that push.
---

# Git Push Safety Rules

## Authorization

Execute `git push` only when the user explicitly requests or approves it for the current push. A conceptual question requires explanation only; never run a command for it.

## Required checks

Before every push:

1. Confirm current-turn user authorization.
2. Check the current local branch with `git branch --show-current`.
3. Check upstream tracking with `git branch -vv`.
4. Verify the intended remote branch is not `main`, `master`, or `develop`.
5. Use an explicit `<local>:<remote>` refspec.
6. For a force push, inspect `git log --oneline -5` and use only `--force-with-lease`.

## Commands

Use explicit source and destination branches:

```bash
git push origin feat/my-feature:feat/my-feature
git push --force-with-lease origin feat/my-feature:feat/my-feature
```

Never use branch-only push forms such as `git push origin feat/my-feature`. Never push directly to a protected branch; use a PR or MR.

If the current branch tracks a protected or otherwise incorrect remote branch, do not push in that state. Reset upstream only when the intended remote branch already exists and is verified safe:

```bash
git branch --set-upstream-to=origin/<current-branch>
```

If it does not exist, publish it with the explicit refspec.

## Response rule

Before executing a user-requested push, report the checked local branch, upstream status, destination refspec, and whether it is a force push.
