# Input Handling Spec

The Sudoku application must be versatile in how it accepts puzzle data, supporting multiple input sources like file paths and standard input (stdin). Additionally, it must handle various common Sudoku representation formats, including grid strings, human-readable ASCII, and advanced library formats like Hodoku.

1. **Input Source Selection**
    - **User Story**: As a CLI user, I want the application to read puzzles from files or stdin, so that I can easily integrate it with other tools or enter puzzles manually.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall accept input from a specified file path.
        2. **Ubiquitous**: The application shall accept input from standard input (stdin).
        3. **When** standard input is an interactive TTY and no input file is provided, **then** the application shall prompt the user to input the puzzle.
        4. **Ubiquitous**: The application shall ignore any line or part of a line starting with a `#` character.

2. **Supported Input Formats**
    - **User Story**: As a developer, I want support for various formats, so that I can use puzzles from different sources and databases without conversion.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall support single 81-character strings using '0' or '.' for empty cells.
        2. **Ubiquitous**: The application shall support 9x9 grids with various delimiters and whitespace.
        3. **Ubiquitous**: The application shall support ASCII-formatted grids with borders (e.g., `+---+---+---+`).
        4. **Ubiquitous**: The application shall support pencil-mark grids showing possible candidates.

3. **Hodoku Library Format Support**
    - **User Story**: As a regression tester, I want the solver to support the Hodoku Library Format, so that I can verify specific solving techniques on pre-configured board states.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall support strings in the format `:technique:candidates:givens:deleted:eliminations:placements:extra`.
        2. **Ubiquitous**: The loader shall distinguish between givens and placed values (prefixed with `+`).
        3. **Ubiquitous**: The loader shall support explicit candidate deletion defined in the `<deleted candidates>` field.
        4. **Unwanted Behavior**: IF a Hodoku string is malformed or has invalid segments, THEN the loader shall return an error.
    - **Resources:** Use the ID column of the techniques table in @docs/TECHNIQUES.md to determine the mapping from the `technique` field values to the corresponding technique implementation function.
