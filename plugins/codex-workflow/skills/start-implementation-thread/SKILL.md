---
name: start-implementation-thread
description: Apply after a technical specification, bug report, or UI implementation document has been presented and the user explicitly approves starting implementation, a fix, or debugging in a later message. Choose an appropriate GPT-5.6 model and reasoning effort, then hand implementation to a new Codex task. Do not apply when implementation was requested from another task.
---

# Start an Implementation Task

## Output Language

All user-facing output and every handoff prompt produced by this skill must use the language requested by the user. If no language is specified, use the language of the user's request. Translate every prompt section heading, label, reason, and report; keep code, paths, identifiers, commands, model IDs, API names, and required proper names unchanged.

Re-read the approved implementation-basis document, choose the minimum capable model and reasoning effort, and start implementation in a separate Codex task.

## Start gate

Require all of the following: the technical specification, bug report, or UI implementation document path (and version, if relevant) is known; the document was presented in an earlier response that ended; the user explicitly approved that identified document in a later message; and no unincorporated feedback or blocking decision remains. A PRD alone is not an implementation basis. An initial request such as “document it and then implement it” is intent, not approval.

## Language policy

Write all generated task instructions and documents in the language the user requested; if no language was specified, use the language of the user's request. This applies to every section heading and section label as well as prose, tables, templates, and diagram labels. The English labels named below are semantic slots, not literal output: translate each one into the target language and keep the section names consistent throughout the handoff. Preserve code, commands, paths, identifiers, and required proper names.

## Selection

Read the document from disk immediately before creating the task. Assess scope, remaining design decisions, external systems and platforms, verification difficulty, failure cost, and independently separable work.

| Model | Use when |
| --- | --- |
| `gpt-5.6-luna` | Completion is clear and the change is narrow and repeatable. |
| `gpt-5.6-terra` | Default for confirmed multi-file implementation, ordinary debugging, and standard verification. |
| `gpt-5.6-sol` | Architecture decisions, difficult cross-system or cross-platform interaction, or material security, data, or compatibility risk. |

Choose the lowest sufficient effort: `low` for short fully defined work; `medium` for ordinary implementation; `high` for meaningful trade-offs or multiple modules and tools; `xhigh` for complex or high-risk work; and `max` only for the hardest indivisible work. Never select `ultra` automatically. Honor an explicit user model or effort request, but report unsupported combinations instead of silently changing them.

## Create the task

Only when `thread/start` and `turn/start` are available to the current model:

1. Derive two concise, meaningful task titles from the approved document and the user's request, written in the user's language: one for the originating task that indicates specification or approval, and one for the new task that indicates implementation. Keep the feature or bug subject consistent and distinguish the lifecycle stage in the title.
2. Tell the user the selected model and effort with a one-sentence reason.
3. Rename the current task using the available task/thread title operation before handing off, using the specification/approval title.
4. Call `thread/start` with the selected model and project's absolute `cwd`.
5. Rename the returned task/thread with the implementation title using the available task/thread title operation before starting its first turn.
6. Call `turn/start` with the returned `threadId`, selected `effort`, and a localized prompt containing semantic sections for **Instructions**, **Goal**, **Work to do**, **Do not**, **Constraints and notes**, **Original user request**, and a list of approved document paths. Render every section heading in the user's language; do not emit these English labels merely because they are listed here.

The task prompt must require implementation, verification, and reporting; forbid stopping after document reading or handing work to another implementation task. Preserve the original request and classify extra user instructions into goal, work, prohibitions, and constraints. Do not invent scope, implementation methods, worktrees, or constraints.

If `thread/start` or `turn/start` is unavailable to the current model, still prepare the complete handoff before reporting the limitation. The fallback report and its labels must follow the `Language policy` above: use the language requested by the user, or the language of the user's request when no language was specified. Output all of the following:

- selected model;
- selected reasoning effort;
- the reason for the selection; and
- the complete prompt that would have been passed to `turn/start`.

Do not summarize, shorten, or rewrite that prompt in the fallback report. Put the exact final `turn/start` input in a Markdown code block, including every localized section, the original user request, and the approved document paths. Also state that the unsupported tool surface prevents creating the new task and that implementation will not continue in the current task.

Use this fallback format, translated into the user's language:

```text
The current session does not provide `thread/start` and `turn/start`, so a new implementation task cannot be created.

Selected model: <model>
Selected reasoning effort: <effort>
Selection reason: <reason>

Prompt intended for the new task:
<the complete, exact turn/start input in a code block>
```

## Subagent boundary

The root model and effort selected here are separate from any later internal delegation. Do not preassign subagent count, roles, models, or task decomposition in the handoff prompt. The receiving task independently applies `orchestrate-subagents` after reading the approved documents and latest user direction; if its creation gate fails, it implements directly.

After successfully creating the task and its first turn, do not continue implementation in the current task. If task-creation tools are unavailable, do not implement as a substitute: report the limitation, selected model, reasoning effort, selection reason, and the complete exact handoff prompt.
