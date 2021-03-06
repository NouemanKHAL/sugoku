package sudoku

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	EMPTY_CELL rune   = '.'
	SEP_CHAR   string = "/"
)

type SudokuGrid struct {
	Size            int      `json:"size"`
	PartitionWidth  int      `json:"partitionWidth"`
	PartitionHeight int      `json:"partitionHeight"`
	Grid            [][]rune `json:"grid"`
	rowsMap         []map[rune]bool
	colsMap         []map[rune]bool
	subGridMap      []map[rune]bool
	allowedValues   []rune
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
	if err := sG.Valid(); err != nil {
		return nil, err
	}

	sG.initMetadata()
	return &sG, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	sG.allowedValues = make([]rune, 0, max(sG.Size, sG.PartitionHeight*sG.PartitionWidth))
	for val := '1'; val <= rune('0'+cap(sG.allowedValues)); val++ {
		sG.allowedValues = append(sG.allowedValues, val)
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
		return errors.New("no solution exists")
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

			// it didnt work, reset the old value
			sG.Set(x, y, oldValue)
		}
	}
	return false
}

// isValidIndex returns true if the coordinates (x, y) represent a valid cell, and false otherwise
func (sG *SudokuGrid) isValidIndex(x, y int) error {
	if x < 0 || x >= sG.Size || y < 0 || y >= sG.Size {
		return fmt.Errorf("cell coordinates out of bounds (%d, %d)", x, y)
	}
	return nil
}

// Get returns the value of the cell with coordinates (x, y)
func (sG *SudokuGrid) Get(x, y int) (rune, error) {
	if err := sG.isValidIndex(x, y); err != nil {
		return EMPTY_CELL, err
	}
	return sG.Grid[x][y], nil
}

// Set sets the value of the cell with coordinates (x, y)
func (sG *SudokuGrid) Set(x, y int, val rune) error {
	if err := sG.isValidIndex(x, y); err != nil {
		return err
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
	if err := sG.isValidIndex(x, y); err != nil {
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
	err = sG.Valid()
	if err != nil {
		return err
	}
	sG.initMetadata()
	return nil
}

// GenerateSudokuGrid returns a SudokuGrid with the given dimensions
func GenerateSudokuGrid(size, partitionWidth, partitionHeight int) (*SudokuGrid, error) {
	sG, err := New(size, partitionWidth, partitionHeight)
	if err != nil {
		return nil, err
	}

	// shuffling the allowed values => random puzzle generation
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(sG.allowedValues), func(i, j int) { sG.allowedValues[i], sG.allowedValues[j] = sG.allowedValues[j], sG.allowedValues[i] })

	log.Debugf("generating sudoku grid using the allowed values: %v\n", sG.allowedValues)

	if generateSudokuGrid(sG, 0, 0) {
		return sG, nil
	}

	return nil, errors.New("could not generate a valid sudoku grid")
}

func generateSudokuGrid(sG *SudokuGrid, i, j int) bool {
	// surpasses the last row or last cell (sG.Size - 1, sG.Size - 1)
	if i >= sG.Size {
		return true
	}
	for _, val := range sG.allowedValues {
		if sG.canSet(i, j, val) {
			// try this value
			sG.Set(i, j, val)

			// determine our next cell
			newI, newJ := i, j+1
			if j+1 >= sG.Size {
				newI = i + 1
				newJ = 0
			}

			// continue backtracking on the next cell
			if generateSudokuGrid(sG, newI, newJ) {
				return true
			}

			// it didnt work, reset the old value
			sG.Set(i, j, EMPTY_CELL)
		}
	}

	return false
}

func getLevelThreshold(level string) (float64, error) {
	switch level {
	case "easy":
		return 0.5, nil
	case "medium":
		return 0.65, nil
	case "hard":
		return 0.8, nil
	case "extreme":
		return 0.95, nil
	case "robot":
		return 1, nil
	}
	return 0, errors.New("invalid level: must be one of the supported levels (easy, medium, hard, extreme, robot)")
}

// SetGridTolevel adds empty cells to match the desired difficulty level
func (sG *SudokuGrid) SetGridToLevel(level string) error {
	threshold, err := getLevelThreshold(level)
	if err != nil {
		return err
	}
	for i := 0; i < sG.Size; i++ {
		for j := 0; j < len(sG.Grid[i]); j++ {
			if rand.Float64() < threshold {
				sG.Set(i, j, EMPTY_CELL)
			}
		}
	}
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
	if sG.Size%sG.PartitionHeight != 0 || sG.Size%sG.PartitionWidth != 0 || sG.Size%(sG.PartitionHeight*sG.PartitionWidth) != 0 {
		return errors.New("size must be divisible by both partitionWidth and partitionHeight")
	}
	if len(sG.Grid) != sG.Size {
		return errors.New("the given grid size does not match the given size property")
	}

	cnt := 0
	for i := 0; i < len(sG.Grid); i++ {
		if len(sG.Grid[i]) != sG.Size {
			cnt++
		}
	}
	if cnt > 0 {
		return fmt.Errorf("%d row(s) sizes do not match the given size property", cnt)
	}

	return nil
}
