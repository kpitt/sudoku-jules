package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindSkyscraper(t *testing.T) {
	// Create an empty puzzle and solver
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	// We'll test with digit 1.
	val := 1

	// Setup Skyscraper pattern on Rows.
	// Base Lines: Row 1 and Row 4.
	// Common Column: Col 1.
	// Tops: (1, 7) and (4, 8).
	// Target to eliminate: (3, 7).
	// (3, 7) sees Top 1 (1, 7) via Column 7.
	// (3, 7) sees Top 2 (4, 8) via Box 5 (both in Box 5).

	// Helper to set candidate
	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Prepare state:
	// Clear all candidates for val 1 first
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	// Also clear nums for cols/boxes for consistency
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// Set up the pattern
	setCandidate(1, 1) // Base 1, Common
	setCandidate(1, 7) // Base 1, Top 1

	setCandidate(4, 1) // Base 2, Common
	setCandidate(4, 8) // Base 2, Top 2

	// Set up target
	setCandidate(3, 7) // Target (Row 3, Col 7)

	// Run Skyscraper
	found := s.findSkyscraper()
	if !found {
		t.Errorf("Skyscraper technique should have been found")
	}

	// Verify eliminations
	if p.Get(3, 7).HasCandidate(val) {
		t.Errorf("Target (3,7) should have eliminated candidate %d", val)
	}
}

func TestRemotePair(t *testing.T) {
	// Example from reglib-1.3.txt
	s := ":0703:58:..+1.+4+9+2+636+3.21+7+9+4.942+63..+7.2634.+17.98.+4.+9..2...+9.+6+2.+34..7..4.9.+4..9+7631..+9+6.+2.+4.7::577 579 594 596 877 879 894 896::7"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, &Options{EnableDebug: true})
	solver.processInitialValues()

	found := solver.findRemotePairs()
	if !found {
		t.Error("Remote Pair should be found")
	}

	// The Hodoku string says it should eliminate candidate 5 and 8 from 
	// several cells.
	// Eliminations: 577 579 594 596 877 879 894 896
	// Note: 577 means value 5 at r7c7 (r6c6). Wait, Hodoku indices are 1-based digit,row,col.
	// "577" => value 5 at r7c7.
	// "877" => value 8 at r7c7.
	
	// r7c7 is index (7-1)*9 + (7-1) = 6*9+6 = 60.
	if p.Cell(60).HasCandidate(5) || p.Cell(60).HasCandidate(8) {
		t.Error("Candidate 5 or 8 should be eliminated from r7c7")
	}
}

func TestWWing(t *testing.T) {
	// W-Wing example from reglib-1.3.txt
	s := ":0803:14:6..+9+5..7...+9.2.....58.+31...+1+64+3+8+9+7+52...1+7+59+46597+24+6..892+54+1+76+8+3...5+6+2.....68+93...::417 427 437 489 499::"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, &Options{EnableDebug: true})
	solver.processInitialValues()

	found := solver.findWWing()
	if !found {
		t.Error("W-Wing should be found")
	}

	// Eliminations: 417 427 437 489 499
	// 417 => digit 4 at r1c7 (index 6)
	if p.Cell(6).HasCandidate(4) {
		t.Error("Candidate 4 should be eliminated from r1c7")
	}
}

func TestEmptyRectangle(t *testing.T) {
	// Empty Rectangle example from reglib-1.3.txt
	s := ":0402:9:7+2+4+956+1381+6842+3+5+9+7+9+3+5+7+1+8+6+2+45..3..+8+1..4..8+17+5..+81.+7.24..+13....+7+2...1...+85.5...7.6+1::986::"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, &Options{EnableDebug: true})
	solver.processInitialValues()

	found := solver.findEmptyRectangle()
	if !found {
		t.Error("Empty Rectangle should be found")
	}

	// Elimination: 986 => digit 9 at r8c6 (index 68)
	if p.Cell(68).HasCandidate(9) {
		t.Error("Candidate 9 should be eliminated from r8c6")
	}
}

func TestXChain(t *testing.T) {
	// Example from reglib-1.3.txt
	s := ":0701:7:3.4+52..8...6.+9.....5..7.3.....68+9.2+3...+734....6+315+27...1.+9+6......9.+4..6.+6.8217..5::742::5"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, &Options{EnableDebug: true})
	solver.processInitialValues()

	found := solver.findXChains()
	if !found {
		t.Error("X-Chain should be found")
	}

	// Elimination: 742 => digit 7 at r4c2
	idx := (4-1)*9 + (2-1)
	if p.Cell(idx).HasCandidate(7) {
		t.Error("Candidate 7 should be eliminated from r4c2")
	}
}

