# One Billion Data Parsing Project

## Table of Contents

- [Introduction](#introduction)
- [How Much is a Billion?](#how-much-is-a-billion)
- [The Problem Statement](#the-problem-statement)
- [Hardware Constraints](#hardware-constraints)
- [Software Constraints](#software-constraints)
- [Configuration](#configuration)
- [How to Run the Project](#how-to-run-the-project)
  - [Directory Structure](#directory-structure)
  - [Environment Variables](#environment-variables)
  - [Running the Project](#running-the-project)
- [Criteria for Success](#criteria-for-success)
- [The Solution](#the-solution)
  - [Stage 1: Base Parsing Solution](#stage-1-base-parsing-solution)
    - [Process](#process)
    - [Points considered](#points-considered)
    - [Results](#results)
- [Resume](#resume)
- [References](#references)
- [License](#license)

## Introduction

As Jockie Stewart once said, "You don't have to be an engineer to be a racing driver, but you do have to have mechanical sympathy." This project focuses on the parsing of a CSV file containing one billion rows of data. It's not just about writing code that works; it's about designing software that respects the hardware and the necessity of scale.

## How Much is a Billion?

A billion is a massive number. To count from 1 to 1 billion, it would take you 31 years, 259 days, 1 hour, 46 minutes, and 40 seconds. If you were to stack a billion pennies, it would reach a height of 870 miles. Driving a billion miles would allow you to circumnavigate the Earth 40,000 times. Clearly, a billion is a significant figure.

This scale is often required by large companies like Google, which has 1.2 billion users. They must store and process vast amounts of data from these users. Dealing with this scale presents a considerable challenge, and this project serves as a small example of how to tackle it.

## The Problem Statement

The problem is straightforward: you have a CSV file with one billion rows, and your task is to parse it. The CSV file follows the format:

```html
<station_name:string>;<temperature:float>
```

You need to parse the file into a list where each row represents a station with its minimum temperature, maximum temperature, and average temperature. The list should be sorted by the station name and have the following format:

```html
<station_name:string>;<min_temperature:float>;<max_temperature:float>;<medium_temperature:float>;<count_station:int>
```

## Hardware Constraints

Here are the hardware constraints for this project:

| Hardware Component | Constraint |
|--------------------|------------|
| CPU                | 8          |
| RAM                | 32GB       |
| Storage            | 1TB SSD    |
| Architecture       | x86_64     |
| Goroutines         | 1          |
| GoMaxProcs         | 8          |

## Software Constraints

Here are the software constraints for this project:

| Software Component | Constraint    |
|--------------------|---------------|
| Go Version         | 1.22.5        |
| Go Arch            | amd64         |
| File System        | ext4          |
| Operating System   | Linux         |
| Number of Files    | 1             |
| File Size          | 16GeGa        |
| File Format        | CSV           |
| File Columns       | 2             |
| File Rows          | 1,000,010,000 |
| File Encoding      | UTF-8         |
| File Delimiter     | `;`           |
| File Compression   | None          |
| File Line Ending   | `\n`          |

## Configuration

we use the following environment variables to configure the algorithm:

- `INPUT_FILE_PATH`: The path to the input data csv file to parse.
- `OUTPUT_FILE_PATH`: The path to the output data csv file to write the parsed data.
- `METRICS_FILE_PATH`: The path to the metrics data csv file to write the metrics data.
- `MAX_WORKERS`: The maximum number of workers to use for parsing the data.
- `CHUNK_SIZE`: The size of the chunk to read from the input file.
- `LOG_LEVEL`: The log level to use for the logger.
- `VERSION`: The version of the algorithm.
- `NUMBER_OF_ROWS`: The number of rows in the input data file.

## How to Run the Project

### Directory Structure

The project has the following directory structure:

```bash
.
.
├── cmd
│   ├── faker
│   │   └── facker.go
│   └── version
│       └── v1_base
│           └── main.go
├── config
│   └── config.go
├── data
│   ├── metrics.log
│   ├── output.csv
│   ├── output_test.csv
│   ├── weather_data.csv
│   ├── weather_data_test.csv
├── go.mod
├── go.sum
├── internal
│   ├── context
│   │   └── context.go
│   ├── hashutil
│   │   ├── hashutil.go
│   │   └── hashutil_test.go
│   └── tracker
│       └── tracker.go
├── LICENSE
├── main.go
├── main_test.go
├── readme.md
└── scripts
    ├── update_and_run.sh
    └── weater_data.sh
```

- All the helper code is in the `internal` directory. The `cmd` directory contains the main entry points for the project. The `config` directory contains the configuration for the project. The `data` directory contains the input and output data files. The `scripts` directory contains the bash scripts to run the project.
- Every version of the project is in a separate branch, and on a separate directory in the `cmd/version` directory.

### Environment Variables

The project uses environment variables to configure the software. You can set the environment variables in a `.env` file in the root directory of the project. Here is an example of the `.env` file:

```bash
INPUT_FILE_PATH = './data/weather_data.csv'
OUTPUT_FILE_PATH = './data/output.csv'
METRICS_FILE_PATH = 'data/metrics.log'
LOG_LEVEL = 'INFO'
MAX_WORKERS = 10
CHUNK_SIZE = 1000
VERSION = '1.0.0'
NUMBER_OF_ROWS = 1000000000
```

### Running the Project

To run the project, their is a bash script `update_and_run.sh` that you can use on a Linux machine. The script update the go version to 1.22.5, update the go modules, and run the project.

```bash
chmod +x ./scripts/update_and_run.sh
./scripts/update_and_run.sh
```

## Measurements and Metrics

The project will measure the following metrics:

- Execution time: The time taken to parse the input data file. max 10 minutes before the project is considered a failure.
- Memory usage: The amount of memory used by the software. max 24GB before the project is considered a failure.

To measure these metrics, we implemented the function `measure` that will write the metrics to a file, stop the process if any of the metrics exceed the constraints, and print the metrics to the console.

Check the [`tracker` module](https://github.com/aminesayagh/1-Billion-row/blob/stage_1_simple_implementation/internal/tracker/tracker.go) in the directory `internal/tracker/tracker.go` for more details.

## Measurements and Metrics on testing

On testing face of the project, the following time metrics inspired from `MCDA` (multi-criteria decision analysis) has been tracked on the comparison between different strategy of parsing the data:

- Median Execution Time: The median time taken to execute the strategy.
- Percentile 90 Execution Time: The time taken to execute the strategy for 90% of the data.
- Minimum Execution Time: The minimum time taken to execute the strategy.
- Maximum Execution Time: The maximum time taken to execute the strategy.

## Criteria for Success

The project will be considered successful if it meets the following criteria:

- The solution parses the input data file correctly.
- The solution writes the output data file correctly.
- The solution respects the hardware constraints.
- The solution respects the solution constraints.
- The solution is well-documented.

## The Solution

To achieve the best performance, we separate the stages of the development of our parsing solution on the following steps (separate branches `<stage_{number}_{title}>`):

### Stage 1: Base Parsing Solution

In this stage, the parsing solution is a simple one threaded solution, we took on consideration this following points to achieve the best performance:

#### Process

1. Read the input file in chunks.
2. Parse the data in each chunk.
3. Store the parsed data in memory (in a map).
4. Write the output data to a file.

#### Points considered

- Use of pointers to avoid copying data.
- Use of the `bufio` package to read the file in chunks, streamlining the reading process, reducing I/O operations.
- Use `strings.Index` Instead of `strings.Split` to avoid memory allocation.
- Use maps to store the parsed data, with a size of 100,000 records, referencing the estimated number of stations in the data.
- Increase the Buffer size for the `bufio.Scanner` to 64 * 1024 bytes.
- Process file on a byte level, to avoid memory allocation.
- Compare the performance of the generation of map keys based on the station name, between the use of `strings` and `bytes` of 16 bytes, and a Hash (FNV-1a), based on multi-criteria decision analysis, the choice was to use the `string` type.
- Cache the temperature values generated on float64 to avoid conversion on each iteration.

#### Results

- Execution time: 2m40.36641632s.
- Allocated memory: 46.02 MB.
- Total memory Allocated: 16724.07 MB.
- System memory used: 67.19 MB.
- Heap memory used: 46.02 MB.

### Stage 2: Assembly Parsing Solution

#### Stage 2: Points considered

- Use of Assembly to parse the data.
- Structure the implementation of a Automat of float processing (see the `cmd/version/v2_assembly/automator/README.md` readme).
- (X) Use of SIMD (Single Instruction, Multiple Data) instructions to parse the data (abandoned due to the complexity of the implementation).

## Resume

| Stage                   | Execution Time    | Total Memory Allocated |
|-------------------------|-------------------|------------------------|
| Stage 1: Simple Parsing | 2m40.36641632s    | 16724.07 MB            |

## References

- Go Programming Language Documentation
- MCDA Techniques
- FNV-1a Hash Algorithm

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
