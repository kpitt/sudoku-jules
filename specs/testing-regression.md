# Testing and Regression Spec

To ensure the reliability and correctness of the Sudoku solver, especially when implementing complex deductive techniques, the application must provide dedicated testing modes and regression tools. These features allow for high-confidence updates and ensure that new optimizations do not introduce regressions.

1. **Testing Mode and Regression Runner**
    - **User Story**: As a developer, I want a specialized testing mode and a regression runner, so that I can automatically verify that all implemented techniques work correctly on real-world test cases.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall provide a dedicated testing mode to simplify the verification of deductive techniques.
        2. **Ubiquitous**: The system shall provide a CLI option to run all tests in a specified file, treating each line as a single test case in the Hodoku Library Format.
        3. **When** the regression runner completes, **then** it shall report the total number of tests run, the number of passed tests, and the number of failed tests.
    - **Resouces:** @test/lib/reglib-1.3.txt contains the regression test library for the Hodoku solver. Use this as an example of the expected structure for a "Hodoku Library Format" file.

2. **Deductive Verification Logic**
    - **User Story**: As a tester, I want the system to precisely verify each deductive step, so that I can ensure the techniques are applied correctly and only remove the intended candidates.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall initialize the grid with the specified givens and placed values from a test case.
        2. **Ubiquitous**: The application shall apply candidate deletions defined in the test case before running the solver.
        3. **Ubiquitous**: The system shall verify that the deductive engine performs the EXACT eliminations or placements defined in the test case.
        4. **Unwanted Behavior**: IF the deductive engine makes an incorrect elimination or fails to perform a specified one, THEN the test case shall fail.

3. **Unit Test Setup Utility**
    - **User Story**: As a developer writing unit tests, I want a simple utility to set up complex puzzle states, so that I can write concise and readable tests for specific techniques.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall provide a library function to initialize a puzzle grid from a single string, including those in Hodoku format.
        2. **Ubiquitous**: This utility shall be available for use within the internal test suite to simplify test setups.
