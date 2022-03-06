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
	Short: "Sugoku is a sudoku REST API written in Go that allows to solve and generate sudoku grids with custom grid and subgrid sizes",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
