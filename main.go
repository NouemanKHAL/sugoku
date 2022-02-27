package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/NouemanKHAL/sudoku-solver-rest-api/sudoku"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
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

func LogMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		log.Debug(r.Body)
		h.ServeHTTP(w, r)
	})
}

func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/sudoku", LogMiddleware(SudokuSolverHandler)).Methods("POST")
}

func StartServer() {
	r := mux.NewRouter()
	SetupHandlers(r)

	port := GetEnvWithDefault("SUDOKU_SERVER_PORT", DEFAULT_PORT)

	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func initLogger() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% - [%lvl%] - %msg%\n",
	})
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	}
	if os.Getenv("TRACE") == "true" {
		log.SetLevel(log.TraceLevel)
	}
}

func main() {
	initLogger()
	StartServer()
}

func GetEnvWithDefault(name, fallback string) string {
	if envVar, ok := os.LookupEnv(name); ok {
		return envVar
	}
	return fallback
}
