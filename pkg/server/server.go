package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/NouemanKHAL/sugoku/pkg/config"
	"github.com/NouemanKHAL/sugoku/pkg/middleware"
	"github.com/NouemanKHAL/sugoku/pkg/sudoku"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	"github.com/hashicorp/go-multierror"
)

func initLogger(cfg config.Config) {
	log.SetReportCaller(true)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% - [%lvl%] - %msg%\n",
	})

	switch cfg.LogLevel {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Sudoku REST API v0.0.1"))
}

/*
TODO:Generator will product same output for the same input... we need to randomize it
*/
func sudokuGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pretty := params.Get("pretty")
	level := params.Get("level")

	var result error
	size, err := strconv.Atoi(params.Get("size"))
	if err != nil {
		result = multierror.Append(result, err)
	}
	partitionWidth, err := strconv.Atoi(params.Get("partitionWidth"))
	if err != nil {
		result = multierror.Append(result, err)
	}
	partitionHeight, err := strconv.Atoi(params.Get("partitionHeight"))
	if err != nil {
		result = multierror.Append(result, err)
	}

	if result != nil {
		log.Errorf("error validating request params: %v", result)
		http.Error(w, result.Error(), http.StatusBadRequest)
	}

	sG, err := sudoku.GenerateSudokuGrid(size, partitionWidth, partitionHeight)
	if err != nil {
		log.Errorf("error generating sudoku grid: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sG.SetGridToLevel(level)
	if err != nil {
		log.Errorf("error setting the grid to the difficulty level: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res []byte

	if pretty == "true" {
		res = []byte(sG.ToStringPrettify())
	} else {
		res, err = json.Marshal(sG)
		if err != nil {
			log.Errorf("error marshalling the response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write(res)
}

func sudokuSolverHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pretty := params.Get("pretty")
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error reading the body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sG := sudoku.SudokuGrid{}
	err = json.Unmarshal(body, &sG)
	if err != nil {
		log.Errorf("error unmarshalling the body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(sG)
	if err != nil {
		log.Warnf("error marshalling body: %s", string(b))
	}

	log.Debugf("body: %v", string(b))

	if err = sG.Valid(); err != nil {
		log.Errorf("error validating the sudoku grid: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = sG.Solve(); err != nil {
		log.Errorf("error solving the sudoku puzzle: %v", err)
		w.Write([]byte(fmt.Sprintf("error solving the sudoku puzzle: %v", err)))
		return
	}

	var res []byte

	if pretty == "true" {
		w.Header().Set("Content-Type", "plain/text")
		res = []byte(sG.ToStringPrettify())
	} else {
		w.Header().Set("Content-Type", "application/json")
		res, err = json.Marshal(sG)
		if err != nil {
			log.Errorf("error marshalling the response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Write(res)
}

func SetupHandlers(r *mux.Router) {
	publicMiddleware := []middleware.Middleware{
		middleware.LogMiddleware,
	}
	// TODO: add support for authentication => privateMiddleware
	r.HandleFunc("/", middleware.Chain(homeHandler, publicMiddleware...)).Methods("GET")
	r.HandleFunc("/sudoku", middleware.Chain(sudokuSolverHandler, publicMiddleware...)).Methods("POST")
	r.HandleFunc("/sudoku", middleware.Chain(sudokuGeneratorHandler, publicMiddleware...)).Methods("GET")
}

func StartServer(cfg config.Config) {
	r := mux.NewRouter()
	SetupHandlers(r)

	initLogger(cfg)

	log.Printf("Server listening on port %d", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
