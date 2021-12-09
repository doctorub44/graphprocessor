package graphproc

import (
	"fmt"
	"testing"
)

func TestMatrixCut(t *testing.T) {

	matrix := [][]string{
		{"0", "1", "2", "3", "4"},
		{"00", "11", "22", "33", "44"},
		{"000", "111", "222", "333", "444"},
	}

	cutrow := []bool{false, true, false, false, false}
	newmatrix := MatrixCut(matrix, cutrow)
	fmt.Println(newmatrix)

	cutrow = []bool{true, false, false, false}
	newmatrix = MatrixCut(matrix, cutrow)
	fmt.Println(newmatrix)

	cutrow = []bool{false, false, true}
	newmatrix = MatrixCut(matrix, cutrow)
	fmt.Println(newmatrix)
}

func TestMatrixSelect(t *testing.T) {

	matrix := [][]string{
		{"0", "1", "2", "3", "4"},
		{"00", "11", "22", "33", "44"},
		{"000", "111", "222", "333", "444"},
	}
	selrow := []bool{false, true, false, false, false}
	newmatrix := MatrixSelect(matrix, selrow)
	fmt.Println(newmatrix)

	matrix = [][]string{
		{"0", "1", "2", "3", "4"},
		{"00", "11", "22", "33", "44"},
		{"000", "111", "222", "333", "444"},
	}
	selrow = []bool{true, false, false, false, false, false}
	newmatrix = MatrixSelect(matrix, selrow)
	fmt.Println(newmatrix)

	matrix = [][]string{
		{"0", "1", "2", "3", "4"},
		{"00", "11", "22", "33", "44"},
		{"000", "111", "222", "333", "444"},
	}
	selrow = []bool{false, false, false, false, true}
	newmatrix = MatrixSelect(matrix, selrow)
	fmt.Println(newmatrix)

	matrix = [][]string{
		{"0", "1", "2", "3", "4"},
		{"00", "11", "22", "33", "44"},
		{"000", "111", "222", "333", "444"},
	}
	selrow = []bool{true, false, true, false, true}
	newmatrix = MatrixSelect(matrix, selrow)
	fmt.Println(newmatrix)
}
