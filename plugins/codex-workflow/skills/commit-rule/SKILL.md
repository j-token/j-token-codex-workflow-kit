---
name: commit-rule
description: Apply when the user explicitly asks Codex to create a Git commit. Require explicit authorization for every commit.
---

# Commit Rules

## Authorization

- Create a commit only after the user explicitly authorizes it for the current request.
- Do not infer permission to commit from a request to change code.
- Treat committing and pushing as separate actions; confirm authorization for each.

## Language

Write commit messages in the language the user requested for this task. If the user did not specify one, use the language of their request.

## Scope

Make each commit independently reviewable and limited to one coherent change. Split unrelated work into separate commits.

Use one of these prefixes:

- `feat:`
- `bug:`
- `perf:`
- `refactor:`
- `misc:` for minor unrelated changes.
