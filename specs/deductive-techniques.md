# Deductive Solving Techniques Spec

This document provides detailed descriptions and logical requirements for the deductive solving techniques to be implemented in the Sudoku-Jules solver. These techniques are categorized by their difficulty level as defined in `REQUIREMENTS.md`.

## 1. Priority 3: Hard Techniques

### 1.1 Full House
- **Logic**: When a house (row, column, or box) has exactly one empty cell, that cell must contain the only remaining candidate for that house.
- **Requirement**: Implement as a distinct technique with higher priority than Naked/Hidden Singles for educational clarity.
- **Verification**: The solver must identify the single missing value and place it.

### 1.2 Remote Pair
- **Logic**: A chain of bivalue cells (cells with exactly two candidates) where each cell contains the same two candidates (e.g., {x, y}). If the chain has an even number of links, the start and end cells see each other, any cell that sees both can have candidates x and y eliminated.
- **Verification**: Detect chains of identical bivalue cells and eliminate candidates from their common peers.

### 1.3 BUG + 1 (Bivalue Universal Grave)
- **Logic**: A state where every unsolved cell has exactly two candidates, except for one cell which has three candidates, and every candidate appears exactly twice in every house. To avoid a non-unique solution (the BUG), the cell with three candidates must be set to the candidate that appears three times in its row, column, and box.
- **Verification**: Identify the BUG+1 pattern and place the correct value in the tri-value cell.

### 1.4 Empty Rectangle
- **Logic**: Uses a box where a candidate is restricted to a "conjugate pair" in a row and column within that box (forming an L-shape or similar). If this box-structure is linked to another strong link for the same candidate in a different house, an elimination can be made.
- **Verification**: Identify intersection-based eliminations for a single candidate.

### 1.5 W-Wing
- **Logic**: Two bivalue cells with the same candidates {x, y} that are linked by a strong link of one of the candidates (e.g., x) in another house. If a cell sees both bivalue cells, the other candidate (y) can be eliminated from it.
- **Verification**: Identify linked bivalue cells and perform eliminations.

### 1.6 Uniqueness Tests (Types 2-6)
- **Logic**: Building on Type 1, these tests use the property that a valid Sudoku must have a unique solution to eliminate candidates that would lead to a "Deadly Pattern" (a rectangle of 4 cells across 2 rows, 2 columns, and 2 boxes where candidates are interchangeable).
    - **Type 2**: One extra candidate in two of the cells.
    - **Type 3**: Two or more extra candidates in two of the cells, forming a virtual naked subset.
    - **Type 4**: A conjugate pair in the rectangle for one of the candidates.
    - **Types 5/6**: More complex variations involving multiple extra candidates or links.
- **Verification**: Correctly identify the specific type and perform the corresponding elimination.

### 1.7 Finned/Sashimi X-Wing
- **Logic**: An X-Wing that is nearly perfect except for one or more "finned" cells in one of the base lines.
    - **Finned**: If the fins are in the same box as a cover-line intersection, eliminations can be made in that box.
    - **Sashimi**: A finned fish where the fish would be incomplete without the fin.
- **Verification**: Identify the base/cover lines and the fin, then eliminate from appropriate peers.

## 2. Priority 4: Unfair & Extreme Techniques

### 2.1 Finned/Sashimi Swordfish & Jellyfish
- **Logic**: Extension of Finned/Sashimi logic to 3x3 (Swordfish) and 4x4 (Jellyfish) structures.
- **Verification**: Identify complex fish structures with fins and perform eliminations.

### 2.2 Chains (X-Chain, XY-Chain)
- **Logic**:
    - **X-Chain**: A chain of strong and weak links for a single candidate. If the chain starts and ends with a strong link, the candidate can be eliminated from any cell that sees both ends.
    - **XY-Chain**: A chain of bivalue cells. If the chain starts with candidate 'a' and ends with candidate 'a', then 'a' can be eliminated from any cell seeing both ends.
- **Verification**: Robust search for chains and valid eliminations.

### 2.3 Nice Loops (Standard & Grouped)
- **Logic**: A sequence of links (strong or weak) that forms a closed loop. Depending on the types of links, candidates can be placed or eliminated.
    - **Grouped**: Links can involve groups of cells (like a candidate in a box-row intersection).
- **Verification**: Detect closed loops and apply the appropriate inference.

### 2.4 ALS-XZ and ALS-XY-Wing
- **Logic**:
    - **ALS (Almost Locked Set)**: A set of N cells with N+1 candidates.
    - **ALS-XZ**: Two ALSs linked by a "restricted common" candidate. Depending on the connection, other common candidates can be eliminated.
    - **ALS-XY-Wing**: Three ALSs linked in a wing-like structure.
- **Verification**: Identify ALS sets and their connections to perform complex eliminations.
