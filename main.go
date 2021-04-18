package main

import (
	"fmt"

	"github.com/NouemanKHAL/sudoku-solver-rest-api/sudoku"
)

func main() {
	sG := sudoku.New()
	fmt.Printf("%v\n\n", sG.ToString())

	inputGrid := [][]rune{
		{'3', '2', '1', '7', '.', '4', '.', '.', '.'},
		{'6', '4', '.', '.', '9', '.', '.', '.', '7'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '4', '5', '9', '.', '.'},
		{'.', '.', '5', '1', '8', '7', '4', '.', '.'},
		{'.', '.', '4', '9', '6', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'2', '.', '.', '.', '7', '.', '.', '1', '9'},
		{'.', '.', '.', '6', '.', '9', '5', '8', '2'},
	}

	sG.Copy(inputGrid)
	fmt.Printf("%v\n\n", sG.ToString())

	sG.Solve()
	fmt.Printf("%v\n\n", sG.ToString())
}
