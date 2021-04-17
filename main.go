package main

import (
	"fmt"

	"github.com/NouemanKHAL/sudoku-solver-rest-api/sudoku"
)

func main() {
	sG := sudoku.New()
	fmt.Print(sG.ToString())
}