func TestXYChain(t *testing.T) {
	// Example from reglib-1.3.txt
	s := ":0702:7:76+2+8+1+3+59+4...+7+6+9+1+2+8.+9.42536+7.+765+9+82.+39..1+32......+67+4+8.9..+7+9+5643+1.49+38+1..2.1.+2+47+9+8.::758 787::19"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, nil)
	solver.processInitialValues()

	found := false
	for {
		if !solver.findXYChains() {
			break
		}
		found = true
		// Check if we already eliminated what we wanted
		idx1 := (5-1)*9 + (8-1)
		idx2 := (8-1)*9 + (7-1)
		if !p.Cell(idx1).HasCandidate(7) && !p.Cell(idx2).HasCandidate(7) {
			break
		}
	}

	if !found {
		t.Error("XY-Chain should be found")
	}

	// Eliminations: 758 and 787
	idx1 := (5-1)*9 + (8-1)
	idx2 := (8-1)*9 + (7-1)
	if p.Cell(idx1).HasCandidate(7) || p.Cell(idx2).HasCandidate(7) {
		t.Error("Candidate 7 should be eliminated from r5c8 and r8c7")
	}
}

func TestUniqueRectangleType2(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// UR Type 2 Setup: {1, 2} with extra 3
	// Corners: (0,0), (0,3), (1,0), (1,3)
	// (0,0) and (1,0) are in Box 0.
	// (0,3) and (1,3) are in Box 1.
	
	// Bivalue cells
	setCandidate(0, 0, 1); setCandidate(0, 0, 2)
	setCandidate(0, 3, 1); setCandidate(0, 3, 2)
	
	// Extra cells (sharing same row or column)
	// Let's put them in column 0: (0,0) and (1,0)? No, those are not corners.
	// Corners are (0,0), (0,3), (1,0), (1,3).
	// Extra cells must be (0,0) and (0,3) [same row] or (0,0) and (1,0) [same col].
	
	// Reset and set correctly
	for r := 0; r < 2; r++ {
		for c := 0; c < 4; c++ {
			p.Get(r, c).Candidates.Clear()
		}
	}
	
	// Bivalue cells: (1,0) and (1,3)
	setCandidate(1, 0, 1); setCandidate(1, 0, 2)
	setCandidate(1, 3, 1); setCandidate(1, 3, 2)
	
	// Extra cells: (0,0) and (0,3)
	setCandidate(0, 0, 1); setCandidate(0, 0, 2); setCandidate(0, 0, 3)
	setCandidate(0, 3, 1); setCandidate(0, 3, 2); setCandidate(0, 3, 3)
	
	// Target cell: must see (0,0) and (0,3)
	// (0,1) is in the same row as both.
	setCandidate(0, 1, 3)

	found := s.findUniqueRectangleType2()
	if !found {
		t.Error("Unique Rectangle Type 2 should be found")
	}

	if p.Get(0, 1).HasCandidate(3) {
		t.Error("Candidate 3 should be eliminated from (0,1)")
	}
	}

	func TestAvoidableRectangle(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// AR Type 1 Setup: {1, 2}
	// Corners: (0,0), (0,3), (1,0), (1,3)
	// (0,0) and (1,0) in Box 0.
	// (0,3) and (1,3) in Box 1.
	
	// Solved corners (non-givens)
	// Corner (0,0) = 1
	p.Get(0, 0).PlaceValue(1)
	// Corner (0,3) = 2
	p.Get(0, 3).PlaceValue(2)
	// Corner (1,0) = 2
	p.Get(1, 0).PlaceValue(2)
	
	// Unsolved corner (1,3) has {1, 2, 3}
	// If (1,3) is 1, we have a deadly pattern {1, 2, 2, 1}.
	setCandidate(1, 3, 1)
	setCandidate(1, 3, 2)
	setCandidate(1, 3, 3)

	found := s.findAvoidableRectangles()
	if !found {
		t.Error("Avoidable Rectangle Type 1 should be found")
	}

	if p.Get(1, 3).HasCandidate(1) {
		t.Error("Candidate 1 should be eliminated from (1,3)")
	}
}

	func TestFindTwoStringKite(t *testing.T) {

	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 1

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Clear candidates for val 1
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	setCandidate(0, 0)
	setCandidate(0, 4)
	setCandidate(1, 2)
	setCandidate(5, 2)
	setCandidate(5, 4) // Target

	found := s.findTwoStringKite()
	if !found {
		t.Errorf("2-String Kite should have been found")
	}

	if p.Get(5, 4).HasCandidate(val) {
		t.Errorf("Candidate 1 should be eliminated from (5,4)")
	}
}

