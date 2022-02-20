# Sudoku Solver REST API

A simple sudoku solver REST API that handles sudoku grids with custom grid size and subgrid size.

## Prerequisites

- Go 1.16 or later. [See the install instructions for Go](https://go.dev/doc/install).

## Usage

To run the server:

- run `go build` in the root folder.
- start the server 
```console
./sudoku-solver-rest-api
```

This will launch a server that listens on the default port `7007`, if you wish to use a different port number, make use of the environment variable `SUDOKU_SERVER_PORT`, e.g.:

```console
SUDOKU_SERVER_PORT=7500 ./sudoku-solver-rest-api
```

