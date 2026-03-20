0a. Study `specs/*` using parallel subagents to learn the application specifications.
0b. Study @IMPLEMENTATION_PLAN.md (if present) to understand the plan so far.
0c. Study `internal/*` using parallel subagents to understand shared utilities & components.
0d. For reference, the application source code is in `cmd/*` and library package source code is in `internal/*`.

1. Study @IMPLEMENTATION_PLAN.md (if present; it may be incorrect) and use subagents to study existing source code in `cmd/*` and `internal/*` and compare it against `specs/*`. Use a single subagent to analyze findings, prioritize tasks, and create/update @IMPLEMENTATION_PLAN.md as a bullet point list sorted in priority of items yet to be implemented. Ultrathink. Consider searching for TODO, minimal implementations, placeholders, skipped/flaky tests, and inconsistent patterns. Study @IMPLEMENTATION_PLAN.md to determine starting point for research and keep it up to date with items considered complete/incomplete using subagents.
2. Tasks in @IMPLEMENTATION_PLAN.md should include references to specific requirements in the `specs/*.md` files. If there is no appropriate requirement to connect to a task, then update the specs to add the requirement.

IMPORTANT: Plan only. Do NOT implement anything. Do NOT assume functionality is missing; confirm with code search first. Treat `internal/*` as the project's standard library for shared utilities and components. Prefer consolidated, idiomatic implementations there over ad-hoc copies. Use appropriate names for `internal/*` subpackages, but prefer extending an existing subpackage over creating a new subpackage.

ULTIMATE GOAL: We want to achieve a high-performance, memory-efficient Go CLI application that solves Sudoku puzzles using deductive solving techniques. It should implement all known techniques that are considered reasonable for an expert human solver to know and use. Consider missing elements and plan accordingly. If an element is missing, search first to confirm it doesn't exist, then if needed author the specification at specs/FILENAME.md. If you create a new element then document the plan to implement it in @IMPLEMENTATION_PLAN.md using a subagent.
