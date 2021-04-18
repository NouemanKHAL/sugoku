package main

import (
	"fmt"

	"github.com/NouemanKHAL/sudoku-solver-rest-api/sudoku"
)

func main() {
	sG := sudoku.New()
	fmt.Printf("%v\n\n", sG.ToString())

	inputGrid := [][]rune{
		{'.', '.', '.', '.', '.', '3', '2', '1', '.'},
		{'1', '2', '.', '.', '.', '.', '4', '3', '6'},
		{'.', '5', '4', '.', '2', '1', '.', '.', '.'},
		{'2', '9', '5', '1', '.', '7', '3', '.', '8'},
		{'.', '3', '5', '8', '5', '.', '6', '.', '.'},
		{'7', '.', '6', '.', '9', '4', '.', '.', '2'},
		{'8', '.', '.', '4', '.', '5', '9', '.', '3'},
		{'.', '7', '.', '.', '.', '.', '.', '2', '.'},
		{'.', '4', '9', '2', '.', '6', '7', '8', '.'},
	}

	sG.Copy(inputGrid)
	fmt.Printf("%v\n\n", sG.ToString())

	sG.Solve()
	fmt.Printf("%v\n\n", sG.ToString())
}
