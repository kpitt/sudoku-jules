package puzzle

import (
	"reflect"
	"testing"
)

func TestNewCell(t *testing.T) {
	c := NewCell(2, 5)

	if c.Row != 2 {
		t.Errorf("expected row 2, got %d", c.Row)
	}
	if c.Col != 5 {
		t.Errorf("expected col 5, got %d", c.Col)
	}
	if c.IsGiven {
		t.Error("new cell should not be given")
	}
	if c.IsSolved() {
		t.Error("new cell should not be solved")
	}
	if c.NumCandidates() != 9 {
		t.Errorf("new cell should have 9 candidates, got %d", c.NumCandidates())
	}
}

func TestCell_IsSolved(t *testing.T) {
	tests := []struct {
		name     string
		cell     Cell
		expected bool
	}{
		{
			name:     "new cell",
			cell:     NewCell(0, 0),
			expected: false,
		},
		{
			name: "solved cell",
			cell: func() Cell {
				c := NewCell(0, 0)
				c.PlaceValue(5)
				return c
			}(),
			expected: true,
		},
		{
			name: "given cell",
			cell: func() Cell {
				c := NewCell(0, 0)
				c.GivenValue(5)
				return c
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cell.IsSolved(); got != tt.expected {
				t.Errorf("Cell.IsSolved() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCell_Index(t *testing.T) {
	tests := []struct {
		name     string
		r, c     int
		expected int
	}{
		{"top left", 0, 0, 0},
		{"top right", 0, 8, 8},
		{"middle", 4, 4, 40},
		{"bottom left", 8, 0, 72},
		{"bottom right", 8, 8, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := NewCell(tt.r, tt.c)
			if got := cell.Index(); got != tt.expected {
				t.Errorf("Cell.Index() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCell_Value(t *testing.T) {
	c := NewCell(0, 0)
	if got := c.Value(); got != 0 {
		t.Errorf("new cell should have value 0, got %d", got)
	}

	c.PlaceValue(5)
	if got := c.Value(); got != 5 {
		t.Errorf("cell should have value 5, got %d", got)
	}
}

func TestCell_PlaceValue(t *testing.T) {
	c := NewCell(0, 0)
	c.PlaceValue(7)

	if !c.IsSolved() {
		t.Error("cell should be solved after PlaceValue")
	}
	if c.Value() != 7 {
		t.Errorf("expected value 7, got %d", c.Value())
	}
	if c.NumCandidates() != 0 {
		t.Errorf("solved cell should have 0 candidates, got %d", c.NumCandidates())
	}
	if c.IsGiven {
		t.Error("placed value should not be marked as given")
	}
}

func TestCell_GivenValue(t *testing.T) {
	c := NewCell(0, 0)
	c.GivenValue(3)

	if !c.IsSolved() {
		t.Error("cell should be solved after GivenValue")
	}
	if c.Value() != 3 {
		t.Errorf("expected value 3, got %d", c.Value())
	}
	if c.NumCandidates() != 0 {
		t.Errorf("given cell should have 0 candidates, got %d", c.NumCandidates())
	}
	if !c.IsGiven {
		t.Error("cell should be marked as given after GivenValue")
	}
}

func TestCell_Candidates(t *testing.T) {
	c := NewCell(0, 0)

	// Test initial candidates
	if got := c.NumCandidates(); got != 9 {
		t.Errorf("expected 9 candidates, got %d", got)
	}

	expectedInitial := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if got := c.CandidateValues(); !reflect.DeepEqual(got, expectedInitial) {
		t.Errorf("expected initial candidates %v, got %v", expectedInitial, got)
	}

	for i := 1; i <= 9; i++ {
		if !c.HasCandidate(i) {
			t.Errorf("new cell should have candidate %d", i)
		}
	}

	// Test removing candidates
	c.RemoveCandidate(5)
	if c.HasCandidate(5) {
		t.Error("cell should not have candidate 5 after removal")
	}
	if got := c.NumCandidates(); got != 8 {
		t.Errorf("expected 8 candidates after removal, got %d", got)
	}

	expectedAfterRemoval := []int{1, 2, 3, 4, 6, 7, 8, 9}
	if got := c.CandidateValues(); !reflect.DeepEqual(got, expectedAfterRemoval) {
		t.Errorf("expected candidates after removal %v, got %v", expectedAfterRemoval, got)
	}

	// Removing a candidate that doesn't exist should be a no-op
	c.RemoveCandidate(5)
	if got := c.NumCandidates(); got != 8 {
		t.Errorf("expected 8 candidates after redundant removal, got %d", got)
	}
}

func TestCell_Box(t *testing.T) {
	tests := []struct {
		r, c, expected int
	}{
		{0, 0, 0},
		{0, 8, 2},
		{4, 4, 4},
		{8, 0, 6},
		{8, 8, 8},
		{1, 5, 1},
		{7, 3, 7},
	}

	for _, tt := range tests {
		c := NewCell(tt.r, tt.c)
		if got := c.Box(); got != tt.expected {
			t.Errorf("Cell(%d,%d).Box() = %d; want %d", tt.r, tt.c, got, tt.expected)
		}
	}
}
