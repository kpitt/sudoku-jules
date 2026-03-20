package puzzle

import (
	"fmt"

	"github.com/kpitt/sudoku/internal/bitset"
)

// Puzzle represents a 9x9 Sudoku grid and its current state.
type Puzzle struct {
	Cells [81]Cell

	// Holds counts of how many of each digit still needs to be placed.  If the
	// count for a digit reaches 0, then that digit is completely solved.
	// Index 0 holds the total count of unsolved grid cells.  When this value
	// reaches 0, the puzzle is completely solved.
	unsolvedCounts [10]int
}

// NewPuzzle creates and initializes a new, empty 9x9 Sudoku puzzle.
func NewPuzzle() *Puzzle {
	p := &Puzzle{}
	for r := range 9 {
		for c := range 9 {
			p.Cells[r*9+c] = NewCell(r, c)
		}
	}

	for digit := range 10 {
		if digit == 0 {
			// Digit 0 represents the total count of unsolved cells.
			p.unsolvedCounts[digit] = 9 * 9
		} else {
			p.unsolvedCounts[digit] = 9
		}
	}

	return p
}

// Get returns the cell at the specified row and column.
func (p *Puzzle) Get(r, c int) *Cell {
	return &p.Cells[r*9+c]
}

// Cell returns the cell at the specified 0-indexed position.
func (p *Puzzle) Cell(idx int) *Cell {
	return &p.Cells[idx]
}

// IsSolved returns true if all cells in the puzzle have been correctly filled.
func (p *Puzzle) IsSolved() bool {
	return p.unsolvedCounts[0] == 0
}

// IsDigitSolved returns true if all 9 instances of the specified digit have
// been placed in the puzzle.
func (p *Puzzle) IsDigitSolved(digit int) bool {
	return p.unsolvedCounts[digit] == 0
}

// GivenValue sets the initial value for a cell, marking it as a "given" that
// cannot be changed.
func (p *Puzzle) GivenValue(idx, val int) {
	p.Cell(idx).GivenValue(val)
	p.updatePuzzleState(idx, val)
}

// PlaceValue sets the value for a cell during the solving process.
func (p *Puzzle) PlaceValue(idx, val int) bool {
	cell := p.Cell(idx)
	if cell.IsSolved() {
		puzzleStateError(fmt.Sprintf("cell %s is already solved (value=%d)",
			FormatCell(idx), cell.Value()))
		return false
	}

	cell.PlaceValue(val)
	p.updatePuzzleState(idx, val)
	return true
}

// ValidateSolution checks if the current puzzle state is a valid Sudoku solution.
func (p *Puzzle) ValidateSolution() error {
	// Check if all cells are filled
	for i := range 81 {
		if !p.Cells[i].IsSolved() {
			return fmt.Errorf("cell %s is not filled", FormatCell(i))
		}
	}

	// Check row constraints
	for r := range 9 {
		var seen bitset.BitSet16
		for c := range 9 {
			val := p.Cells[r*9+c].Value()
			if val < 1 || val > 9 {
				return fmt.Errorf("invalid value %d in cell r%dc%d", val, r+1, c+1)
			}
			if seen.Contains(val) {
				return fmt.Errorf("duplicate value %d in row %d", val, r+1)
			}
			seen.Add(val)
		}
	}

	// Check column constraints
	for c := range 9 {
		var seen bitset.BitSet16
		for r := range 9 {
			val := p.Cells[r*9+c].Value()
			if seen.Contains(val) {
				return fmt.Errorf("duplicate value %d in column %d", val, c+1)
			}
			seen.Add(val)
		}
	}

	// Check box constraints
	for box := range 9 {
		var seen bitset.BitSet16
		boxRow, boxCol := box/3, box%3
		for i := range 9 {
			r, c := boxRow*3+i/3, boxCol*3+i%3
			val := p.Cells[r*9+c].Value()
			if seen.Contains(val) {
				return fmt.Errorf("duplicate value %d in box %d", val, box+1)
			}
			seen.Add(val)
		}
	}

	return nil
}

// updatePuzzleState updates the valid candidates and unsolved counts after a
// value of val is placed at index idx.
func (p *Puzzle) updatePuzzleState(idx, val int) {
	p.removeConflictingCandidates(idx, val)
	p.updateUnsolvedCounts(idx, val)
}

func (p *Puzzle) removeConflictingCandidates(idx, val int) {
	r, c := idx/9, idx%9
	rb, cb := r/3*3, c/3*3
	for i := range 9 {
		p.Cells[r*9+i].RemoveCandidate(val) // remove candidate from row r
		p.Cells[i*9+c].RemoveCandidate(val) // remove candidate from column c
		// remove candidate from the box that contains (r,c)
		br, bc := rb+i/3, cb+i%3
		p.Cells[br*9+bc].RemoveCandidate(val)
	}
}

func (p *Puzzle) updateUnsolvedCounts(idx, val int) {
	p.unsolvedCounts[0]--
	p.unsolvedCounts[val]--
	if p.unsolvedCounts[val] < 0 {
		puzzleStateError(fmt.Sprintf("too many instances of digit %d when placing cell %s",
			val, FormatCell(idx)))
	}
}