func TestFindXWing(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 1

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Clear candidates for val 1
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// X-Wing Setup (Rows as Base)
	// R1C1, R1C8
	// R4C1, R4C8
	setCandidate(1, 1)
	setCandidate(1, 8)
	setCandidate(4, 1)
	setCandidate(4, 8)

	// Target to eliminate (Col 1, not in base rows)
	setCandidate(0, 1)
	// Add noise to Row 0 so it has > 2 candidates, preventing it from being
	// selected as a base line for X-Wing (Size 2).
	setCandidate(0, 2)
	setCandidate(0, 3)

	found := s.findXWings()
	if !found {
		t.Errorf("X-Wing technique should have been found")
	}

	if p.Get(0, 1).HasCandidate(val) {
		t.Errorf("Target (0,1) should have eliminated candidate %d", val)
	}
}

func TestFindSwordfish(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 2

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Clear candidates for val 2
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// Swordfish Setup (Rows as Base)
	// R1: C1, C4
	// R4: C4, C7
	// R7: C1, C7
	setCandidate(1, 1)
	setCandidate(1, 4)
	setCandidate(4, 4)
	setCandidate(4, 7)
	setCandidate(7, 1)
	setCandidate(7, 7)

	// Target to eliminate (Col 1, Row 0)
	setCandidate(0, 1)
	// Add noise to Row 0 so it has > 3 candidates, preventing it from being
	// selected as a base line for Swordfish (Size 3).
	setCandidate(0, 2)
	setCandidate(0, 3)
	setCandidate(0, 5)

	found := s.findSwordfish()
	if !found {
		t.Errorf("Swordfish technique should have been found")
	}

	if p.Get(0, 1).HasCandidate(val) {
		t.Errorf("Target (0,1) should have eliminated candidate %d", val)
	}
}

func TestFindJellyfish(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 3

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Clear candidates for val 3
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// Jellyfish Setup (Rows as Base)
	// R1: C1, C2, C4, C5
	// R2: C1, C2, C4, C5
	// R4: C1, C2, C4, C5
	// R5: C1, C2, C4, C5
	rows := []int{1, 2, 4, 5}
	cols := []int{1, 2, 4, 5}

	for _, r := range rows {
		for _, c := range cols {
			setCandidate(r, c)
		}
	}

	// Target to eliminate (Col 1, Row 0)
	setCandidate(0, 1)
	// Add noise to Row 0 so it has > 4 candidates, preventing it from being
	// selected as a base line for Jellyfish (Size 4).
	setCandidate(0, 2)
	setCandidate(0, 3)
	setCandidate(0, 6)
	setCandidate(0, 7)
	setCandidate(0, 8)

	found := s.findJellyfish()
	if !found {
		t.Errorf("Jellyfish technique should have been found")
	}

	if p.Get(0, 1).HasCandidate(val) {
		t.Errorf("Target (0,1) should have eliminated candidate %d", val)
	}
}

