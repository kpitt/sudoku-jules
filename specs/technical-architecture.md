# Technical Architecture Spec

The Sudoku solver must be built with a focus on high performance, memory efficiency, and long-term maintainability. This involves adopting data-oriented design principles, minimizing heap allocations, and using efficient bitwise operations for set manipulations.

1. **Performance and Memory Efficiency**
    - **User Story**: As a systems engineer, I want the solver to use minimal resources and high-performance algorithms, so that it can run efficiently on low-power hardware or in high-throughput environments.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The application shall prioritize stack allocation over heap allocation wherever possible.
        2. **Ubiquitous**: The system shall use bitmasks (e.g., `uint16`) to represent candidate sets and perform O(1) set operations.
        3. **Ubiquitous**: The data model shall use flattened arrays (e.g., `[81]Cell`) to improve cache locality.

2. **Extensibility and Maintainability**
    - **User Story**: As a maintainer, I want a modular architecture for solving techniques, so that I can easily add new human-like techniques as they are discovered or requested.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The system shall provide a standard interface for deductive techniques.
        2. **Ubiquitous**: The codebase shall maintain a clear separation between the puzzle data model, the solving logic, and utility bitsets.
        3. **Ubiquitous**: The solver shall be designed to allow easy addition of new techniques without modifying the core solving loop.

3. **Correctness and Reliability**
    - **User Story**: As a developer, I want to ensure that performance optimizations do not compromise the correctness of the solver, so that the results are always reliable.
    - **Acceptance Criteria**:
        1. **Ubiquitous**: The system shall include a comprehensive test suite covering all difficulty levels.
        2. **Ubiquitous**: The solver shall maintain a "zero regression" policy for both performance and logic.
