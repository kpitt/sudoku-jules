package solver

import (
	"github.com/kpitt/sudoku/internal/bitset"
	"github.com/kpitt/sudoku/internal/puzzle"
)

func (s *Solver) findXYWings() bool {
	// Iterate every cell to act as the potential "Pivot"
	for pivotIdx := range 81 {
		pivotCell := s.puzzle.Cell(pivotIdx)

		// Pivot must have exactly 2 candidates
		if pivotCell.NumCandidates() != 2 {
			continue
		}

		// Extract Pivot candidates (X and Y)
		// We can get them from CandidateValues() which returns sorted int slice
		vals := pivotCell.CandidateValues()
		x, y := vals[0], vals[1]

		// Iterate Pivot's peers to find potential Pincers
		peers := puzzle.GetPeers(pivotIdx) // Take a pointer to avoid array copy

		// We need to find two distinct pincers among the peers
		for i := range 20 {
			pincerAIdx := peers[i]
			pincerA := s.puzzle.Cell(pincerAIdx)

			if pincerA.NumCandidates() != 2 {
				continue
			}

			// It must look like {X, Z} or {Y, Z}
			matchX := pincerA.HasCandidate(x)
			matchY := pincerA.HasCandidate(y)

			// XOR: It must match X or Y, but not both (that would be identical to Pivot)
			if matchX == matchY {
				continue
			}

			// Determine Z (the non-shared candidate in Pincer A)
			// Mask of PincerA minus Pivot Mask leaves only Z
			zSet := pincerA.Candidates.Difference(pivotCell.Candidates)
			if zSet.Size() != 1 {
				continue
			} // Should not happen if count is 2
			zVal := zSet.Value()

			// Now look for Pincer B
			// If Pincer A matched X, Pincer B must match Y (and share Z).
			// If Pincer A matched Y, Pincer B must match X (and share Z).
			targetSharedVal := y
			if matchY {
				targetSharedVal = x
			}

			// Inner loop to find Pincer B
			for j := i + 1; j < 20; j++ {
				pincerBIdx := peers[j]
				pincerB := s.puzzle.Cell(pincerBIdx)

				if pincerB.NumCandidates() != 2 {
					continue
				}

				// Pincer B must have {targetSharedVal, zVal}
				if !pincerB.HasCandidate(targetSharedVal) || !pincerB.HasCandidate(zVal) {
					continue
				}

				// Found a valid XY-Wing!
				// Pivot: {X,Y}, PincerA: {X,Z}, PincerB: {Y,Z}
				// Elimination: Remove Z from common peers of Pincer A and Pincer B
				step := NewStep(kindXYWing)
				if s.eliminateFromIntersection(pincerAIdx, pincerBIdx, -1, zVal, step) {
					s.applyStep(step.
						WithIndices(pivotIdx, pincerAIdx, pincerBIdx).
						WithValues(x, y, zVal))
					return true
				}
			}
		}
	}
	return false
}

// findXYZ searches for 3 cells that fit the "XYZ-Wing" pattern.  An XYZ-Wing
// consists of a pivot cell with 3 candidate values x,y,z, and two pincer cells
// that each see the pivot cell but don't see each other.  One pincer must have
// candidate values x,z and the other must have candidate values y,z.  One of
// these cells must have the value z, so z can be eliminated as a candidate for
// any cell that sees all three.  Note that one pincer *MUST* be in the same
// box as the pivot cell in order for it to be possible for any cell to see the
// pivot and both pincers.
func (s *Solver) findXYZWings() bool {
	for pivotIdx := range 81 {
		pivotCell := s.puzzle.Cell(pivotIdx)

		if pivotCell.NumCandidates() != 3 {
			continue
		}

		pivotMask := pivotCell.Candidates
		peers := puzzle.GetPeers(pivotIdx)

		// Search for Pincer A
		for i := range 20 {
			pAIdx := peers[i]
			pA := s.puzzle.Cell(pAIdx)

			if pA.NumCandidates() != 2 {
				continue
			}

			maskA := pA.Candidates
			if !maskA.IsSubsetOf(pivotMask) {
				continue
			}

			// Search for Pincer B
			for j := i + 1; j < 20; j++ {
				pBIdx := peers[j]
				pB := s.puzzle.Cell(pBIdx)

				if pB.NumCandidates() != 2 {
					continue
				}
				maskB := pB.Candidates
				if !maskB.IsSubsetOf(pivotMask) {
					continue
				}

				union := maskA.Union(maskB)
				intersection := maskA.Intersection(maskB)

				if union == pivotMask && intersection.Size() == 1 {
					zVal := bitset.BitSet16(intersection).Value()

					step := NewStep(kindXYZWing)
					if s.eliminateFromIntersection(pAIdx, pBIdx, pivotIdx, zVal, step) {
						s.applyStep(step.
							WithIndices(pivotIdx, pAIdx, pBIdx).
							WithValues(pivotCell.CandidateValues()...))
						return true
					}
				}
			}
		}
	}
	return false
}
