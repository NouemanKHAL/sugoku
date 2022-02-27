# Sudoku Solver REST API

A simple sudoku solver REST API that handles sudoku grids with custom grid and subgrid sizes.

## Prerequisites

- Go 1.16 or later. [See the install instructions for Go](https://go.dev/doc/install).

## Usage

To run the server:

- run `go build` in the root folder.
- start the server by running the generated binary.
```console
./sudoku-solver-rest-api
```

This will launch a server that listens on the default port `7007`. 

If you wish to use a different port number, make sure to pass the environment variable `SUDOKU_SERVER_PORT` as follows:

```console
SUDOKU_SERVER_PORT=7331 ./sudoku-solver-rest-api
```

### Example: Solve a 4x4 sudoku grid with 2x2 subgrids

We will use the grid below as input for the example.

```console
 1  . | .  4
 .  . | 1  .
--------------
 2  . | .  .
 4  . | 2  .
```

which corresponds to the following `json` payload:

 ```json
 {
    "Size": 4,
    "PartitionWidth": 2,
    "PartitionHeight": 2,
    "Grid":[
        [49,46,46,52],
        [46,46,49,46],
        [50,46,46,46],
        [52,46,50,46]
    ]
 }
```

- Send a POST request to `/sudoku` with the previous `json` payload in the body.

*You may use `?pretty=true` for a human readable output.*
```console
curl -X POST http://localhost:7007/sudoku?pretty=true -d {\"Size\":4,\"PartitionWidth\":2,\"PartitionHeight\":2,\"Grid\":[[49,46,46,52],[46,46,49,46],[50,46,46,46],[52,46,50,46]]}
```

- Server responds with:
```console
 1  2 | 3  4
 3  4 | 1  2
--------------
 2  1 | 4  3
 4  3 | 2  1
```

- Done!

## TO DO

- Add more unit tests
- ~~Improve errors handling in the route handlers~~
- ~~Improve request logging~~
- Add a Sudoku puzzle generator endpoint
- Add authentication?
- Support HTTPS?
- Build a Front End
- ...
