---
name: pr-rule
description: Apply when the user explicitly asks Codex to create a pull request. Create a draft PR and confirm the target branch with the user.
---

# Pull Request Rules

## Required rules

- All user-facing output and every artifact produced by this skill must use the language requested by the user. If no language is specified, use the language of the user's request.
- Write the PR title and body in the language requested by the user. If none was specified, use the language of the user's request.
- Create the PR as a draft.
- Ask the user to choose the target branch; never select it arbitrarily.
- Put the title only in GitHub's PR title field. Do not repeat it in the body.
- Preserve the template below. Do not add overlapping sections such as separate "Expected result", "Actual result", "Facts", "Assumptions", or a diff-redundant "Changes made" section.

## PR body template

```md
## TL;DR

## Background / context

## Observed problem or goal

## Approach

## Changed files

## Acceptance criteria

## Tests

## Official documentation or references

## Related PRs (optional)
```

Translate the template headings and all prose into the user's requested language when creating the PR.

## Writing rules

- Explain why the PR is needed and the expected behavior in **Background / context**.
- List the issue or goal that the reviewer must assess, with evidence, in **Observed problem or goal**.
- Explain the chosen implementation in **Approach**.
- State the conditions that demonstrate success in **Acceptance criteria**.
- Do not duplicate implementation details already visible in the diff.

## Title and branch format

Use `<branch-prefix>: <concise title>` for the PR title. Branches must use one of `feat/`, `bug/`, `perf/`, `refactor/`, or `misc/` and have a meaningful descriptive suffix.
