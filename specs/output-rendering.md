# Output Rendering Spec

The Sudoku application must provide clear and informative feedback to the user via the console. This includes rendering the grid with visual cues for different types of values and optionally logging the detailed steps taken by the deductive engine to arrive at the solution.

1. **Console Grid Rendering**
    - **User Story**: As a player, I want a well-formatted grid in my console, so that I can easily read the puzzle and its solution.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The system shall render the 9x9 puzzle grid and solution steps to the console.
        2. **Where** the output is a TTY console, **then** the application shall use ANSI color codes to distinguish between givens, placed values, and pencil marks.
        3. **Ubiquitous**: The system shall use distinct visual markers or colors for givens (initial values), placed values (solved by logic), and pencil marks (candidate values).

2. **Solution Logging**
    - **User Story**: As a student of Sudoku, I want to see the reasoning behind each step, so that I can learn how to apply different techniques myself.
    - **Acceptance Criteria**:
        1. **Optional**: WHERE the user enables step logging, **then** the system shall log each deductive step taken.
        2. **Ubiquitous**: The step logs shall include the technique used (e.g., "X-Wing"), the target cell(s), and the candidates that were eliminated.
        3. **Ubiquitous**: The log output must be readable and sequential to follow the solver's progress.
