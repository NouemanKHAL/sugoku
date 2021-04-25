package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NouemanKHAL/sudoku-solver-rest-api/sudoku"
)

func solveSudoku(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Body Error: received %s", string(body))
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	grid := string(body)
	sG, err := sudoku.Deserialize(grid)
	if err != nil {
		fmt.Fprintf(w, "Deserialization Error: received %s", string(body))
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Initial Sudoku:\n%s", sG.ToStringPrettify())
	sG.Solve()
	fmt.Fprintf(w, "Solved:\n%s", sG.ToStringPrettify())
}

func main() {
	http.HandleFunc("/sudoku/solve", solveSudoku)
	http.ListenAndServe("localhost:8080", nil)
}

/*
sG, err := sudoku.New(9, 3, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n\n", sG.ToStringPrettify())

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
	sG.Copy(inputGrid)
	serialized := sudoku.Serialize(sG)
	fmt.Println(serialized)
	deserialized, _ := sudoku.Deserialize(serialized)
	fmt.Println(deserialized.ToStringPrettify())
	deserialized.Solve()
	fmt.Println(deserialized.ToStringPrettify())
*/
