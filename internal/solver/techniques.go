package solver

import (
	"github.com/kpitt/sudoku/internal/puzzle"
)

type techniqueKind int

const (
	// Enum constants for each solving technique, which must be defined in the
	// order that they should be applied.
	kindNakedSingle techniqueKind = iota
	kindHiddenSingle
	kindLockedCandidatesPointing
	kindLockedCandidatesClaiming
	kindNakedPair
	kindNakedTriple
	kindHiddenPair
	kindHiddenTriple
	kindNakedQuadruple
	kindHiddenQuadruple
	kindXWing
	kindSwordfish
	kindJellyfish
	kindRemotePair
	kindBUG
	kindSkyscraper
	kindTwoStringKite
	kindEmptyRectangle
	kindWWing
	kindXYWing
	kindXYZWing
	kindAvoidableRectangle
	kindUniqueRectangle1
	kindUniqueRectangle2
	kindUniqueRectangle3
	kindUniqueRectangle4
	kindHiddenRectangle
	kindFinnedXWing
	kindFinnedSwordfish
	kindFinnedJellyfish
	kindSueDeCoq
	kindSimpleColoring
	kindMultiColoring
	kindXChain
	kindXYChain
	kindAIC
	kindALS
	kindBruteForce
)

// A Technique represents a single Sudoku solving technique, represented by a
// display name and a function to check for the technique in the puzzle.
type Technique struct {
	Name  string
	Check func() bool
}

// initTechniques initializes a list of the known solving techniques and
// assigns it to the solver.
// The order must match the order of the techniqueKind constants.
// TODO: Implement Bitwise Techniques: Rewrite all solver techniques to use bitwise logic (POPCNT, AND, OR, XOR).
func (s *Solver) initTechniques() {
	// The order of this list must match the order of the techniqueKind constants.
	s.techniques = []Technique{
		{"Naked Single", nil}, // checked during candidate removal
		{"Hidden Single", s.findHiddenSingles},
		{"Locked Candidates Type 1 (Pointing)", s.findPointingTuples},
		{"Locked Candidates Type 2 (Claiming)", s.findClaimingTuples},
		{"Naked Pair", s.findNakedPairs},
		{"Naked Triple", s.findNakedTriples},
		{"Hidden Pair", s.findHiddenPairs},
		{"Hidden Triple", s.findHiddenTriples},
		{"Naked Quadruple", s.findNakedQuadruples},
		{"Hidden Quadruple", s.findHiddenQuadruples},
		{"X-Wing", s.findXWings},
		{"Swordfish", s.findSwordfish},
		{"Jellyfish", s.findJellyfish},
		{"Remote Pair", nil}, // TODO: Implement "Remote Pair" technique
		{"BUG+1", nil},       // TODO: Implement "BUG+1" technique
		{"Skyscraper", s.findSkyscraper},
		{"2-String Kite", s.findTwoStringKite},
		{"Empty Rectangle", nil}, // TODO: Implement "Empty Rectangle" technique
		{"W-Wing", nil},          // TODO: Implement "W-Wing" technique
		{"XY-Wing", s.findXYWings},
		{"XYZ-Wing", s.findXYZWings},
		{"Avoidable Rectangle", s.findAvoidableRectangles}, // TODO: Implement "Avoidable Rectangle" technique
		{"Unique Rectangle Type 1", s.findUniqueRectangleType1},
		{"Unique Rectangle Type 2", nil}, // TODO: Implement "Unique Rectangle Type 2" technique
		{"Unique Rectangle Type 3", nil}, // TODO: Implement "Unique Rectangle Type 3" technique
		{"Unique Rectangle Type 4", nil}, // TODO: Implement "Unique Rectangle Type 4" technique
		{"Hidden Rectangle", nil},        // TODO: Implement "Hidden Rectangle" technique
		{"Finned X-Wing", nil},           // TODO: Implement "Finned X-Wing" technique
		{"Finned Swordfish", nil},        // TODO: Implement "Finned Swordfish" technique
		{"Finned Jellyfish", nil},        // TODO: Implement "Finned Jellyfish" technique
		{"Sue de Coq", nil},              // TODO: Implement "Sue de Coq" technique
		{"Simple Coloring", nil},         // TODO: Implement "Simple Coloring" technique
		{"Multi-Coloring", nil},          // TODO: Implement "Multi-Coloring" technique
		{"X-Chain", nil},                 // TODO: Implement "X-Chain" technique
		{"XY-Chain", nil},                // TODO: Implement "XY-Chain" technique
		{"AIC", nil},                     // TODO: Implement "AIC" technique
		{"ALS", nil},                     // TODO: Implement "ALS" technique
		{"Brute Force", nil},             // custom check as last resort
	}
}

