/*
Copyright © 2022 Noueman KHALIKINE <noueman.khal@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sugoku",
	Short: "Sugoku is a simple Sudoku REST API (written in Go) designed to solve and generate Sudoku puzzles with customizable grid dimensions as well as subgrid dimensions",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
