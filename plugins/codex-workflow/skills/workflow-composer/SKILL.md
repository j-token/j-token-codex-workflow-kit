---
name: workflow-composer
description: Apply when requirements discovery, debugging, and UI or Figma implementation are mixed in one request. Combine the relevant workflow modules without forcing one classification, keep working documents in `.codex/temp`, and sequence questions clearly.
---

# Workflow Composer

## Output Language

All user-facing output and every working document or prompt produced by this skill must use the language requested by the user. If no language is specified, use the language of the user's request. Translate headings, labels, questions, templates, prompts, and reports; keep code, paths, identifiers, commands, skill names, and required proper names unchanged.

Use this as the entry point for compound requests. Detect and apply every relevant module:

| Signal | Module |
| --- | --- |
| Context, desired outcome, rough feature, policy, product or feature specification | `requirements-to-spec` → `prd-writer` → `technical-spec-writer` |
| Symptom, error, screenshot, log, reproduction, or “it does not work” | `bug-report-to-fix` |
| Figma link, screen, design, UI flow, screenshot, or visual material | `figma-flow-to-implementation` |

Use supporting skills as their triggers apply: `cognitive-writing`, `branch-rule`, `commit-rule`, `pr-rule`, `git-push-safety`, `start-implementation-thread`, and `orchestrate-subagents`. Disable a module only when the user explicitly excludes it or asks for a narrower mode.

## Language policy

Write all generated documents in the language the user requested; if no language was specified, use the language of the user's request. Apply this to every section heading and section label as well as prose, tables, templates, and diagram labels. Translate semantic template labels too; do not leave English section names in an otherwise localized artifact. Preserve code, commands, identifiers, and required proper names.

## Sequencing and documents

Follow the order in which the user introduced work items. Finish the required questions and document update for one item before moving to the next; do not mix unrelated questions unless one response genuinely resolves multiple modules.

Unless the user specifies another location, create `.codex/temp/` at the repository root and use `cognitive-writing`:

```text
.codex/temp/YYYYMMDD-HHMM-<type>-<topic>-workflow.md
```

Use `feat`, `bug`, `perf`, `refactor`, or `misc` for `<type>`, a recognizable slug for `<topic>`, and numeric suffixes to prevent same-minute collisions.

## Final-document and approval rules

- Create final documents only when the user requests them or explicitly confirms the working document.
- For product and feature work, write a PRD first, present it, wait for separate approval, then write and present a technical specification. Keep them separate.
- For bug and UI work, maintain the single implementation-basis document specified by the corresponding skill.
- Separate facts, assumptions, open questions, decisions, and exclusions.
- Presenting a document ends the current turn; future-tense intent in the original request is not approval.
- On later explicit approval, hand the approved technical specification or relevant bug/UI document to `start-implementation-thread`; do not implement in the current task.
- The receiving implementation task, not this task, independently decides whether to use subagents through `orchestrate-subagents`.
