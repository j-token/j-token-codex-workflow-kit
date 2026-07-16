---
name: prd-writer
description: Create or refine a product requirements document (PRD), product proposal, feature requirements, or product-facing technical plan. Use to define why and what to build, scope, decisions, and acceptance criteria before technical design.
---

# Product Requirements Document Writer

Create a decision-ready PRD. Define the user and business problem, desired outcome, scope, requirements, and observable acceptance criteria. Do not turn the PRD into a technical implementation specification.

## Language policy

Write the complete PRD in the language the user requested. If no language was specified, use the language of the user's request. Localize headings, tables, templates, examples, and diagram labels. Preserve code, commands, APIs, identifiers, product names, and required proper names.

## Process

1. Read the request, existing working document, and relevant product context. Research facts that materially affect the decision.
2. Separate confirmed requirements, assumptions, open questions, decisions needed from the user, and exclusions.
3. Ask only blocking questions. Do not invent platform support, policy, user segments, legal requirements, metrics, or integrations.
4. State the smallest viable scope and explicitly record out-of-scope items.
5. Write concrete, testable acceptance criteria. Distinguish functional, UX, reliability, privacy, security, accessibility, and analytics requirements where relevant.
6. Present the PRD and end the turn. A later technical-specification workflow must re-read the confirmed PRD from disk.

## Required PRD structure

Use these sections, localized for the user:

```md
# <Product or feature> PRD

## TL;DR
## Problem and background
## Goals and non-goals
## Target users and scenarios
## Scope
## User journeys / requirements
## Functional requirements
## Non-functional requirements
## Success metrics
## Dependencies and constraints
## Risks and open questions
## Decisions required
## Acceptance criteria
## Rollout / measurement plan
## References
```

Omit a section only when it is genuinely not applicable, and say why when its absence could otherwise be misleading.

## Requirement quality

- Use unambiguous language: identify actor, trigger, expected behavior, boundary, and failure behavior.
- Prefer examples and scenarios for ambiguous behavior.
- Make success metrics measurable and specify the baseline, target, owner, and measurement window when known; otherwise label them as open.
- Do not promise implementation details, timelines, or outcomes unsupported by evidence.
- Use a Mermaid flow only when a journey, state, or dependency cannot be understood clearly from concise prose.

## Handoff

The PRD is the authoritative product decision record. It becomes input to `technical-spec-writer`, not authorization to implement. Preserve user decisions and excluded approaches so the technical specification can trace them.
