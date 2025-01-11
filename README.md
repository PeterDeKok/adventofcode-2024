# Advent of Code 2024

Archive of puzzle solutions and tui for managing the [Advent of Code 2024](https://adventofcode.com/2024) challenges.

No problem defitions, personal input or solutions will be archived in this repo.

## TLDR
The puzzle solutions can be found at [/src/puzzles](/src/puzzles)

## Result summary

| day        | stars | part 1       | part 2       |
|------------|-------|--------------|--------------|
| 2024-12-01 | 2     | 410.344µs    | 552.994µs    |
| 2024-12-02 | 2     | 358.328µs    | 754.861µs    |
| 2024-12-03 | 2     | 414.663µs    | 648.858µs    |
| 2024-12-04 | 2     | 27.262693ms  | 32.280411ms  |
| 2024-12-05 | 2     | 1.530678ms   | 2.408762ms   |
| 2024-12-06 | 1     | 2.122639ms   | -            |
| 2024-12-07 | 2     | 330.376587ms | 21.457967ms  |
| 2024-12-08 | 2     | 402.943µs    | 361.918µs    |
| 2024-12-09 | 2     | 7.57408ms    | 232.337639ms |
| 2024-12-10 | 2     | 2.037859ms   | 259.54µs     |
| 2024-12-11 | 2     | 11.221525ms  | 639.067508ms |
| 2024-12-12 | 2     | 1.335515ms   | 3.077066ms   |
| 2024-12-13 | 0     | -            | -            |
| 2024-12-14 | 0     | -            | -            |
| 2024-12-15 | 0     | -            | -            |
| 2024-12-16 | 1     | 2.727663ms   | -            |
| 2024-12-17 | 0     | -            | -            |
| 2024-12-18 | 2     | 557.447µs    | 997.415593ms |
| 2024-12-19 | 0     | -            | -            |
| 2024-12-20 | 0     | -            | -            |
| 2024-12-21 | 0     | -            | -            |
| 2024-12-22 | 0     | -            | -            |
| 2024-12-23 | 0     | -            | -            |
| 2024-12-24 | 0     | -            | -            |
| 2024-12-25 | 0     | -            | -            |

Cumulative runtime
```
2.318956081s
```

## Tools

This year's solutions and tools are written in [Go](https://go.dev/).

This repo contains a TUI for managing the challenge and viewing statistics.

Input files, test input, problem statements, boilerplate challenge runners, etc. can be created when available.

### TUI
A compiled binary is not supplied.
A compiled binary is currenly untested.
A properly set environment is expected.

The TUI can be started with:
```
go run ./src
```

The TUI enables the alternate screen buffer,
therefore normal logging to stdout is not available.
For debugging, multiple log files are available.

The main log file set by the [environment](#environment) and
part specific log files for the individual puzzle solutions build and
run stages. Which can be found in the output directory
of every puzzle part after building, running, etc.

#### Environment
| VAR                        | Default | Required | Description                                                                                                                                              |
|----------------------------|---------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| AOC_LOG_FILE               | -       | Yes      | The filepath to the log file. it is advisible to use an absolute path, however it will load relative paths as well. the parent directory should exist.   |
| AOC_PUZZLES_DIR            | -       | Yes      | The filepath to the puzzles directory. it is advisible to use an absolute path, however it will load relative paths as well. the directory should exist. |
| AOC_SESSION_COOKIE_EXPIRES | -       | Yes      | The expiration date of the aoc session cookie.                                                                                                           |
| AOC_SESSION_COOKIE_VALUE   | -       | Yes      | the content of the aoc session cookie.                                                                                                                   |

### Remote
Retrieving and parsing remote data is provided.
The creator and maintainer of the [Advent of Code](https://adventofcode.com/)
challenges has asked users to limit their requests to a decent rate.
Therefore this tool also includes base ratelimits.
The default values for each individual remote operation have been chosen to provide
a balance between user experience.
Wherever possible, the published ratelimits are adhered to, however not all endpoints have defined ratelimits.

### Testing
Some tests are provided. But at time of writing, the project is under-tested.

Utility functions are mostly tested, as they can be used in solutions.
The manager and tui contain some tests, but these require a lot more updates.

```
go test ./...
```

Please note that the puzzle part boilerplate also contains test scaffolding.
The above command will run all tests, including these.
Invalid or incomplete solutions will therefore also fail these tests.

Depending on the solution efficiency, these tests might massively
increase test runtime. To the point where the go tools default timeout
might be reached (10 mminutes).

## Solutions
After the TUI (or tools [1](#appendix-1---ideas)) have created the
boilerplate fo puzzle parts, the solution should be created in the `solution.go` file.
Other files can be created and imported.

The `solution_test.go` file can be used to verify the solution in a standalone manner.
The default boilerplate requires the input and sample input & expected
files to be filled.

The TUI or tools can retrieve these values, or they can be copied manually.
If these files contain additional new lines after the official content,
these need to be handled in the solution explicitly. Any trailing new lines in
the expected sample file can lead to unexpected results as the solution returns
the same format for the sample and actual inputs.

# Appendix
## Appendix 1 - Ideas
See [IDEAS.md](IDEAS.md) for any future ideas and in progress or upcomming tasks.
