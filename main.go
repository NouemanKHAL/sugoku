package main

import (
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
	w.Write([]byte("Welcome to the Sudoki REST API v0.0.1"))
}

func SudokuGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement a sudoku generator
	params := r.URL.Query()
	size := params.Get("size")
	partitionHeight := params.Get("partitionHeight")
	partitionWidth := params.Get("partitionWidth")
	w.Write([]byte("generator under construction\nreceived " + fmt.Sprintf("size:%s,pHeight:%s,pWidth:%s", size, partitionHeight, partitionWidth)))
}

func SudokuSolverHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error reading the body: %s", err)))
	}

	sG := &sudoku.SudokuGrid{}
	err = sG.UnmarshalJSON(body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error unmarshaling the response: %s", err)))
	}

	solvedGridBytes, err := sG.MarshalJSON()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error marshalling the solution: %s", err)))
	}
	w.Write(solvedGridBytes)
}

func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/sudoku", SudokuGeneratorHandler).Methods("GET")
	r.HandleFunc("/sudoku", SudokuSolverHandler).Methods("POST")
}

func StartServer() {
	r := mux.NewRouter()
	SetupHandlers(r)

	port := GetEnvWithDefault("SUDOKU_SERVER_PORT", DEFAULT_PORT)

	log.Printf("starting server under localhost:%s\n", port)
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
