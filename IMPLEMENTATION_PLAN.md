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

- [ ] **Finned/Sashimi Swordfish & Jellyfish**: Extend fish logic to complex variants. (Ref: specs/deductive-techniques.md 2.1)
- [ ] **Chains (X-Chain, XY-Chain)**: Implement robust search for single-candidate and bivalue chains. (Ref: specs/deductive-techniques.md 2.2)
- [ ] **Nice Loops (Standard & Grouped)**: Implement closed-loop detection and inference. (Ref: specs/deductive-techniques.md 2.3)
- [ ] **ALS (Almost Locked Sets)**: Implement ALS-XZ and ALS-XY-Wing. (Ref: specs/deductive-techniques.md 2.4)

## Priority 5: Polish & Performance
Optimizing the solver for high throughput and improving the quality of explanation output.

- [ ] **Enhanced Step Logging**: Improve `internal/solver/solution.go` to provide descriptive, human-readable output for complex techniques. (Ref: specs/output-rendering.md 2.2, 2.3)
- [ ] **Performance Optimization**: Profile the solver and optimize hot paths in bitset operations and house lookups. (Ref: specs/technical-architecture.md 1.1, 1.2)
