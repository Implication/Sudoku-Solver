package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Data struct {
	Seed       int
	difficulty string
	Puzzle     string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	url := "https://sudoku-generator1.p.rapidapi.com/sudoku/generate"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", os.Getenv("X-RapidAPI-Key"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("X-RapidAPI-Host"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var payload Data
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	var row []int
	var puzz [][]int
	for i, v := range payload.Puzzle {
		var val int
		if string(v) == "." {
			val = 0
		} else {
			val, err = strconv.Atoi(string(v))
			if err != nil {
				log.Fatal("Error during string conversion: ", err)
			}
		}
		row = append(row, val)
		if (i+1)%9 == 0 {
			puzz = append(puzz, row)
			row = nil
		}
	}

	fmt.Println("Initial board pulled from rapid api")
	displayBoard(puzz)
	fmt.Println()

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
		for ; i < 10; i++ {
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
		emptySpaces = append([][2]int{getEmptySpace(puzz)}, emptySpaces...)
	}

	fmt.Println("Solution to generated sudoku puzzle")
	fmt.Println()
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

//Version 2: Pull from an api, and solve the puzzle

/*
	TODO:
	- Curate a list of random sudoku puzzles to read from if a person does not have a rapid api key to use
	- Create a home page that will pull in an initial puzzle from either the api or a given json file list
	- On that home page, add in a "solve" button, that will take the user to a solution of the puzzle
	- Organize this project properly based on golang standards
	- Add in testing for the solver methods, and get method
	-(Maybes) Create a way to be able to make a pdf file or something of the sudoku file (maybe a browsers print is enough?)
	-(Maybe) Create a way to be able to email to a user a sudoku puzzle file

	- Long term stuff to maybe do
		- Introduce svelte and javascript for reactive solver
		- Make the solver reactive, fire off multiple requests to change each position to show it changing in real time.
*/
