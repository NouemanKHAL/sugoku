package sudoku

import (
	"errors"
	"fmt"
)

const (
	EMPTY_CELL rune = '.'
)

type sudokuGrid struct {
	Grid        [][]rune
	IsGenerated [][]bool
	validRows   []bool
	validCols   []bool
}

// New Returns an empty 9 x 9 Grid
func New() *sudokuGrid {
	sG := sudokuGrid{}
	sG.Grid = make([][]rune, 9)
	sG.IsGenerated = make([][]bool, 9)
	for i := 0; i < 9; i++ {
		sG.Grid[i] = make([]rune, 9)
		sG.IsGenerated[i] = make([]bool, 9)
		for j := 0; j < 9; j++ {
			sG.Grid[i][j] = EMPTY_CELL
		}
	}
	return &sG
}

// Get returns the value of the cell with coordinates (x, y)
func (sG *sudokuGrid) Get(x, y int) (rune, error) {
	if x < 0 || x > 8 || y < 0 || y > 8 {
		return EMPTY_CELL, errors.New(fmt.Sprintf("Cell coordinates out of bounds. (%d, %d)", x, y))
	}
	return sG.Grid[x][y], nil
}

// Set sets the value of the cell with coordinates (x, y)
func (sG *sudokuGrid) Set(x, y int, val rune) error {
	if x < 0 || x > 8 || y < 0 || y > 8 {
		return errors.New(fmt.Sprintf("Cell coordinates out of bounds. (%d, %d)", x, y))
	}
	if sG.IsGenerated[x][y] {
		return errors.New(fmt.Sprintf("Cannot Set generated cell (%d, %d)", x, y))
	}
	sG.Grid[x][y] = val
	return nil
}

func (sG *sudokuGrid) IsValid() bool {
	if len(sG.Grid) != 9 {
		return false
	}
	for _, row := range sG.Grid {
		if len(row) != 9 {
			return false
		}
	}
	return true
}

func (sG *sudokuGrid) ToString() string {
	res := ""
	for _, row := range sG.Grid {
		for _, cell := range row {
			res += string(cell)
		}
		res += "\n"
	}
	return res
}
