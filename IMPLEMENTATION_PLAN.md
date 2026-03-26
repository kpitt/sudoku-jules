# Sudoku-Jules Implementation Plan

This document outlines the prioritized implementation plan for the Sudoku-Jules project, aiming for a high-performance, expert-level deductive solver ("Ultra" goal).

## Priority 3: Expert-Level Deductive Techniques (Hard)
Expanding the deductive engine to cover standard advanced human-like strategies.

- [x] **Full House**: Implement as a distinct technique with higher priority than Hidden Single. (Ref: specs/deductive-techniques.md 1.1)
- [x] **Remote Pair**: Implement bivalue cell chain detection and elimination. (Ref: specs/deductive-techniques.md 1.2)
- [x] **BUG + 1**: Implement detection for this uniqueness-based pattern. (Ref: specs/deductive-techniques.md 1.3)
- [x] **Empty Rectangle**: Implement intersection-based eliminations. (Ref: specs/deductive-techniques.md 1.4)
- [x] **W-Wing**: Implement eliminations based on linked bivalue cells. (Ref: specs/deductive-techniques.md 1.5)
- [x] **Unique Rectangle Type 2** (Ref: specs/deductive-techniques.md 1.6)
- [x] **Uniqueness Tests (3-6)**: Implement Type 3 through Type 6 and **Hidden Rectangle**. (Ref: specs/deductive-techniques.md 1.6)
- [x] **Finned/Sashimi X-Wing**: Implement Finned and Sashimi variants of X-Wing. (Ref: specs/deductive-techniques.md 1.7)
- [x] **Avoidable Rectangle (Type 1)** (Ref: docs/REQUIREMENTS.md 2.1)
- [x] **Avoidable Rectangle (Type 2)**: Implement Type 2. (Ref: docs/REQUIREMENTS.md 2.1)

## Priority 4: Expert-Level Deductive Techniques (Unfair & Extreme)
The "Ultra" goal of implementing advanced sets and complex chain logic to solve the most difficult puzzles.

- [x] **Finned/Sashimi Swordfish & Jellyfish**: Extend fish logic to complex variants. (Ref: specs/deductive-techniques.md 2.1)
- [x] **Chains (X-Chain, XY-Chain)**: Implement robust search for single-candidate and bivalue chains. (Ref: specs/deductive-techniques.md 2.2)
- [x] **Nice Loops (Standard & Grouped)**: Implement closed-loop detection and inference. (Ref: specs/deductive-techniques.md 2.3)
- [x] **ALS (Almost Locked Sets)**: Implement ALS-XZ and ALS-XY-Wing. (Ref: specs/deductive-techniques.md 2.4)

## Priority 5: Polish & Performance
Optimizing the solver for high throughput and improving the quality of explanation output.

- [x] **Enhanced Step Logging**: Improve `internal/solver/solution.go` to provide descriptive, human-readable output for complex techniques. (Ref: specs/output-rendering.md 2.2, 2.3)
- [x] **Performance Optimization**: Profile the solver and optimize hot paths in bitset operations and house lookups. (Ref: specs/technical-architecture.md 1.1, 1.2). **Note: Achieved 2x-4x performance speedup across benchmark suites.**

## Priority 6: Remaining Expert/Extreme Techniques
The final frontier of deductive techniques as outlined in the project roadmap and requirements.

- [x] **Simple Coloring**: Single-digit conjugate chains (A/B logic). (Ref: docs/REQUIREMENTS.md 2.1, docs/ROADMAP.md 3)
- [x] **Multi-Coloring**: Multi-digit conjugate chains, implemented using the generalized AIC engine. (Ref: docs/REQUIREMENTS.md 2.1, docs/ROADMAP.md 3)
- [x] **Sue de Coq**: Implement complex set intersection logic. (Ref: docs/REQUIREMENTS.md 2.1, docs/ROADMAP.md 4)
- [x] **ALS-XY-Chain**: Generalized Almost Locked Set chains, implemented using ALS nodes within the AIC engine. (Ref: docs/REQUIREMENTS.md 2.1)
- [x] **Franken & Mutant Fish**: Implement fish logic using non-standard houses (rows/columns + boxes), including Finned and Sashimi variants. (Ref: docs/REQUIREMENTS.md 2.1)
- [x] **Forcing Chains & Nets**: Complex multi-digit branching chains. (Ref: docs/REQUIREMENTS.md 2.1)

## Priority 7: Code Quality & Standards
Ensuring codebase health and maintainability through rigorous static analysis and refactoring.

- [x] **Address Linting Issues**: Resolve the 98 issues identified by `golangci-lint`, including:
    - [x] **Unchecked Errors**: Address `errcheck` failures in `main.go`, `print.go`, and `regression.go`.
    - [x] **Error Handling**: Modernize error logic using `errors.As` and `%w` formatting (`errorlint`).
    - [x] **Code Idioms**: Apply `gocritic` recommendations such as `paramTypeCombine`.

## Findings & Next Steps

### Recent Findings
- **Technique Coverage**: The solver now supports most "Unfair" and "Extreme" techniques, including ALS-XZ, ALS-XY-Wing, and a robust AIC (Alternating Inference Chain) engine.
- **Forcing Chains & Nets**: Implemented comprehensive Forcing Net (Cell and House types) and Forcing Chain (Contradiction type) detection using the AIC engine.
- **AIC Engine Extensions**: The AIC engine has been extended to support **Multi-Coloring** (via single-digit conjugate chains) and **ALS-XY-Chain** (by integrating ALS nodes as chain links).
- **Simple Coloring**: Implementation includes grouped strong links, providing more elimination power than basic coloring.
- **Sue de Coq**: Successfully implemented with box-row and box-column intersection logic.
- **Franken & Mutant Fish**: Implemented generalized fish logic (X-Wing, Swordfish, Jellyfish) using non-standard houses (rows/columns + boxes), including Finned and Sashimi variants.
- **Performance**: High-performance bitset operations and efficient house lookups maintain the 2x-4x speedup over the baseline solver.
- **Code Quality**: All 98+ linting issues have been resolved; the codebase now follows modern Go idioms and has robust error handling.

### Next Steps
- **Performance Optimization (Phase 2)**: Transition to data-oriented design with pre-computed peer and house tables as outlined in `docs/ROADMAP.md`.
- **Formal Verification**: Add more regression cases for the newest techniques (Franken/Mutant Fish, Multi-Coloring, ALS-XY-Chain, Sue de Coq, Forcing Nets).

