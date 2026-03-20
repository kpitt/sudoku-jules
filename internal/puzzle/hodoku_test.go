package puzzle

import (
	"testing"
)

func TestFromHodokuString(t *testing.T) {
	// Simple Full House example
	s := ":0000:1:+6+3+42+8+9...8+7531+6+2+9+42+91+7+5+4+6+38+349+16+58+7+2+5+8+79+42..+6+1+62+83+754+97+5+8+4+2+39+61+9.+6+5784+23+4+2+3+6+91...:::182:"
	p, err := FromHodokuString(s)
	if err != nil {
		t.Fatalf("FromHodokuString failed: %v", err)
	}

	// Verify some cells
	// r1c1 should be 6, and NOT a given (placed with +)
	c1 := p.Get(0, 0)
	if !c1.IsSolved() || c1.Value() != 6 {
		t.Errorf("r1c1: expected solved 6, got %v", c1.Value())
	}
	if c1.IsGiven {
		t.Errorf("r1c1: expected NOT given, but is marked as given")
	}

	// r1c4 should be 2, and a given
	c4 := p.Get(0, 3)
	if !c4.IsSolved() || c4.Value() != 2 {
		t.Errorf("r1c4: expected solved 2, got %v", c4.Value())
	}
	if !c4.IsGiven {
		t.Errorf("r1c4: expected given, but is NOT marked as given")
	}

	// r1c7 should be empty
	c7 := p.Get(0, 6)
	if c7.IsSolved() {
		t.Errorf("r1c7: expected unsolved, got %v", c7.Value())
	}
}

func TestFromHodokuString_DeletedCandidates(t *testing.T) {
	// Hodoku string with deleted candidates
	// Digit 1 at r1c1 is a given.
	// Digit 2 at r1c2 is empty, but we'll delete candidate 3 from it.
	// :technique:candidates:givens:deleted:eliminations:placements:extra
	s := ":0000:1:1................................................................................:312:::112:"
	// Wait, format for deleted is <candidate><line><col>. So 312 means digit 3 at r1c2.
	
	p, err := FromHodokuString(s)
	if err != nil {
		t.Fatalf("FromHodokuString failed: %v", err)
	}

	// r1c2 (index 1) should NOT have candidate 3
	c2 := p.Get(0, 1)
	if c2.HasCandidate(3) {
		t.Errorf("r1c2: expected candidate 3 to be deleted")
	}
	// It should still have other candidates, e.g., 4
	if !c2.HasCandidate(4) {
		t.Errorf("r1c2: expected candidate 4 to be present")
	}
}
