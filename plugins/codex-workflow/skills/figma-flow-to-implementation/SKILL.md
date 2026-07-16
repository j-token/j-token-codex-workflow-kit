---
name: figma-flow-to-implementation
description: Apply when a user supplies a Figma link, screenshot, visual material, or asks to implement a UI or screen. Infer screen roles and transitions, confirm a Mermaid flow, and create one UI specification / implementation document before implementation.
---

# Figma Flow to Implementation

## Workflow

1. Inspect all provided visual material. Do not treat link order as screen-transition order.
2. Give screens meaningful user-facing names rather than raw Figma IDs.
3. Identify needed frames, icons, logos, illustrations, and other assets.
4. When Figma access is available, obtain needed node JSON, screen images, and assets from the REST API after the user completes OAuth authorization.
5. Write the inferred transitions as a Mermaid flowchart and ask the user to correct screen order, missing items, and transition conditions.
6. Maintain a temporary UI working document.
7. Create the final UI specification / implementation document only when requested or confirmed, present it, and end the turn.
8. On later explicit approval identifying that document, use `start-implementation-thread` to implement it in a separate task.

## Language policy

Write every generated document in the language the user requested; if none was specified, use the language of the user's request. Localize headings, templates, tables, and Mermaid labels in addition to prose. Preserve code, commands, Figma IDs, identifiers, and required proper names.

## Inputs and documents

Prefer the strongest available input: Figma plugin context, Figma link, screenshots, screen recordings, then a user-written screen description. Create `.codex/temp/` when needed, follow `cognitive-writing`, and use:

```text
.codex/temp/YYYYMMDD-HHMM-feat-<topic>-ui-workflow.md
```

Use this localized template:

```md
# UI work document

## Screen inventory
## Inferred screen flow
## Mermaid flowchart
## Key elements by screen
## Required Figma assets
## Figma REST API collection plan
## User confirmation needed
## Confirmed transition rules
## Excluded transitions
## TL;DR
## Implementation scope
## Requirements by screen
## Loading / empty / error / disabled / success states
## Existing-code integration points
## Acceptance criteria
```

## Figma asset collection

- Parse `file_key` and `node-id` from Figma links.
- Before an OAuth flow, confirm the target Figma account, requested scopes, and redirect URL. Use least privilege: normally `file_content:read`; add `current_user:read` only when needed.
- Validate OAuth `state`, use PKCE when possible, and exchange the short-lived authorization code promptly.
- Never expose access tokens, refresh tokens, or client secrets in logs, documents, commits, or PRs.
- Inspect files or nodes before export. Export frames with `GET /v1/images/:key?ids=...`; prefer SVG for suitable icons and vectors; discover image fills through `GET /v1/files/:key/images`.
- Store files according to repository conventions and record node ID, filename, location, export options, and acquisition status in the documents.
- Do not recreate inaccessible assets without asking the user; report the asset, reason it is unavailable, and whether a temporary substitute is acceptable.

## Implementation gate and rules

A request such as “organize it and implement it” is not approval. After presenting the document, require a later message that identifies its path or version and explicitly approves implementation; then hand off through `start-implementation-thread` rather than editing code in the current task.

Before implementation, inspect existing routing, components, and design-system conventions. Use retrieved original visual assets where available; implement only code-reproducible shapes, layout, and typography directly. Include natural loading, empty, error, disabled, and success states, retain app patterns, and visually verify frontend changes when a local browser target is available.