func TestFindXYWing(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	// Clear all candidates
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// XY-Wing Setup
	// Pivot: (4,4) with candidates {1, 2}
	setCandidate(4, 4, 1)
	setCandidate(4, 4, 2)

	// Pincer 1: (4,0) with candidates {1, 3} (Same Row)
	setCandidate(4, 0, 1)
	setCandidate(4, 0, 3)

	// Pincer 2: (0,4) with candidates {2, 3} (Same Col)
	setCandidate(0, 4, 2)
	setCandidate(0, 4, 3)

	// Target: (0,0) with candidate {3} (Sees both pincers)
	setCandidate(0, 0, 3)

	found := s.findXYWings()
	if !found {
		t.Errorf("XY-Wing technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(3) {
		t.Errorf("Target (0,0) should have eliminated candidate 3")
	}
}

func TestFindXYZWing(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	// Clear all candidates
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// XYZ-Wing Setup
	// Pivot: (4,4) with candidates {1, 2, 3}
	setCandidate(4, 4, 1)
	setCandidate(4, 4, 2)
	setCandidate(4, 4, 3)

	// Pincer 1: (4,5) with candidates {1, 3} (Same Row & Box)
	setCandidate(4, 5, 1)
	setCandidate(4, 5, 3)

	// Pincer 2: (2,4) with candidates {2, 3} (Same Col, DIFFERENT Box)
	// (2,4) is in Box 1 (Row 2, Col 4)
	setCandidate(2, 4, 2)
	setCandidate(2, 4, 3)

	// Target: (5,4) with candidate {3} (Sees all three)
	// Sees (4,4) via Col 4 / Box 4
	// Sees (4,5) via Box 4
	// Sees (2,4) via Col 4
	setCandidate(5, 4, 3)

	found := s.findXYZWings()
	if !found {
		t.Errorf("XYZ-Wing technique should have been found")
	}

	if p.Get(5, 4).HasCandidate(3) {
		t.Errorf("Target (5,4) should have eliminated candidate 3")
	}
}

func TestFindNakedPairs(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	// Clear all candidates
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Naked Pair in Row 0: Cells (0,0) and (0,1) with candidates {1, 2}
	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 1)
	setCandidate(0, 1, 2)

	// Other cells in Row 0 that have candidates 1 and 2
	setCandidate(0, 2, 1) // Target to eliminate
	setCandidate(0, 3, 2) // Target to eliminate

	// Add other candidates so these cells aren't empty
	setCandidate(0, 2, 3)
	setCandidate(0, 3, 4)

	found := s.findNakedPairs()
	if !found {
		t.Errorf("Naked Pair technique should have been found")
	}

	if p.Get(0, 2).HasCandidate(1) {
		t.Errorf("Target (0,2) should have eliminated candidate 1")
	}
	if p.Get(0, 3).HasCandidate(2) {
		t.Errorf("Target (0,3) should have eliminated candidate 2")
	}
}

func TestFindNakedTriples(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Naked Triple: {1,2}, {2,3}, {1,3}
	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 1)
	setCandidate(0, 2, 3)

	setCandidate(0, 4, 1) // Target
	setCandidate(0, 4, 2) // Target
	setCandidate(0, 4, 4) // Noise

	found := s.findNakedTriples()
	if !found {
		t.Errorf("Naked Triple technique should have been found")
	}

	if p.Get(0, 4).HasCandidate(1) || p.Get(0, 4).HasCandidate(2) {
		t.Errorf("Target (0,4) should have eliminated candidates 1 and 2")
	}
}

func TestFindNakedQuadruples(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 3)
	setCandidate(0, 2, 4)
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 4)

	setCandidate(0, 5, 1) // Target
	setCandidate(0, 5, 5) // Noise

	found := s.findNakedQuadruples()
	if !found {
		t.Errorf("Naked Quadruple technique should have been found")
	}

	if p.Get(0, 5).HasCandidate(1) {
		t.Errorf("Target (0,5) should have eliminated candidate 1")
	}
}

func TestFindHiddenPairs(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 1)
	setCandidate(0, 1, 2)

	setCandidate(0, 0, 3) // Target
	setCandidate(0, 1, 4) // Target

	setCandidate(0, 2, 3)
	setCandidate(0, 3, 4)
	setCandidate(0, 4, 5)

	found := s.findHiddenPairs()
	if !found {
		t.Errorf("Hidden Pair technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(3) {
		t.Errorf("Target (0,0) should have eliminated candidate 3")
	}
	if p.Get(0, 1).HasCandidate(4) {
		t.Errorf("Target (0,1) should have eliminated candidate 4")
	}
}

func TestFindHiddenTriples(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 1)
	setCandidate(0, 2, 3)

	setCandidate(0, 0, 4) // Target
	setCandidate(0, 1, 5) // Target
	setCandidate(0, 2, 6) // Target

	setCandidate(0, 3, 4)
	setCandidate(0, 3, 5)
	setCandidate(0, 4, 5)
	setCandidate(0, 4, 6)

	found := s.findHiddenTriples()
	if !found {
		t.Errorf("Hidden Triple technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(4) {
		t.Errorf("Target (0,0) should have eliminated candidate 4")
	}
	if p.Get(0, 1).HasCandidate(5) {
		t.Errorf("Target (0,1) should have eliminated candidate 5")
	}
	if p.Get(0, 2).HasCandidate(6) {
		t.Errorf("Target (0,2) should have eliminated candidate 6")
	}
}

