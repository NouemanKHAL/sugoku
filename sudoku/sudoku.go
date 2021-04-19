package sudoku

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	EMPTY_CELL rune   = '.'
	SEP_CHAR   string = "/"
)

type sudokuGrid struct {
	Size            int
	PartitionWidth  int
	PartitionHeight int
	Grid            [][]rune
	rowsMap         []map[rune]bool
	colsMap         []map[rune]bool
	subGridMap      []map[rune]bool
}

type coord struct {
	x, y int
}

// New Returns an empty sG.Size x sG.Size Grid
func New(size, partitionWidth, partitionHeight int) (*sudokuGrid, error) {

	if size%partitionHeight != 0 || size%partitionWidth != 0 || size%(partitionHeight*partitionWidth) != 0 {
		return nil, errors.New("Size must be divisible by both Width and Height.")
	}

	sG := sudokuGrid{
		Size:            size,
		PartitionWidth:  partitionWidth,
		PartitionHeight: partitionHeight,
	}

	sG.Grid = make([][]rune, sG.Size)

	for i := 0; i < sG.Size; i++ {
		sG.Grid[i] = make([]rune, sG.Size)
		for j := 0; j < sG.Size; j++ {
			sG.Grid[i][j] = EMPTY_CELL
		}
	}

	sG.rowsMap = make([]map[rune]bool, sG.Size)
	sG.colsMap = make([]map[rune]bool, sG.Size)
	sG.subGridMap = make([]map[rune]bool, sG.Size)

	for i := 0; i < sG.Size; i++ {
		sG.rowsMap[i] = make(map[rune]bool)
		sG.colsMap[i] = make(map[rune]bool)
		sG.subGridMap[i] = make(map[rune]bool)
	}

	return &sG, nil
}

func (sG *sudokuGrid) Reset() {
	sG, _ = New(sG.Size, sG.PartitionWidth, sG.PartitionHeight)
}

func (sG *sudokuGrid) Solve() error {
	missingCells := make([]coord, 0, sG.Size*sG.Size)

	for i := 0; i < sG.Size; i++ {
		for j := 0; j < sG.Size; j++ {
			if sG.Grid[i][j] == EMPTY_CELL {
				missingCells = append(missingCells, coord{x: i, y: j})
			}
		}
	}

	if !sG.solve(missingCells) {
		return errors.New("No solution exists")
	}
	return nil
}

func (sG *sudokuGrid) solve(cells []coord) bool {
	if len(cells) == 0 {
		return true
	}

	x := cells[0].x
	y := cells[0].y

	for val := '1'; val <= rune('0'+sG.Size); val++ {
		if !sG.checkIfExists(x, y, val) {
			oldValue := sG.Grid[x][y]

			// try this value
			sG.Set(x, y, val)

			// continue backtracking on the next cell
			if sG.solve(cells[1:]) {
				return true
			}

			/// it didnt work, reset the old value
			sG.Set(x, y, oldValue)
		}
	}
	return false
}

func (sG *sudokuGrid) isValidIndex(x, y int) bool {
	return x >= 0 || x < sG.Size || y >= 0 || y < sG.Size
}

// Get returns the value of the cell with coordinates (x, y)
func (sG *sudokuGrid) Get(x, y int) (rune, error) {
	if !sG.isValidIndex(x, y) {
		return EMPTY_CELL, errors.New(fmt.Sprintf("Cell coordinates out of bounds. (%d, %d)", x, y))
	}
	return sG.Grid[x][y], nil
}

// Set sets the value of the cell with coordinates (x, y)
func (sG *sudokuGrid) Set(x, y int, val rune) error {
	if !sG.isValidIndex(x, y) {
		return errors.New(fmt.Sprintf("Cell coordinates out of bounds. (%d, %d)", x, y))
	}

	sG.updateCount(x, y, val)
	sG.Grid[x][y] = val

	return nil
}

func (sG *sudokuGrid) GetSubgridIndex(x, y int) int {
	// floor(x/A) * ROW_SIZE + floor(y/B)
	// ROW_SIZE is the size of the compressed matrix where each cell represents a subgrid
	subGridMatrixWidth := sG.Size / sG.PartitionWidth
	return (x/sG.PartitionHeight)*(subGridMatrixWidth) + y/sG.PartitionWidth
}

func (sG *sudokuGrid) updateCount(x, y int, newValue rune) {
	oldValue := sG.Grid[x][y]

	// decrement row, col, subgrid count of the oldValue
	sG.rowsMap[x][oldValue] = false
	sG.colsMap[y][oldValue] = false
	sG.subGridMap[sG.GetSubgridIndex(x, y)][oldValue] = false

	// increment row, col, subgrid count of the newValue
	sG.rowsMap[x][newValue] = true
	sG.colsMap[y][newValue] = true
	sG.subGridMap[sG.GetSubgridIndex(x, y)][newValue] = true
}

// checkIfExists returns true if the value exists in the same row, or same column, or same subgrid
func (sG *sudokuGrid) checkIfExists(x, y int, val rune) bool {
	return sG.rowsMap[x][val] || sG.colsMap[y][val] || sG.subGridMap[sG.GetSubgridIndex(x, y)][val]
}

func Serialize(sG *sudokuGrid) string {
	res := fmt.Sprintf("%d/%d/%d/", sG.Size, sG.PartitionHeight, sG.PartitionWidth)
	for _, row := range sG.Grid {
		for _, cell := range row {
			res += string(cell)
		}
		res += SEP_CHAR
	}
	return res
}

func Deserialize(data string) (*sudokuGrid, error) {
	tokens := strings.Split(data, SEP_CHAR)
	if len(tokens) < 4 {
		return nil, errors.New("Missing data or Incorrect format")
	}
	size, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, err
	}
	partitionWidth, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, err
	}
	partitionHeight, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, err
	}
	sG, err := New(size, partitionWidth, partitionHeight)
	if err != nil {
		return nil, err
	}
	sG.CopyString(tokens[2:])
	return sG, nil
}

func (sG *sudokuGrid) ToStringPrettify() string {
	res := ""
	for i := 0; i < sG.Size; i++ {
		if i > 0 && i%sG.PartitionHeight == 0 {
			res += fmt.Sprintf(strings.Repeat("-", sG.Size*3+2)) + "\n"
		}
		for j := 0; j < sG.Size; j++ {
			if j > 0 && j%sG.PartitionWidth == 0 {
				res += fmt.Sprintf("|")
			}
			res += fmt.Sprintf("%2c ", sG.Grid[i][j])
		}
		res += "\n"
	}
	return res
}

func (sG *sudokuGrid) Copy(grid [][]rune) {
	sG.Reset()
	for i := 0; i < sG.Size; i++ {
		for j := 0; j < sG.Size; j++ {
			sG.Set(i, j, grid[i][j])
		}
	}
}

func (sG *sudokuGrid) CopyString(grid []string) {
	sG.Reset()
	for i := 0; i < sG.Size; i++ {
		for j := 0; j < sG.Size; j++ {
			sG.Set(i, j, rune(grid[i][j]))
		}
	}
}

/*

TODO:
	- Make Sudoku Grid Size Customizable SudokuGrid(N, A, B) 	-- DONE
	- Remove isGenerated										-- DONE
	- keep rune? use []int instead of map[rune]int ?



*/
