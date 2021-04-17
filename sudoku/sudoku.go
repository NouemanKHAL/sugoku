package sudoku

const (
	EMPTY_CELL rune = '.'
)

type sudokuGrid struct {
	Grid [][]rune
}

// New Returns an empty 9 x 9 Grid
func New() *sudokuGrid {
	sG := sudokuGrid{}
	sG.Grid = make([][]rune, 9)
	for i := 0; i < 9; i++ {
		sG.Grid[i] = make([]rune, 9)
		for j := 0; j < 9; j++ {
			sG.Grid[i][j] = EMPTY_CELL
		}
	}
	return &sG
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
