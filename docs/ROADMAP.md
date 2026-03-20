# Sudoku Solver Roadmap

This document tracks the implementation status of various Sudoku solving techniques,
classified by human difficulty.

## Status Legend

- [x] Implemented
- [ ] Todo / Planned
- [-] Stubbed (Defined but not implemented)

## 3. Advanced Techniques

*Complex patterns and chaining logic.*

- [x] **Simple Coloring (Color Chain)**: Single-digit conjugate chains (A/B logic).
- [-] **Multi-Coloring**: Multi-digit conjugate chains.

## 4. Expert Techniques

*Generalized set theory and long chains.*

- [x] **AIC (Alternating Inference Chains)**: Generalized chaining (Strong/Weak links).
- [x] **Sue de Coq**: Set intersection logic.

## 6. Performance Optimizations

*Transitioning from object-oriented/pointer-heavy logic to high-performance data-oriented design.*

### Phase 2: Tables & Lookups

- [ ] **Pre-computed Peers**: Implement `PeerLookup [81][20]uint8` to replace dynamic peer calculation.
- [ ] **Pre-computed Houses**: Implement `HouseLookup [27][9]uint8` to replace `House` structs and maps.
- [ ] **Remove Maps**: Eliminate `House.Unsolved` map in favor of direct array iteration.

### Phase 3: Algorithm Optimization

- [ ] **Bitwise Techniques**: Rewrite all solver techniques to use bitwise logic (POPCNT, AND, OR, XOR).
- [ ] **Branchless Iteration**: Update loops to use fixed-size array iteration.
- [ ] **Stack Allocation**: Ensure the main `Solver` and `Board` structs fit entirely on the stack.
