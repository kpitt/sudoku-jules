package solver

import (
	"time"
)

// BRUTE FORCE (LAST RESORT) TECHNIQUE

// findBruteForce uses a brute force search to find a solution for any remaining
// unsolved cells.  The search uses Donald Knuth's Algorithm X with the "Dancing
// Links" technique.
func (s *Solver) findBruteForce() bool {
	dl := NewDancingLinks(s.puzzle)
	dlOptions := &DancingLinksOptions{
		EnableDebug: s.EnableDebug,
		TimeLimit:   5 * time.Second,
	}

	solved, _ := dl.SolveWithStats(dlOptions)
	if solved {
		s.applyBruteForceSteps(dl)
	}

	return solved
}

func (s *Solver) applyBruteForceSteps(dl *DancingLinks) {
	steps := dl.GetSolution()
	for _, step := range steps {
		s.appendNextStep(NewStep(kindBruteForce).WithPlacedValue(step.Index, step.Value))
		// Place the value in the puzzle.
		s.puzzle.PlaceValue(step.Index, step.Value)
	}
}
