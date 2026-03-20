// Package solver implements various Sudoku solving techniques.
package solver

import (
	"fmt"

	"github.com/fatih/color"
)

func (s *Solver) printProgress(format string, a ...any) {
	if s.EnableDebug {
		printDebug(format, a...)
	}
}

func (s *Solver) printChecking(name string) {
	s.printProgress("Checking %q technique", name)
}

func printDebug(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintln(color.Error, color.HiBlackString(">>> %s", msg))
}

// PrintStep prints a single solution step to stdout.
func (s *Solver) PrintStep(step *SolutionStep) {
	fmt.Println(s.FormatStep(step))
}

// PrintSolution prints all steps of the solution to stdout.
func (s *Solver) PrintSolution() {
	for i, step := range s.solution {
		fmt.Printf("%2d. %s\n", i+1, s.FormatStep(step))
	}
}
