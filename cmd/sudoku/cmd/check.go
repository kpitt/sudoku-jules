package cmd

import (
	"fmt"
	"os"

	"github.com/kpitt/sudoku/internal/puzzle"
	"github.com/kpitt/sudoku/internal/solver"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check a Sudoku puzzle for unique solvability",
	Long:  `Read a Sudoku puzzle from standard input and check if it has zero, one, or multiple solutions.`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := puzzle.FromFile(os.Stdin)
		if err != nil {
			fatalError(err.Error())
		}

		valid, err := solver.CheckPuzzle(p)
		if err != nil {
			if err == puzzle.ErrNoSolution {
				fmt.Println("Puzzle has no solution")
				os.Exit(1)
			} else if err == puzzle.ErrMultipleSolutions {
				fmt.Println("Puzzle has multiple solutions")
				os.Exit(1)
			}
			fatalError(err.Error())
		}

		if valid {
			fmt.Println("Puzzle has one unique solution")
			os.Exit(0)
		}

		// This part should theoretically not be reached because CheckPuzzle returns error for invalid cases
		fmt.Println("Puzzle is invalid")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
