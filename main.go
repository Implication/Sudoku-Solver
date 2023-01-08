package main

import "fmt"

// Version 1, take a static sudoku puzzle and come up with a solution, print out the rows in CLI format
func main() {
	puzz := [][]int{
		{3, 4, 2, 0, 1, 0, 0, 6, 0},
		{5, 9, 8, 0, 0, 3, 0, 0, 0},
		{0, 6, 1, 2, 0, 5, 8, 0, 4},
		{4, 5, 0, 0, 7, 8, 0, 9, 2},
		{8, 0, 9, 6, 0, 4, 0, 5, 7},
		{1, 3, 7, 5, 2, 0, 6, 4, 0},
		{6, 0, 0, 9, 0, 0, 7, 8, 3},
		{2, 0, 5, 0, 8, 6, 0, 0, 0},
		{0, 8, 3, 0, 4, 1, 5, 2, 0},
	}

	emptySpaces := make([][2]int, 2)
	//Append the first empty space before beginning the loop
	emptySpaces[0] = getEmptySpace(puzz)

	for emptySpaces[0][0] != -1 {
		validNumberFound := false
		/*
			The index should equal the current number in the empty space plus one
			If the number is zero we start with 1 and go up until we have a hit,
			If the number we hit was valid, but we need to backtrack another guess, we try another guess
		*/
		i := puzz[emptySpaces[0][0]][emptySpaces[0][1]] + 1
		//Void out whatever number is currently in the puzzle, it is not valid
		puzz[emptySpaces[0][0]][emptySpaces[0][1]] = 0
		for ;i < 10; i++ {
			if isNumValid(puzz, i, emptySpaces[0][0], emptySpaces[0][1]) {
				puzz[emptySpaces[0][0]][emptySpaces[0][1]] = i
				validNumberFound = true
				break
			}
		}
			if !validNumberFound {
				emptySpaces = emptySpaces[1:] // Pop off the stack
				continue                      //Go to the previous position, i.e. backtrack and try a different number
			}
			//Append to the top of the stack a new empty space
			emptySpaces = append([][2]int{getEmptySpace(puzz)},emptySpaces...)
		}

	displayBoard(puzz)
}

func displayBoard(puzzle [][]int) {
	for _, line := range puzzle {
		fmt.Println(line)
	}
}

func isNumValid(puzz [][]int, guess, row, column int) bool {
	for index := range puzz {
		/*
			Sudoku rules: A number entered has to be different from
			the numbers in the row, column and grid.

			If the row or columm
			has a similar number to our guess, then that number is not valid

			If the grid (3x3 square) contains our guessed number, then that number is also not valid
		*/
		if puzz[row][index] == guess && column != index {
			return false
		}
		if puzz[index][column] == guess && row != index {
			return false
		}
	}

	gridRowIndex := ((row / 3) * 3)
	gridColIndex := ((column / 3) * 3)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if puzz[gridRowIndex][gridColIndex] == guess && row != gridRowIndex && column != gridColIndex {
				return false
			}
			gridColIndex++
		}
		gridRowIndex++
		gridColIndex = ((column / 3) * 3)
	}

	return true
}

func getEmptySpace(puzz [][]int) [2]int {
	for i := range puzz {
		for j := range puzz[i] {
			if puzz[i][j] == 0 {
				return [2]int{i, j}
			}
		}
	}
	return [2]int{-1, -1}
}
