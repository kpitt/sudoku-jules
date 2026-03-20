# Puzzle Solving Spec

The Sudoku solver is designed to solve standard 9x9 puzzles by prioritizing human-like deductive reasoning techniques to provide speed and educational value. If a puzzle cannot be solved using these human-like techniques, the system must employ a high-performance brute-force algorithm (Dancing Links) to guarantee a solution.

1. **Standard 9x9 Sudoku Support**
    - **User Story**: As a player, I want the solver to support standard 9x9 grids, so that I can solve typical Sudoku puzzles.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The system shall solve standard 9x9 Sudoku puzzles.
        2. **Unwanted Behavior**: IF the input grid is not a valid 9x9 Sudoku (e.g., incorrect dimensions or rule violations), THEN the system shall return an error message.

2. **Deductive Engine and Priority**
    - **User Story**: As a player who wants to learn, I want the solver to use human-like techniques in a specific priority order, so that I can follow the simplest possible logic for each step.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The solver shall implement human-like solving techniques as defined in the REQUIREMENTS.md document (Easy, Medium, Hard, Unfair, Extreme).
        2. **Ubiquitous**: The solver shall apply techniques in the specified priority order (e.g., Full House before Naked Single).
        3. **When** a deductive technique successfully identifies an elimination or placement, **then** the solver shall restart the technique search from the highest priority.

3. **Brute-Force Fallback (Algorithm X)**
    - **User Story**: As a user with an extremely difficult puzzle, I want the solver to fall back to a robust algorithm when deductive reasoning fails, so that I always get a solution.
    - **Acceptance Criteria**:
        1. **When** the deductive engine can no longer make progress on a puzzle, **then** the system shall employ Knuth's Algorithm X (using Dancing Links) to find the solution.
        2. **Ubiquitous**: The brute-force algorithm shall find all possible solutions (if the puzzle is not unique) or determine if no solution exists.

4. **Performance Targets**
    - **User Story**: As a high-frequency user, I want the solver to be extremely fast, so that I can solve thousands of puzzles per second.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The solver shall aim for sub-millisecond solve times for easy to advanced puzzles using deductive techniques.
        2. **Ubiquitous**: The Dancing Links implementation shall be optimized for high performance on puzzles that cannot be solved deductively.
