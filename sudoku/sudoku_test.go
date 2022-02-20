package sudoku

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sudoku", func() {
	Context("Creating a Sudoku Grid", func() {
		It("returns a SudokuGrid if the input is valid", func() {
			sG, err := New(9, 3, 3)
			Expect(err).To(BeNil())
			Expect(sG).NotTo(BeNil())
			Expect(sG.Size).To(Equal(9))
			Expect(sG.PartitionWidth).To(Equal(3))
			Expect(sG.PartitionHeight).To(Equal(3))
		})
		It("returns an error if the input is invalid", func() {
			sG, err := New(9, 4, 3)
			Expect(err).NotTo(BeNil())
			Expect(sG).To(BeNil())

			sG, err = New(9, 2, 3)
			Expect(err).NotTo(BeNil())
			Expect(sG).To(BeNil())

			sG, err = New(9, 9, 9)
			Expect(err).NotTo(BeNil())
			Expect(sG).To(BeNil())
		})
	})

	Context("Sudoku Grid serialization", func() {
		It("successfully serializes and deserializes a sudoku grid", func() {
			sG1, err := New(9, 3, 3)
			Expect(err).To(BeNil())

			sudokuGridBytes, err := json.Marshal(sG1)
			Expect(err).To(BeNil())
			Expect(string(sudokuGridBytes)).To(Equal(`{"Size":9,"PartitionWidth":3,"PartitionHeight":3,"Grid":[[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46],[46,46,46,46,46,46,46,46,46]]}`))

			sG2 := &SudokuGrid{}
			err = json.Unmarshal(sudokuGridBytes, sG2)
			Expect(err).To(BeNil())
			Expect(sG1.Size).To(Equal(sG2.Size))
			Expect(sG1.PartitionHeight).To(Equal(sG2.PartitionHeight))
			Expect(sG1.PartitionWidth).To(Equal(sG2.PartitionWidth))
			Expect(sG1.Grid).To(Equal(sG2.Grid))
			Expect(len(sG2.rowsMap)).To(Equal(9))
			Expect(len(sG2.colsMap)).To(Equal(9))
			Expect(len(sG2.subGridMap)).To(Equal(9))
		})
	})

	Context("Helper functions", func() {

		var (
			sG *SudokuGrid
		)
		BeforeEach(func() {
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

			var err error
			sG, err = New(9, 3, 3)
			Expect(err).To(BeNil())
			Expect(sG).NotTo(BeNil())

			for i := 0; i < len(inputGrid); i++ {
				for j := 0; j < len(inputGrid[i]); j++ {
					sG.Set(i, j, inputGrid[i][j])
				}
			}
		})

		Context("isValidIndex method", func() {
			It("isValidIndex returns true if the given coordinates are not out of range", func() {
				Expect(sG.isValidIndex(0, 0)).To(BeTrue())
				Expect(sG.isValidIndex(sG.Size-1, 0)).To(BeTrue())
				Expect(sG.isValidIndex(0, sG.Size-1)).To(BeTrue())
			})

			It("isValidIndex returns false if the given coordinates are out of range", func() {
				Expect(sG.isValidIndex(-1, 0)).To(BeFalse())
				Expect(sG.isValidIndex(0, -1)).To(BeFalse())
				Expect(sG.isValidIndex(sG.Size, 0)).To(BeFalse())
				Expect(sG.isValidIndex(0, sG.Size)).To(BeFalse())
			})
		})

		Context("canSet method", func() {
			It("canSet returns true if the value doesn't exist in the same row, the same column and the same subgrid", func() {
				Expect(sG.canSet(0, 0, '9')).To(BeTrue()) // only valid value
				Expect(sG.canSet(0, 0, '9')).To(BeTrue()) // only valid value
			})

			It("canSet returns false if the value exists in the same row, same column or the same subgrid", func() {
				Expect(sG.canSet(0, 0, '1')).To(BeFalse()) // same row & column
				Expect(sG.canSet(0, 0, '2')).To(BeFalse()) // same row & column
				Expect(sG.canSet(0, 0, '3')).To(BeFalse()) // same row
				Expect(sG.canSet(0, 0, '4')).To(BeFalse()) // same subgrid
				Expect(sG.canSet(0, 0, '5')).To(BeFalse()) // same subgrid
				Expect(sG.canSet(0, 0, '7')).To(BeFalse()) // same column
				Expect(sG.canSet(0, 0, '8')).To(BeFalse()) // same column
			})

			It("canSet returns false if the given coordinates are not valid", func() {
				Expect(sG.canSet(sG.Size, 0, '1')).To(BeFalse())
				Expect(sG.canSet(0, sG.Size, '1')).To(BeFalse())
				Expect(sG.canSet(-1, 0, '1')).To(BeFalse())
				Expect(sG.canSet(0, -1, '2')).To(BeFalse())
			})
		})

	})
})
