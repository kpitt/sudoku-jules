0a. Study `specs/*` using parallel subagents to learn the application specifications.
0b. Study @IMPLEMENTATION_PLAN.md.
0c. For reference, the application source code is in `cmd/*` and library package source code is in `internal/*`.

1. Your task is to implement functionality per the specifications using parallel subagents. Follow @IMPLEMENTATION_PLAN.md and choose the most important item to address. Before making changes, search the codebase (don't assume not implemented) using subagents. Use parallel subagents for searches/reads but only 1 subagent for build/tests. Use subagents when complex reasoning is needed (debugging, architectural decisions).
2. After implementing functionality or resolving problems, run the tests for that unit of code that was improved. If functionality is missing then it's your job to add it as per the application specifications. Ultrathink.
3. When you discover issues, immediately update @IMPLEMENTATION_PLAN.md with your findings using a subagent. When resolved, update and remove the item.
4. When the tests pass, update @IMPLEMENTATION_PLAN.md, then `git add -A` then `git commit` with a message describing the changes. After the commit, `git push`.

**Operational Guidelines:**

- You may add extra logging if required to debug issues.
- Keep @IMPLEMENTATION_PLAN.md current with learnings using a subagent — future work depends on this to avoid duplicating efforts. Update especially after finishing your turn.
- When @IMPLEMENTATION_PLAN.md becomes large periodically clean out the items that are completed from the file using a subagent.
- Implement functionality completely. Placeholders and stubs waste efforts and time redoing the same work.
- For any bugs you notice, resolve them or document them in @IMPLEMENTATION_PLAN.md using a subagent even if it is unrelated to the current piece of work.
- If you find inconsistencies in the `specs/*` then use a subagent with 'ultrathink' requested to update the specs.
- When you learn something new about how to run the application, update @AGENTS.md using a subagent but keep it brief. For example if you run commands multiple times before learning the correct command then that file should be updated.
- IMPORTANT: Keep @AGENTS.md operational only — status updates and progress notes belong in `IMPLEMENTATION_PLAN.md`. A bloated AGENTS.md pollutes every future loop's context.
