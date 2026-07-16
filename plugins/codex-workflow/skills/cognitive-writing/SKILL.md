---
name: cognitive-writing
description: Structure plans, PRs, technical documents, reports, and implementation notes to minimize reader cognitive load. Apply when creating or revising written engineering artifacts.
---

# Cognitive Writing

Write for a reader who should understand the decision and verify the work without reconstructing the author's reasoning.

## Language policy

Write the artifact in the language the user requested. If no language was explicitly requested, use the language of the user's request. This includes headings, tables, checklists, examples, and Mermaid labels. Preserve code, commands, API names, identifiers, and required proper names.

## Core rules

1. Start with a TL;DR of at most three lines: purpose, essential change, and verification state.
2. Separate facts, assumptions, decisions, open questions, and excluded alternatives.
3. Organize information in the reader's order of understanding: motivation, structure, decisions, entry points, then verification.
4. Give each change or plan step one compact block containing purpose, files or entry point, key decision, and verification.
5. Explain only non-obvious decisions and their trade-offs. State external constraints with a source when relevant.
6. Name a clear entry point for large changes: the first file, module, or section to inspect.
7. Define uncommon domain terms on first use.
8. Separate mechanical noise—formatting, import sorting, unused-variable removal, renames, and comment-only changes—from behavioral changes.

## Markdown

Use valid GitHub-flavored Markdown. Leave blank lines around headings, lists, and code blocks. Give code fences a language where possible. Use descriptive headings and avoid vague titles such as “Update” or “Cleanup”.

## Visuals

Add a small visual when it materially reduces ambiguity: before/after images for visual UI changes, a short recording for interaction, and a Mermaid diagram for architecture, data flow, state, or multi-step navigation. Put the diagram before the detailed explanation and split complex topics into small diagrams by viewpoint.

## Final check

- Can the reader identify why the work exists and what success means from the first screen?
- Are facts distinct from assumptions and pending decisions?
- Does every non-trivial step identify an entry point and verification?
- Are unrelated changes separated?
- Is the artifact localized to the user-requested language and free of broken Markdown?
