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
	rowsMap     []map[rune]int
	colsMap     []map[rune]int
	subGridMap  []map[rune]int
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

	sG.rowsMap = make([]map[rune]int, 9)
	sG.colsMap = make([]map[rune]int, 9)
	sG.subGridMap = make([]map[rune]int, 9)

	for i := 0; i < 9; i++ {
		sG.rowsMap[i] = make(map[rune]int)
		sG.colsMap[i] = make(map[rune]int)
		sG.subGridMap[i] = make(map[rune]int)
	}

	return &sG
}

func (sG *sudokuGrid) reset() {
	sG = New()
}

func (sG *sudokuGrid) Solve() error {
	missingCells := [][2]int{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sG.Grid[i][j] != EMPTY_CELL {
				continue
			}
			missingCells = append(missingCells, [2]int{i, j})
		}
	}

	fmt.Printf("#Missing Cells: %d\n", len(missingCells))

	if !sG.solve(missingCells, 0) {
		return errors.New("No solution exists")
	}
	return nil
}

func (sG *sudokuGrid) solve(cells [][2]int, curIndex int) bool {
	if curIndex == len(cells) {
		return true
	}

	x := cells[curIndex][0]
	y := cells[curIndex][1]

	for val := '1'; val <= '9'; val++ {
		if !sG.checkIfExists(x, y, val) {
			oldValue := sG.Grid[x][y]

			// try this value
			sG.Set(x, y, val)

			// continue backtracking on the next cell
			if sG.solve(cells, curIndex+1) {
				return true
			}

			/// it didnt work, reset the old value
			sG.Set(x, y, oldValue)
		}
	}
	return false
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

	sG.updateCount(x, y, val)
	sG.Grid[x][y] = val

	return nil
}

func (sG *sudokuGrid) getSubgridIndex(x, y int) int {
	return (x/3)*3 + (y / 3)
}

func (sG *sudokuGrid) updateCount(x, y int, newValue rune) {
	oldValue := sG.Grid[x][y]

	sG.rowsMap[x][oldValue]--
	sG.colsMap[y][oldValue]--
	sG.subGridMap[sG.getSubgridIndex(x, y)][oldValue]--

	sG.rowsMap[x][newValue]++
	sG.colsMap[y][newValue]++
	sG.subGridMap[sG.getSubgridIndex(x, y)][newValue]++
}

func (sG *sudokuGrid) checkIfExists(x, y int, val rune) bool {
	if sG.rowsMap[x][val] > 0 || sG.colsMap[y][val] > 0 {
		return true
	}

	if sG.subGridMap[sG.getSubgridIndex(x, y)][val] > 0 {
		return true
	}

	return false
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

func (sG *sudokuGrid) Copy(grid [][]rune) {
	sG.reset()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			sG.Set(i, j, grid[i][j])
			if grid[i][j] != EMPTY_CELL {
				sG.IsGenerated[i][j] = true
			}
		}
	}
}
