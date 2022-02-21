package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/NouemanKHAL/sudoku-solver-rest-api/sudoku"
	"github.com/gorilla/mux"
)

var (
	DEFAULT_PORT = "7007"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Sudoku Solver REST API v0.0.1"))
}

func SudokuSolverHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pretty := params.Get("pretty")
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error reading the body: %v", err)))
	}

	sG := sudoku.SudokuGrid{}
	err = json.Unmarshal(body, &sG)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error unmarshaling the response: %v", err)))
	}
	err = sG.Solve()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error solving the sudoku grid: %v", err)))
	}

	var res []byte

	if pretty == "true" {
		res = []byte(sG.ToStringPrettify())
	} else {
		res, err = json.Marshal(sG)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error marshalling the solution: %v", err)))
		}
	}
	w.Write(res)
}

func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/sudoku", SudokuSolverHandler).Methods("POST")
}

func StartServer() {
	r := mux.NewRouter()
	SetupHandlers(r)

	port := GetEnvWithDefault("SUDOKU_SERVER_PORT", DEFAULT_PORT)

	log.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func main() {
	StartServer()
}

func GetEnvWithDefault(name, fallback string) string {
	if envVar, ok := os.LookupEnv(name); ok {
		return envVar
	}
	return fallback
}
