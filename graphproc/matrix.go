package graphproc

import (
	"strings"
)

//MartrixCut : delete fields from a matrix
func MatrixCut(matrix [][]string, cutrow []bool) [][]string {
	for i, row := range matrix {
		newrow := make([]string, 0, 128)
		for j, c := range cutrow {
			if !c {
				newrow = append(newrow, row[j])
			}
		}
		matrix[i] = newrow
	}
	return matrix
}

//MatrixSelect : select fields from a matrix
func MatrixSelect(matrix [][]string, selectrow []bool) [][]string {
	for i, row := range matrix {
		var newrow []string
		for j, c := range selectrow {
			if c {
				newrow = append(newrow, row[j:j+1]...)
			}
		}
		matrix[i] = newrow
	}
	return matrix
}

//BytesMatrix : create matrix from text stored in a byte slice
func BytesMatrix(buffer []byte) [][]string {
	matrix := make([][]string, 0, 128)
	textlines := strings.Split(string(buffer), "\n")
	for i, line := range textlines {
		matrix = append(matrix, make([]string, 0, 128))
		matrix[i] = append(matrix[i], strings.Split(line, ",")...)
	}
	return matrix
}

//MatrixBytes : convert matrix to text in a byte slice
func MatrixBytes(matrix [][]string, buffer []byte) ([]byte, error) {
	for _, row := range matrix {
		byterow := make([]byte, 0, 2048)
		for j, col := range row {
			if j != 0 {
				byterow = append(byterow, []byte(",")...)
			}
			byterow = append(byterow, []byte(col)...)
		}
		buffer = append(buffer, byterow...)
		buffer = append(buffer, []byte("\n")...)
	}
	return buffer, nil
}
