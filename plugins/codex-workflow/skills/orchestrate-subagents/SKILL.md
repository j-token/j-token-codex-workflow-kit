---
name: orchestrate-subagents
description: Apply before spawning subagents or beginning implementation that may benefit from delegation. Delegate only independent, bounded work after a creation gate, with explicit roles, file ownership, validation, and integration by the root agent.
---

# Subagent Orchestration

Do not create subagents by default. The root agent remains responsible for validation, conflict resolution, and final integration. This skill never bypasses document-approval or new-implementation-task gates.

## Creation gate

Spawn a subagent for a current step only when all conditions hold: delegation has a clear context, speed, or quality benefit; the work can run independently from verified inputs; the prompt can state a complete input and output contract; the root can validate and integrate the result; and the benefit exceeds token, concurrency, and coordination cost. Otherwise use zero subagents.

Do not spawn for explanation, status checks, or a local one-file edit. Do not treat `ultra`, file count, or task count alone as justification. If spawning is unavailable, work directly unless the user explicitly requires delegation, in which case report the limitation.

## Roles and selection

Respect explicit user model, headcount, and effort requests, then applicable `AGENTS.md`, then these defaults:

| Role | Model | Best use |
| --- | --- | --- |
| Sol | `gpt-5.6-sol` | Orchestration; ambiguous or high-risk decisions; security, data, compatibility, bug, logic, and edge-case review. |
| Terra | `gpt-5.6-terra` | Evidence-based hypotheses, research interpretation, planning, and specification-conformance validation. |
| Luna | `gpt-5.6-luna` | Broad repository exploration, mechanical transformations, test execution, and fully specified low-judgment implementation. |

Record the selection rationale in the internal work contract. Only claim an actual model assignment when the tool schema supports it and it is observable.

## Prompt and ownership

Default `fork_turns` to `none`; pass only the minimum needed context. Every subagent prompt must have these semantic sections, with every heading translated into the user's language: **Instructions**, **Goal**, **Work to do**, **Do not**, and **Constraints and notes**. These English names are placeholders for the required meanings, not literal headings to emit. State input scope, expected evidence, validation method, ownership, and exact writable files.

For read-only, review, and validation work, state that no files may be changed. For writing work, list allowed files exactly. A subagent must report—not edit—any file outside that list which it thinks requires modification.

Do not give two writers overlapping files. Before spawning, compare their expected write sets; if they overlap, assign one writer and make the others read-only. When later work depends on earlier work, follow: delegate current step → root validates → establish next input and ownership → re-evaluate the gate.

## Integration

Wait for requested results, validate evidence, tests, and changed files against the contract, then integrate. On failure, narrow and retry once or complete the small gap directly. Report whether delegation occurred, why, what was delegated, validation results, conflicts or failures, and any unverified model assignment.
