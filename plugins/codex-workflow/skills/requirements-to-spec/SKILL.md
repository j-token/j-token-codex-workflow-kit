---
name: requirements-to-spec
description: Apply when a user brings rough requirements, context, open questions, or asks for product or feature documentation. Create a working document in `.codex/temp`, resolve uncertainty, write and confirm a PRD, then write a technical specification.
---

# Requirements to Specification

## Workflow

1. Restate the current understanding.
2. Classify items as confirmed, needs research, needs a user decision, or excluded.
3. Inspect code, documentation, or external sources when a fact affects implementation.
4. Update the temporary document in `.codex/temp`.
5. Once requirements are ready, apply `prd-writer` and present the PRD.
6. End the turn. In a later message, require approval that identifies the PRD path or version before applying `technical-spec-writer`.
7. Re-read the approved PRD from disk and write a separate technical specification.
8. Present it and end the turn. Require a later explicit implementation approval before applying `start-implementation-thread`.

## Language policy

Write all generated documents in the language the user requested; if none was stated, use the language of the user's request. Localize headings, tables, templates, and diagram labels as well as prose. Do not translate code, commands, identifiers, or required proper names.

## Documents

Create `.codex/temp/` when needed and follow `cognitive-writing`.

```text
.codex/temp/YYYYMMDD-HHMM-feat-<topic>-workflow.md
.codex/temp/<product>-prd.md
.codex/temp/<product>-technical-spec.md
```

Use `misc` instead of `feat` for a small mixed cleanup. Keep the PRD and technical specification as separate files: the PRD establishes why and what to build, scope, and acceptance criteria; the technical specification establishes implementation, contracts, and verification.

Use this localized working-document template:

```md
# Requirements / implementation work document

## Current goal
## Background
## Confirmed requirements
## Open questions
## Research findings
## Decisions
## Excluded approaches
## Implementation plan
## Acceptance criteria
## Next confirmations
```

## Rules

- Do not finalize a specification from a single rough message.
- Ask only blocking questions; research other unknowns first.
- Move user-approved items to **Decisions** and rejected ideas to **Excluded approaches**.
- Do not merge the PRD and technical specification.
- Do not write the technical specification in the PRD presentation response, and do not implement in the technical-specification presentation response.
- A first request such as “write the PRD and spec, then implement” is not later approval. On later explicit approval identifying the technical specification, use `start-implementation-thread` rather than changing code in the current task.
