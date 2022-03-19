<p align="center">
  <img src="https://user-images.githubusercontent.com/25211181/156906917-fab62386-0f7b-4d8a-b004-6d8c6a2ebc1d.png"> 
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

## Examples

### Solve a sudoku puzzle

We will use the grid below as input for the example.

```console
 1  . | .  4
 .  . | 1  .
--------------
 2  . | .  .
 4  . | 2  .
```

In order to solve a sudoku puzzle:

1. Send a POST request to `/sudoku` with the puzzle in `json` format in the body, and optionally you may set the query parameter `pretty=true` for a human readable output.

```console
curl -X POST http://localhost:7007/sudoku?pretty=true -d '{"size":4,"partitionWidth":2,"partitionHeight":2,"grid":[[49,46,46,52],[46,46,49,46],[50,46,46,46],[52,46,50,46]]}'
```

2. Server responds with a valid solution:

```console
 1  2 | 3  4
 3  4 | 1  2
--------------
 2  1 | 4  3
 4  3 | 2  1
```

3. Done!

### Generate a sudoku puzzle

In order to generate a 9x9 hard sudoku puzzle:

1. Send a GET Request to `/sudoku` endpoint with the following query parameters `size=9`, `partitionWidth=3`, `partitionHeight=3` and optionally add `pretty=true` for a human readable output

```console
curl 'http://localhost:7007/sudoku?pretty=true&size=9&partitionWidth=3&partitionHeight=3&level=hard'
```

2. Server responds with a human readable output of the puzzle


```console
 1  .  7 | .  .  . | .  .  .
 .  .  . | 3  8  . | .  .  .
 3  .  . | .  9  . | .  .  5
-----------------------------
 .  1  . | .  .  2 | .  .  .
 .  .  2 | .  .  . | .  .  .
 .  .  . | .  1  . | 7  .  .
-----------------------------
 2  .  1 | .  .  . | .  .  .
 5  .  . | .  3  . | .  .  .
 6  .  . | .  .  1 | .  .  .
```

3. Done!

## TO DO

- Add more unit tests
- ~~Improve errors handling in the route handlers~~
- ~~Improve request logging~~
- ~~Add a Sudoku puzzle generator endpoint~~
- Add authentication?
- Support HTTPS?
- Build a Front End
- ...