// ***** IMPORTANT NOTE *****
//
// When processing a check against a set of viable candidates, _always_
// short-circuit the remaining checks after making a change that could
// invalidate the remaining candidates.  Checks should not be combined in a
// single pass unless the checks are completely independent.  If it isn't
// clear whether or not the checks are independent, go ahead and short-ciruit,
// and any additional candidates will get checked in the next solver pass.
//
// For checks that should _not_ short-circuit, be careful when using OR
// expressions to combine the results.  The safest approach is to use an
// accumulator variable and the pattern `found = check(...) || found` for each
// check.  The `||` operator will short-circuit after the first term that
// evalues to `true`, so only the first term is guaranteed to be evaluated.

// eliminateFromOtherLocs removes the candidates listed in values from all
// cells that are not listed in locs.
func (s *Solver) eliminateFromOtherLocs(
	h *House, values ValSet, locs LocSet, step *SolutionStep,
) bool {
	found := false
	for l := range 9 {
		if locs.Contains(l) {
			continue
		}
		c := h.Cells[l]

		// Check for intersection first to avoid inner loop/allocations.
		common := c.Candidates.Intersection(values)
		if common.Empty() {
			continue
		}

		for v := range common.All() {
			step.DeleteCandidate(c.Index(), v)
			found = true
		}
	}

	return found
}

// eliminateFromOtherLocsMulti removes the candidates listed in values from all
// cells from each house in houses whose index is not listed in locs.  Returns
// true if at least one candidate was eliminated.
func (s *Solver) eliminateFromOtherLocsMulti(
	houses []*House, values ValSet, locs LocSet, step *SolutionStep,
) bool {
	updated := false
	for _, g := range houses {
		updated = s.eliminateFromOtherLocs(g, values, locs, step) || updated
	}

	return updated
}

// eliminateOtherValues removes candidates that are not listed in values from
// the cells in locs.
func (s *Solver) eliminateOtherValues(
	h *House, values ValSet, locs LocSet, step *SolutionStep,
) bool {
	found := false
	for l := range locs.All() {
		c := h.Cells[l]

		// We want to remove any candidates in `c` that are NOT in `values`.
		toRemove := c.Candidates.Difference(values)
		if toRemove.Empty() {
			continue
		}

		for v := range toRemove.All() {
			step.DeleteCandidate(c.Index(), v)
			found = true
		}
	}

	return found
}

// eliminateFromIntersection removes `val` from any cell that sees ALL provided indices.
// If pivotIdx is -1, it finds intersection of (idx1, idx2).
// If pivotIdx is >= 0, it finds intersection of (idx1, idx2, pivotIdx).
func (s *Solver) eliminateFromIntersection(
	idx1, idx2, pivotIdx, val int, step *SolutionStep,
) bool {
	changed := false

	peers1 := puzzle.GetPeers(idx1)

	// Build fast lookup table for peers of idx2.
	isPeerOf2 := s.getPeerSet(idx2)

	var isPeerOfPivot [81]bool
	hasPivot := false
	if pivotIdx != -1 {
		isPeerOfPivot = s.getPeerSet(pivotIdx)
		hasPivot = true
	}

	for i := range 20 {
		candidateIdx := int(peers1[i])

		// Must be peer of idx2
		if !isPeerOf2[candidateIdx] {
			continue
		}

		// Must be peer of pivot
		if hasPivot && !isPeerOfPivot[candidateIdx] {
			continue
		}

		// Don't eliminate from the pattern cells themselves
		if candidateIdx == idx1 || candidateIdx == idx2 || candidateIdx == pivotIdx {
			continue
		}

		// Remove the candidate if it exists
		if s.puzzle.Cell(candidateIdx).HasCandidate(val) {
			changed = true
			step.DeleteCandidate(candidateIdx, val)
		}
	}
	return changed
}

// getPeerSet returns a fast lookup table for peers.
func (s *Solver) getPeerSet(idx int) [81]bool {
	peers := puzzle.GetPeers(idx)

	var lookup [81]bool
	for i := range 20 {
		peerIdx := peers[i]
		lookup[peerIdx] = true
	}
	return lookup
}

// eliminateValuesFromCell removes all candidates listed in values from the cell
// at (r,c).
func (s *Solver) eliminateValuesFromCell(
	r, c int, values ValSet, step *SolutionStep,
) bool {
	cell := s.puzzle.Get(r, c)
	found := false
	for _, v := range values.Values() {
		if cell.HasCandidate(v) {
			step.DeleteCandidate(r*9+c, v)
			found = true
		}
	}
	return found
}
