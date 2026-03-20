package puzzle

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	borderTop    = "┌─────┬─────┬─────╥─────┬─────┬─────╥─────┬─────┬─────┐"
	borderBot    = "└─────┴─────┴─────╨─────┴─────┴─────╨─────┴─────┴─────┘"
	dividerMinor = "├─────┼─────┼─────╫─────┼─────┼─────╫─────┼─────┼─────┤"
	dividerMajor = "╞═════╪═════╪═════╬═════╪═════╪═════╬═════╪═════╪═════╡"
	edgeMinor    = "│"
	edgeMajor    = "║"
)

var (
	givenColor     = color.New(color.FgHiBlue)
	givenLegend    = givenColor.Sprint("Blue")
	solvedColor    = color.New(color.FgHiGreen)
	solvedLegend   = solvedColor.Sprint("Green")
	unsolvedColor  = color.New(color.FgHiBlack)
	unsolvedLegend = unsolvedColor.Sprint("Gray")
)

// Print outputs the current state of the puzzle to the console in a
// human-readable grid format.
func (p *Puzzle) Print() {
	fmt.Println("┌───────┬───────┬───────┐")
	for r := range 9 {
		if r == 3 || r == 6 {
			fmt.Println("├───────┼───────┼───────┤")
		}
		fmt.Print("│ ")
		for c := range 9 {
			if c == 3 || c == 6 {
				fmt.Print("│ ")
			}
			cell := p.Get(r, c)
			if cell.IsSolved() {
				if cell.IsGiven {
					givenColor.Printf("%d ", cell.Value())
				} else {
					solvedColor.Printf("%d ", cell.Value())
				}
			} else {
				unsolvedColor.Print("· ")
			}
		}
		fmt.Println("│")
	}
	fmt.Println("└───────┴───────┴───────┘")
	fmt.Printf("Legend: %s = Given, %s = Solved, %s = Empty\n",
		givenLegend, solvedLegend, unsolvedLegend)
}

// PrintUnsolvedCounts outputs the number of unsolved cells and the remaining
// counts for each digit.
func (p *Puzzle) PrintUnsolvedCounts() {
	color.HiWhite("Unsolved Digits (%d cells):", p.unsolvedCounts[0])
	for i := range 9 {
		digit := i + 1
		if !p.IsDigitSolved(digit) {
			fmt.Printf("%d: %d remaining\n", digit, p.unsolvedCounts[digit])
		} else {
			fmt.Printf("%d: complete\n", digit)
		}
	}
	fmt.Println()
}

// PrintCandidateGrid outputs a detailed 9x9 grid showing all remaining
// candidates for each unsolved cell.
func (p *Puzzle) PrintCandidateGrid() {
	fmt.Println(borderTop)
	for r := range 9 {
		if r != 0 {
			if r%3 == 0 {
				fmt.Println(dividerMajor)
			} else {
				fmt.Println(dividerMinor)
			}
		}
		p.printRow(r)
	}
	fmt.Println(borderBot)
	fmt.Printf("Legend: %s = Given, %s = Solved, %s = Candidate\n",
		givenLegend, solvedLegend, unsolvedLegend)
}

// FormatCell returns a human-readable string representation of a cell's
// position (e.g., "r1c1" for the top-left cell).
func FormatCell(index int) string {
	r, c := index/9, index%9
	return fmt.Sprintf("r%dc%d", r+1, c+1)
}

func (p *Puzzle) printRow(r int) {
	for cr := range 3 {
		p.printCandidateRow(r, cr)
	}
}

func (p *Puzzle) printCandidateRow(r, candidateRow int) {
	for c := range 9 {
		if c != 0 && c%3 == 0 {
			fmt.Print(edgeMajor)
		} else {
			fmt.Print(edgeMinor)
		}
		cell := p.Get(r, c)
		if cell.IsSolved() {
			if candidateRow == 1 {
				if cell.IsGiven {
					givenColor.Printf(" [%d] ", cell.Value())
				} else {
					solvedColor.Printf("  %d  ", cell.Value())
				}
			} else {
				fmt.Print("     ")
			}
		} else {
			cell.printCandidates(candidateRow)
		}
	}
	fmt.Println(edgeMinor)
}

func (c *Cell) printCandidates(candidateRow int) {
	candidateBase := candidateRow*3 + 1
	for col := range 3 {
		if col > 0 {
			// Add a space between candidates.
			fmt.Print(" ")
		}
		candidate := candidateBase + col
		if c.HasCandidate(candidate) {
			unsolvedColor.Printf("%d", candidate)
		} else {
			fmt.Print(" ")
		}
	}
}
