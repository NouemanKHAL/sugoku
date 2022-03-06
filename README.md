<p align="center">
  <img src="https://user-images.githubusercontent.com/25211181/156906917-fab62386-0f7b-4d8a-b004-6d8c6a2ebc1d.png" width="900"> 
</p>

A simple Sudoku REST API (written in Go) designed to solve and generate Sudoku puzzles with customizable grid dimensions as well as subgrid dimensions.

## Prerequisites

- Go 1.16 or later. [See the install instructions for Go](https://go.dev/doc/install).

## Usage
```console
Usage:
  sugoku [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Start a sudoku server

Flags:
  -h, --help   help for sugoku

Use "sugoku [command] --help" for more information about a command.
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
