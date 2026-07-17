---
name: technical-spec-writer
description: Create or refine a technical specification, implementation specification, architecture specification, or API, SDK, CLI, IPC, runtime, build, and test contract. Use after product requirements are confirmed to define how the work will be implemented and verified.
---

# Technical Specification Writer

Turn an approved PRD or equivalent confirmed requirements into an implementable, reviewable technical contract. Re-read the source document and repository context; do not rely on a conversation summary.

## Language policy

All user-facing output and every artifact produced by this skill must use the language requested by the user. If no language is specified, use the language of the user's request.

Write the full specification in the language the user requested. If none was specified, use the language of the user's request. Localize headings, tables, templates, examples, and Mermaid labels. Leave code, commands, API names, identifiers, paths, and required proper names unchanged.

## Process

1. Extract confirmed requirements, acceptance criteria, non-goals, decisions, and constraints from the approved source document.
2. Inspect affected code, interfaces, data models, build and deployment paths, tests, and external documentation as needed.
3. Separate verified facts from design choices and unresolved questions. Do not present assumptions as repository facts.
4. Define the smallest safe design: boundaries, contracts, state changes, errors, compatibility, migration, observability, security, and rollback when relevant.
5. Specify tests and acceptance checks that can demonstrate each important requirement.
6. Present the specification and end the turn. Implementation requires a later, separate explicit approval and must use `start-implementation-thread`.

## Required structure

Use these semantic sections, translating every heading into the user's language (the English names below are placeholders, never mandatory output):

```md
# <Feature> technical specification

## TL;DR
## Source requirements and scope
## Existing-system findings
## Goals and non-goals
## Architecture and data flow
## Components and responsibilities
## Interfaces and contracts
## Data model and state transitions
## Error handling and edge cases
## Security, privacy, and permissions
## Compatibility, migration, and rollback
## Observability
## Implementation plan
## Test plan and acceptance checks
## Risks, trade-offs, and open questions
## References
```

Remove inapplicable sections only when their omission is clear and safe. Use diagrams for genuinely complex flows, with descriptive localized labels.

## Specification rules

- Link requirements to implementation responsibilities and verification.
- For every contract, define inputs, outputs, validation, error behavior, ownership, and compatibility expectations.
- Identify exact files or modules only after inspecting them; otherwise label them as candidates.
- State changes in ordered steps with dependencies and rollback implications.
- Avoid vague directions such as “handle errors properly”; name the failure mode, user-visible behavior, logging, and test.
- Do not silently enlarge scope. Put optional improvements and unresolved decisions in their own section.

## Handoff

This document is the implementation basis. It is distinct from the PRD and must be explicitly approved in a later user message before implementation starts.
