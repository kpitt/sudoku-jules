// Package puzzle provides the core data structures and logic for representing
// and manipulating a Sudoku puzzle.
package puzzle

import "github.com/kpitt/sudoku/internal/bitset"

const allDigitBits = 0b1111111110

// Cell represents a single square in a Sudoku puzzle.
type Cell struct {
	Row, Col int
	IsGiven  bool

	value      int
	Candidates bitset.BitSet16
}

// NewCell initializes a new cell at the given row and column.
func NewCell(r, c int) Cell {
	return Cell{
		Row: r, Col: c,
		Candidates: bitset.BitSet16(allDigitBits),
	}
}

// IsSolved returns true if a solved value has been placed in this cell.
func (c *Cell) IsSolved() bool {
	return c.value > 0
}

// Index returns the 0-indexed position of the cell in a flat 81-cell array.
func (c *Cell) Index() int {
	return c.Row*9 + c.Col
}

// Value returns the value currently placed in the cell, or 0 if unsolved.
func (c *Cell) Value() int {
	return c.value
}

// PlaceValue places a solved value into the cell, clearing any remaining
// candidates.
func (c *Cell) PlaceValue(val int) {
	c.value = val
	c.Candidates.Clear()
}

// GivenValue places an initial value into the cell, marking it as a given
// value that cannot be changed.  This is used for the initial puzzle setup.
func (c *Cell) GivenValue(val int) {
	c.IsGiven = true
	c.PlaceValue(val)
}

// NumCandidates returns the number of potential values that could be placed in
// this cell.
func (c *Cell) NumCandidates() int {
	return c.Candidates.Size()
}

// CandidateValues returns a slice of all potential values that could be placed
// in this cell.
func (c *Cell) CandidateValues() []int {
	return c.Candidates.Values()
}

// HasCandidate returns true if the specified value is a potential candidate for
// this cell.
func (c *Cell) HasCandidate(val int) bool {
	return c.Candidates.Contains(val)
}

// RemoveCandidate removes the specified value from the cell's potential candidates.
func (c *Cell) RemoveCandidate(val int) {
	c.Candidates.Remove(val)
}

// Box returns the index of the 3x3 box that contains this cell.  Boxes are
// numbered from left-to-right and top-to-bottom, with box 0 at the top-left
// and box 8 at the bottom-right.
func (c *Cell) Box() int {
	return c.Row/3*3 + c.Col/3
}
