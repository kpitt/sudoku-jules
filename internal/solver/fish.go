package solver

import (
	"github.com/kpitt/sudoku/internal/bitset"
)

func (s *Solver) findXWings() (found bool) {
	return s.findFishOfSize(2, kindXWing)
}

func (s *Solver) findSwordfish() (found bool) {
	return s.findFishOfSize(3, kindSwordfish)
}

func (s *Solver) findJellyfish() (found bool) {
	return s.findFishOfSize(4, kindJellyfish)
}

func (s *Solver) findFishOfSize(fishSize int, fishKind techniqueKind) (found bool) {
	find := func(baseLines, coverLines []*House) bool {
		return s.findFishInLines(fishSize, fishKind, baseLines, coverLines)
	}
	return find(s.rows[:], s.columns[:]) || find(s.columns[:], s.rows[:])
}

func (s *Solver) findFishInLines(
	fishSize int,
	fishKind techniqueKind,
	baseLines, coverLines []*House,
) (found bool) {
	for _, base := range baseLines {
		for val := 1; val <= 9; val++ {
			locs := base.Unsolved[val]
			if locs.Empty() {
				continue
			}

			// A fish line must have no more than fishSize candidate locations
			// for a value. We assume that all singles and smaller fish have
			// already been found.
			if locs.Size() > fishSize {
				continue
			}

			if s.checkFishForValue(fishSize, fishKind, val, base, baseLines, coverLines) {
				return true
			}
		}
	}

	return false
}

func (s *Solver) checkFishForValue(
	fishSize int,
	fishKind techniqueKind,
	val int,
	base1 *House,
	baseLines, coverLines []*House,
) (found bool) {
	// Find all base lines other than base1 that have either 2 or 3 candidate
	// locations for val.
	candidates := filterSlice(baseLines, func(b2 *House) bool {
		numLocs := b2.NumLocations(val)
		return b2.Index != base1.Index && numLocs >= 2 && numLocs <= fishSize
	})

	valueSet := bitset.FromValues16(val)
	// Variables for storing search results.
	fishLines := []int{base1.Index}
	var coverLocs LocSet

	// Must forward-declare func so we can call it self-recursively.
	var checkLines func(lines []*House, fishLocs LocSet) bool
	checkLines = func(lines []*House, fishLocs LocSet) bool {
		for i, line := range lines {
			locs := bitset.Union(fishLocs, line.Unsolved[val])
			if locs.Size() > fishSize {
				// Too many locations, so this line can't be part of the fish.
				continue
			}
			fishLines = append(fishLines, line.Index)
			if len(fishLines) == fishSize {
				// We found enough lines, so we have a fish.
				// fishLines contains the base lines that make up the fish,
				// but we need to save the set containing the cover lines.
				coverLocs = locs
				return true
			}
			// Check recursively until we have enough lines.
			if checkLines(lines[i+1:], locs) {
				return true
			}

			// No fish found, so backtrack the last line and keep trying.
			fishLines = fishLines[:len(fishLines)-1]
		}
		// No more candidate lines, so we don't have a fish.
		return false
	}

	step := NewStep(fishKind).WithValues(val)
	if checkLines(candidates, base1.Unsolved[val]) {
		// We found a fish.
		bases := transformSlice(fishLines, func(x int) *House {
			return baseLines[x]
		})
		covers := transformSlice(coverLocs.Values(), func(y int) *House {
			return coverLines[y]
		})
		locSet := bitset.FromValues16(fishLines...)
		if s.eliminateFromOtherLocsMulti(covers, valueSet, locSet, step) {
			s.applyStep(step.
				WithBases(bases...).
				WithCovers(covers...))
			return true
		}
	}

	return false
}