func TestFindHiddenQuadruples(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 3)
	setCandidate(0, 2, 4)
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 4)

	setCandidate(0, 0, 5) // Target
	setCandidate(0, 1, 6) // Target

	setCandidate(0, 4, 5)
	setCandidate(0, 4, 6)
	setCandidate(0, 5, 7)
	setCandidate(0, 5, 8)

	found := s.findHiddenQuadruples()
	if !found {
		t.Errorf("Hidden Quadruple technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(5) {
		t.Errorf("Target (0,0) should have eliminated candidate 5")
	}
	if p.Get(0, 1).HasCandidate(6) {
		t.Errorf("Target (0,1) should have eliminated candidate 6")
	}
}

func TestFindHiddenSingles(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Setup: Candidate 1 is only in (0,0) for Row 0
	setCandidate(0, 0, 1)

	// Add other candidates to (0,0) so it's not a Naked Single
	setCandidate(0, 0, 2)
	setCandidate(0, 0, 3)

	// Other cells in Row 0 get candidates 2 and 3, but NOT 1
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 2)
	setCandidate(0, 2, 3)

	found := s.findHiddenSingles()
	if !found {
		t.Errorf("Hidden Single technique should have been found")
	}

	// Because finding a Hidden Single places the value, the cell should be
	// marked as solved with the value 1.
	if !p.Get(0, 0).IsSolved() {
		t.Errorf("Hidden Single should have placed a value in cell (0,0)")
	} else if p.Get(0, 0).Value() != 1 {
		t.Errorf("Hidden Single should have placed value 1 in cell (0,0)")
	}
}

func TestFindPointingTuples(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Setup: Candidate 1 in Box 0 is restricted to Row 0
	setCandidate(0, 0, 1)
	setCandidate(0, 1, 1)

	// Add noise to avoid singles
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 3)

	// Target to eliminate in Row 0, outside Box 0
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 4) // Noise

	found := s.findPointingTuples()
	if !found {
		t.Errorf("Pointing Tuples technique should have been found")
	}

	if p.Get(0, 3).HasCandidate(1) {
		t.Errorf("Target (0,3) should have eliminated candidate 1")
	}
}

func TestFindClaimingTuples(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Setup: Candidate 1 in Row 0 is restricted to Box 0
	setCandidate(0, 0, 1)
	setCandidate(0, 1, 1)

	// Add noise to avoid singles
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 3)

	// Target to eliminate in Box 0, outside Row 0
	setCandidate(1, 0, 1)
	setCandidate(1, 0, 4) // Noise

	found := s.findClaimingTuples()
	if !found {
		t.Errorf("Claiming Tuples technique should have been found")
	}

	if p.Get(1, 0).HasCandidate(1) {
		t.Errorf("Target (1,0) should have eliminated candidate 1")
	}
}

func TestFindUniqueRectangleType1(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Setup: Unique Rectangle (UR) Type 1
	// The 4 corners must be in exactly 2 rows, 2 columns, and 2 boxes.
	// For Type 1, three corners have exactly candidates {X, Y}.
	// The fourth corner has candidates {X, Y, Z}, we can eliminate X and Y.
	// We'll use Cells (0,0), (0,3), (1,0), and (1,3).
	// Row 0 and Row 1 (2 rows).
	// Col 0 and Col 3 (2 columns).
	// Boxes 0 and 1 (2 boxes).

	// Bivalue cells: {1, 2}
	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2) // Base
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 2) // Row Wing
	setCandidate(1, 0, 1)
	setCandidate(1, 0, 2) // Col Wing

	// Fourth corner with extra candidates
	setCandidate(1, 3, 1) // Target to eliminate
	setCandidate(1, 3, 2) // Target to eliminate
	setCandidate(1, 3, 3)
	setCandidate(1, 3, 4)

	found := s.findUniqueRectangleType1()
	if !found {
		t.Errorf("Unique Rectangle Type 1 technique should have been found")
	}

	if p.Get(1, 3).HasCandidate(1) || p.Get(1, 3).HasCandidate(2) {
		t.Errorf("Target (1,3) should have eliminated candidates 1 and 2")
	}
}
