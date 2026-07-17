---
name: bug-report-to-fix
description: Apply when a user reports a bug, unclear symptom, screenshot, log, or debugging request. First create a bug report in `.codex/temp`, gather reproduction information, and begin code changes only after a later explicit approval.
---

# Bug Report to Fix

## Workflow

1. Separate observed facts from assumptions.
2. Gather the minimum reproduction information: action, timing, expected result, actual result, platform, and recent changes.
3. Ask for the smallest concrete action that can narrow the cause when necessary.
4. Maintain a temporary debugging document under `.codex/temp`.
5. Present the report and end the turn. Only after the user explicitly approves starting a fix or debugging in a later message, apply `start-implementation-thread` with the report path.
6. In the new task, verify the fix through reproduction, tests, logs, or the closest available check.

## Language policy

Write every user-facing and generated document—including headings, tables, templates, diagram labels, and prose—in the language the user requested. If no language was explicitly requested, use the language of the user's request. Keep code, commands, identifiers, and required proper names unchanged.

## Platform rule

Do not assume a problem is platform-specific. For example, “tested on iOS” only establishes that iOS was tested; it does not establish that Android is unaffected.

## Temporary document

Create `.codex/temp/` at the repository root when it does not exist. Follow `cognitive-writing` and use:

```text
.codex/temp/YYYYMMDD-HHMM-bug-<topic>-workflow.md
```

Use this localized template:

```md
# Debugging work document

## TL;DR
## Observed symptoms
## User actions
## Expected result
## Actual result
## Confirmed environments
## Environments not yet checked
## Attached material
## Facts
## Assumptions
## User confirmation needed
## Reproduction path
## Candidate causes
## Files / logs to investigate
## Proposed fix
## Acceptance criteria
## Fix result
```

## Code-change gate

Before explicit later approval, limit work to reporting, investigation questions, and a reproduction plan. A first message such as “write a bug report and fix it” is not approval. After presenting the report, require a separate message that identifies its path or version and explicitly approves the fix; then hand off through `start-implementation-thread`, rather than editing code in the current task.
