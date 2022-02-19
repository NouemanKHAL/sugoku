package sudoku_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSudoku(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sudoku Suite")
}
