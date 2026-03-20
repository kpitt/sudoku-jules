# Sudoku-Jules Implementation Plan

This document outlines the prioritized implementation plan for the Sudoku-Jules project, aiming for a high-performance, expert-level deductive solver ("Ultra" goal).

## Priority 1: CLI & Input Infrastructure
Foundational tasks to enable advanced testing, flexible usage, and compatibility with standard Sudoku libraries like Hodoku.

- **CLI Enhancements**
    - [x] Refactor `cmd/sudoku/main.go` to use the standard `flag` package for robust argument parsing. (Ref: specs/input-handling.md 1.1)
    - [x] Support operational flags: `-file` (input file), `--brute-force` (toggle fallback), `--log-level` (info, debug, trace), and `--test` (regression mode). (Ref: specs/input-handling.md 1.2, 1.3; specs/output-rendering.md 2.1; specs/testing-regression.md 1.2)
- **Hodoku Library Format Support**
    - [x] Update `internal/puzzle/reader.go` to support the Hodoku Library Format: `:technique:candidates:givens:deleted:eliminations:placements:extra`. (Ref: specs/input-handling.md 3.1)
    - [x] Correctly distinguish between givens and placed values (prefixed with `+`) during parsing. (Ref: specs/input-handling.md 3.2)
    - [x] Implement support for the `<deleted candidates>` field to initialize board states with specific eliminations for testing. (Ref: specs/input-handling.md 3.3)
- **Test Setup Utility**
    - [x] Add a helper function `FromHodokuString(s string)` in `internal/puzzle` to simplify puzzle initialization in unit tests and the regression runner. (Ref: specs/testing-regression.md 3.1, 3.2)

## Priority 2: Regression Testing & Validation
Ensuring solver correctness through automated verification and improving the robustness of the fallback engine.

- **Regression Test Runner**
    - [x] Implement a dedicated regression runner to verify deductive techniques against real-world test cases. (Ref: specs/testing-regression.md 1.1)
    - [x] Add verification logic to ensure the solver performs the *exact* eliminations or placements specified in a Hodoku test case. (Ref: specs/testing-regression.md 2.2, 2.3)
    - [x] Report detailed summary statistics including total tests, pass/fail counts, and execution time. (Ref: specs/testing-regression.md 1.3)
- **Early Non-Uniqueness Detection**
    - [ ] Integrate `CountSolutions(2)` into the main `Solve` loop in `internal/solver/solver.go` to detect and report non-unique puzzles early. (Ref: specs/puzzle-solving.md 3.2)

## Priority 3: Expert-Level Deductive Techniques (Hard)
Expanding the deductive engine to cover standard advanced human-like strategies.

- [x] **Full House**: Implement as a distinct technique with higher priority than Hidden Single. (Ref: specs/deductive-techniques.md 1.1)
- [ ] **Remote Pair**: Implement bivalue cell chain detection and elimination. (Ref: specs/deductive-techniques.md 1.2)
- [ ] **BUG + 1**: Implement detection for this uniqueness-based pattern. (Ref: specs/deductive-techniques.md 1.3)
- [ ] **Empty Rectangle**: Implement intersection-based eliminations. (Ref: specs/deductive-techniques.md 1.4)
- [ ] **W-Wing**: Implement eliminations based on linked bivalue cells. (Ref: specs/deductive-techniques.md 1.5)
- [ ] **Uniqueness Tests (2-6)**: Implement Type 2 through Type 6 and **Hidden Rectangle**. (Ref: specs/deductive-techniques.md 1.6)
- [ ] **Finned/Sashimi X-Wing**: Implement Finned and Sashimi variants of X-Wing. (Ref: specs/deductive-techniques.md 1.7)
- [ ] **Avoidable Rectangle**: Implement Type 1 and Type 2. (Ref: docs/REQUIREMENTS.md 2.1)

## Priority 4: Expert-Level Deductive Techniques (Unfair & Extreme)
The "Ultra" goal of implementing advanced sets and complex chain logic to solve the most difficult puzzles.

- [ ] **Finned/Sashimi Swordfish & Jellyfish**: Extend fish logic to complex variants. (Ref: specs/deductive-techniques.md 2.1)
- [ ] **Chains (X-Chain, XY-Chain)**: Implement robust search for single-candidate and bivalue chains. (Ref: specs/deductive-techniques.md 2.2)
- [ ] **Nice Loops (Standard & Grouped)**: Implement closed-loop detection and inference. (Ref: specs/deductive-techniques.md 2.3)
- [ ] **ALS (Almost Locked Sets)**: Implement ALS-XZ and ALS-XY-Wing. (Ref: specs/deductive-techniques.md 2.4)

## Priority 5: Polish & Performance
Optimizing the solver for high throughput and improving the quality of explanation output.

- [ ] **Enhanced Step Logging**: Improve `internal/solver/solution.go` to provide descriptive, human-readable output for complex techniques. (Ref: specs/output-rendering.md 2.2, 2.3)
- [ ] **Performance Optimization**: Profile the solver and optimize hot paths in bitset operations and house lookups. (Ref: specs/technical-architecture.md 1.1, 1.2)
