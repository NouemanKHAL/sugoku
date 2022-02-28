package sudoku

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
)

const (
	EMPTY_CELL rune   = '.'
	SEP_CHAR   string = "/"
)

type SudokuGrid struct {
	Size            int      `json:"Size"`
	PartitionWidth  int      `json:"PartitionWidth"`
	PartitionHeight int      `json:"PartitionHeight"`
	Grid            [][]rune `json:"Grid"`
	rowsMap         []map[rune]bool
	colsMap         []map[rune]bool
	subGridMap      []map[rune]bool
}

type coord struct {
	x, y int
}

// New Returns an empty sG.Size x sG.Size SudokuGrid
func New(size, partitionWidth, partitionHeight int) (*SudokuGrid, error) {
	sG := SudokuGrid{
		Size:            size,
		PartitionWidth:  partitionWidth,
		PartitionHeight: partitionHeight,
	}

	sG.Grid = make([][]rune, sG.Size)

	for i := 0; i < sG.Size; i++ {
		sG.Grid[i] = make([]rune, sG.Size)
		for j := 0; j < len(sG.Grid[i]); j++ {
			sG.Grid[i][j] = EMPTY_CELL
		}
	}
	sG.initMetadata()

	if err := sG.Valid(); err != nil {
		return nil, err
	}
	return &sG, nil
}

func (sG *SudokuGrid) initMetadata() {
	sG.rowsMap = make([]map[rune]bool, sG.Size)
	sG.colsMap = make([]map[rune]bool, sG.Size)
	sG.subGridMap = make([]map[rune]bool, sG.Size)

	for i := 0; i < sG.Size; i++ {
		sG.rowsMap[i] = make(map[rune]bool)
		sG.colsMap[i] = make(map[rune]bool)
		sG.subGridMap[i] = make(map[rune]bool)
	}

	// update the state of rowsMap, colsMap, subGridMap
	for i := 0; i < sG.Size; i++ {
		for j := 0; j < len(sG.Grid[i]); j++ {
			sG.Set(i, j, sG.Grid[i][j])
		}
	}
}

// Reset sets all the cells of the SudokuGrid to EMPTY_CELL value
func (sG *SudokuGrid) Reset() {
	sG, _ = New(sG.Size, sG.PartitionWidth, sG.PartitionHeight)
}

// Solve solves the SudokuGrid in-place, returns an error if no solution exist
func (sG *SudokuGrid) Solve() error {
	missingCells := make([]coord, 0, sG.Size*sG.Size)

	for i := 0; i < sG.Size; i++ {
		for j := 0; j < len(sG.Grid[i]); j++ {
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

func (sG *SudokuGrid) solve(cells []coord) bool {
	if len(cells) == 0 {
		return true
	}

	x := cells[0].x
	y := cells[0].y

	for val := '1'; val <= rune('0'+sG.Size); val++ {
		if sG.canSet(x, y, val) {
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

// isValidIndex returns true if the coordinates (x, y) represent a valid cell, and false otherwise
func (sG *SudokuGrid) isValidIndex(x, y int) bool {
	return x >= 0 && x < sG.Size && y >= 0 && y < sG.Size
}

// Get returns the value of the cell with coordinates (x, y)
func (sG *SudokuGrid) Get(x, y int) (rune, error) {
	if !sG.isValidIndex(x, y) {
		return EMPTY_CELL, errors.New(fmt.Sprintf("Cell coordinates out of bounds. (%d, %d)", x, y))
	}
	return sG.Grid[x][y], nil
}

// Set sets the value of the cell with coordinates (x, y)
func (sG *SudokuGrid) Set(x, y int, val rune) error {
	if !sG.isValidIndex(x, y) {
		return errors.New(fmt.Sprintf("Cell coordinates out of bounds. (%d, %d)", x, y))
	}

	sG.updateCount(x, y, val)
	sG.Grid[x][y] = val

	return nil
}

// GetSubgridIndex returns the index of the partition containing the cell with coordinates (x, y) in the partitions grid - subgrid -.
func (sG *SudokuGrid) GetSubgridIndex(x, y int) int {
	// floor(x/Height) * ROW_SIZE + floor(y/W)
	// ROW_SIZE is the size of the compressed matrix where each cell represents a subgrid
	compressedMatrixWidth := sG.Size / sG.PartitionWidth
	return (x/sG.PartitionHeight)*(compressedMatrixWidth) + y/sG.PartitionWidth
}

func (sG *SudokuGrid) updateCount(x, y int, newValue rune) {
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

// canSet returns true if the given value doesn't exist in the same row (x), column (y), or subgrid
func (sG *SudokuGrid) canSet(x, y int, val rune) bool {
	if !sG.isValidIndex(x, y) {
		return false
	}
	return !(sG.rowsMap[x][val] || sG.colsMap[y][val] || sG.subGridMap[sG.GetSubgridIndex(x, y)][val])
}

func (sG *SudokuGrid) MarshalJSON() ([]byte, error) {
	return json.Marshal(*sG)
}

func (sG *SudokuGrid) UnmarshalJSON(data []byte) error {
	type tmpSudoku SudokuGrid
	var tmpStruct tmpSudoku
	err := json.Unmarshal(data, &tmpStruct)
	if err != nil {
		return err
	}
	*sG = SudokuGrid(tmpStruct)
	sG.initMetadata()
	return nil
}

// ToStringPrettify returns a formatted string representation of the SudokuGrid
func (sG *SudokuGrid) ToStringPrettify() string {
	var res strings.Builder
	res.Grow(sG.Size * (2*sG.Size + 1))
	for i := 0; i < sG.Size; i++ {
		if i > 0 && i%sG.PartitionHeight == 0 {
			fmt.Fprintf(&res, "%s\n", strings.Repeat("-", sG.Size*3+2))
		}
		for j := 0; j < len(sG.Grid[i]); j++ {
			if j > 0 && j%sG.PartitionWidth == 0 {
				fmt.Fprintf(&res, "|")
			}
			fmt.Fprintf(&res, "%2c ", sG.Grid[i][j])
		}
		fmt.Fprintf(&res, "\n")
	}
	return res.String()
}

// Valid returns all the errors if the SudokuGrid isn't valid, nil otherwise
func (sG *SudokuGrid) Valid() error {
	var result error

	if sG.Size%sG.PartitionHeight != 0 || sG.Size%sG.PartitionWidth != 0 || sG.Size%(sG.PartitionHeight*sG.PartitionWidth) != 0 {
		result = multierror.Append(result, errors.New("size must be divisible by both partition width and partition height"))
	}

	if len(sG.Grid) != sG.Size {
		result = multierror.Append(result, errors.New("grid size does not match the size attribute"))
	}

	cnt := 0
	for i := 0; i < sG.Size; i++ {
		if len(sG.Grid[i]) != sG.Size {
			cnt++
		}
	}
	if cnt > 0 {
		result = multierror.Append(result, fmt.Errorf("size of %d row(s) does not match the size attribute", cnt))
	}

	return result
}
