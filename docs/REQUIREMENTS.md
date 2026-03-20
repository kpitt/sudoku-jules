# Sudoku Solver Requirements

This document outlines the functional and non-functional requirements for the Sudoku solver project.

## 1. Project Overview

The goal of this project is to develop a high-performance Sudoku solver that employs a hybrid approach: prioritizing human-like deductive reasoning for speed and educational value, and falling back to a robust brute-force algorithm (Dancing Links) when necessary.

## 2. Functional Requirements

### 2.1 Puzzle Solving
- **9x9 Sudoku**: The system must solve standard 9x9 Sudoku puzzles.
- **Deductive Engine**: The solver MUST implement all of the following human-like solving techniques, and they must be applied in the specified priority order.
    - **Easy**:
        1. Full House
        2. Naked Single
        3. Hidden Single
    - **Medium**:
        4. Locked Pair
        5. Locked Triple
        6. Locked Candidates Type 1 (Pointing)
        7. Locked Candidates Type 2 (Claiming)
        8. Naked Pair
        9. Naked Triple
        10. Hidden Pair
        11. Hidden Triple
    - **Hard**:
        12. Naked Quadruple
        13. Hidden Quadruple
        14. X-Wing
        15. Swordfish
        16. Jellyfish
        17. Remote Pair
        18. BUG + 1
        19. Skyscraper
        20. Two String Kite
        21. Turbot Fish
        22. Empty Rectangle
        23. W-Wing
        24. XY-Wing
        25. XYZ-Wing
        26. Uniqueness Test 1
        27. Uniqueness Test 2
        28. Uniqueness Test 3
        29. Uniqueness Test 4
        30. Uniqueness Test 5
        31. Uniqueness Test 6
        32. Hidden Rectangle
        33. Avoidable Rectangle Type 1
        34. Avoidable Rectangle Type 2
        35. Finned X-Wing
        36. Sashimi X-Wing
        37. Simple Colors
        38. Multi Colors
    - **Unfair**:
        39. Finned Swordfish
        40. Sashimi Swordfish
        41. Finned Jellyfish
        42. Sashimi Jellyfish
        43. Sue de Coq
        44. X-Chain
        45. XY-Chain
        46. Nice Loop
        47. Grouped Nice Loop
        48. ALS-XZ
        49. ALS-XY-Wing
        50. ALS-XY-Chain
        51. Franken X-Wing
        52. Franken Swordfish
        53. Finned Franken X-Wing
        54. Finned Franken Swordfish
    - **Extreme**:
        55. Forcing Chain
        56. Forcing Net
        57. Brute Force (as a deductive step)
        58. Give Up (termination state)
- **Brute-Force Fallback**: If deductive techniques fail to solve the puzzle, the system must employ Knuth's Algorithm X (using Dancing Links) to find a solution.
- **Validation**: The system must be able to validate whether a given grid is a valid Sudoku solution.

### 2.2 Input and Output
- **Input Source**: The application must accept input from either a specified file path or from the standard input (stdin).
- **Line Comments**: When reading from a file or stdin, any line or part of a line starting with a `#` character must be ignored.
- **Interactive Input**: If standard input is an interactive TTY and no input file is provided, the application must prompt the user to input the puzzle.
- **Supported Input Formats**: The system must support several common Sudoku input formats:
    - **Single 81-character string**: A compact format using '0' or '.' for empty cells (e.g., `0000000000000...`).
    - **9x9 grid**: 9 lines of 9 characters each, with whitespace or other delimiters allowed.
    - **ASCII-formatted grids**: Common human-readable formats with grid borders (e.g., `+---+---+---+`).
    - **Pencil-mark grids**: Grids that include possible candidates for empty cells, often used for partially solved puzzles.
    - **Hodoku Library Format**: Support strings in the "HoDoKu Library Format" (`:technique:candidates:givens:deleted:eliminations:placements:extra`). This format is essential for regression testing and must support:
        - **Givens vs. Placed Values**: Distinguishing between real givens and values already placed (prefixed with `+`).
        - **Explicit Candidate Deletion**: Supporting the removal of specific candidates as defined in the `<deleted candidates>` field to set up precise board states.
- **Console Output**: Render the puzzle grid and solution steps to the console.
- **TTY Color Support**: When outputting a puzzle grid to a TTY console, use ANSI color codes to distinguish specific elements:
    - **Given values**: Values provided in the initial puzzle state.
    - **Placed values**: Values determined by the solver or user.
    - **Pencil mark values**: Candidate values for empty cells.
- **Solution Logging**: Optionally log each deductive step taken, including the technique used and the candidates eliminated.

### 2.3 Testing and Regression
- **Testing Mode**: The application must provide a dedicated testing mode to simplify the verification of specific deductive techniques.
- **Regression Test Runner**: Provide a CLI option to run all tests in a specified file. Each line in the file should be treated as a single test case in the Hodoku Library Format. Each test run should report the total number of tests run, the number of passed tests, and the number of failed tests.
- **Verification Logic**: For each regression test, the application must:
    - Initialize the grid with the specified givens and placed values.
    - Apply the specified candidate deletions.
    - Verify that the deductive engine performs the exact eliminations or placements defined in the test case.
- **Unit Test Setup Utility**: The application must provide a library function to initialize a puzzle grid from a single string (including Hodoku format), allowing for concise setup in unit tests.

## 3. Non-Functional Requirements

### 3.1 Performance
- **High-Performance Algorithms**: Prioritize algorithms that minimize computational complexity.
- **Sub-millisecond Solving**: Aim for sub-millisecond solve times for easy to advanced puzzles using deductive techniques.
- **Efficient Brute-Force**: The Dancing Links implementation should be highly optimized for puzzles that cannot be solved deductively.

### 3.2 Memory and Resource Efficiency
- **Minimize Heap Allocations**: Use data structures and algorithms that prefer stack allocation and minimize garbage collection overhead.
- **Data-Oriented Design**: Use flattened arrays (e.g., `[81]Cell`) and contiguous memory layouts to improve cache locality.
- **Bitwise Operations**: Use bitmasks (e.g., `uint16`) for representing candidate sets and performing set operations (Union, Intersection, etc.) in O(1) time.

### 3.3 Extensibility and Maintainability
- **Modular Technique Interface**: New deductive techniques should be easy to add by implementing a standard interface.
- **Clear Architecture**: Maintain a clean separation between the data model (`puzzle`), the solving logic (`solver`), and utility structures (`bitset`).

### 3.4 Correctness and Reliability
- **Comprehensive Testing**: Maintain a test suite with puzzles of varying difficulty (Beginner, Expert, "Impossible") and specific technique-testing puzzles.
- **Zero Regressions**: Ensure that performance optimizations do not compromise the correctness of the solver.
