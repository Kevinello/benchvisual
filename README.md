# go-benchvisual

Parse and visualize Golang standard Benchmark output.
mainly based on `regexp2 + go-echart + cobra`

## Motivation

go-benchvisual is a tool for visualizing Golang standard Benchmark output, which will show benchmark metrics in different charts, group by benchmark scenario, make it convenient for user to compare same metric of different benchmark target in a specific scenario

## Features

- piped output of `go test -bench` as input
- file as input
- custom regexp / separator for Benchmark name to recognize "target" and "scenario"
- custom output file path
- json output instead of visualized output for secondary development

## Install

```shell
go install github.com/Kevinello/go-benchvisual@latest
```

## Usage

```shell
go-benchvisual --help
Parse and visualize Golang standard Benchmark output.
go-benchvisual provides pipe mode and file mode, it will work in pipe mode in default, add flag -f to let it work in file mode.
there are two definitions in go-benchvisual -- "targets" and "scenarios", which the visualization works depends on.
"targets" means targets to generate Benchmark on, and the "scenarios" means specific Benchmark scenarios where the Benchmark run in.
In the visualization progress, we will visualize Benchmark in a concepts mapping below:
        - bench.Set(package) -> page(html)
        - metrics(ns/op...)  -> series of metrics value
        - targets            -> series name(x axis)
        - scenarios          -> dummy values in charts(group name)
go-benchvisual also provides json output format for your secondary development, use --json to let it output json file.

Usage:
  go-benchvisual [--version] [--help] [-s <separator> | -r <regexp>] [-f <benchmark path>] [-o <output path>] [--json] [flags]

Examples:
  go test -bench . | go-benchvisual -r '^Bench(mark)?(?<target>\\S+)/(?<scenario>\\S+)$'
  go-benchvisual -s '/' -f "path/to/origin/benchmark/file"

Flags:
  -f, --file string     use file mode instead of pipe mode, Read the original Benchmark output from the given file path
  -h, --help            help for go-benchvisual
      --json            only output parsed Benchmark result in json file
  -o, --output string   directory path to save the output file (default ".")
  -r, --regex string    regexp expression with two sub groups(target and scenario), written in '.NET-style capture groups'--(?<name>re) or (?'name're).
                        e.g., '^Bench(mark)?(?<target>\S+/\S+)/(?<scenario>\S+)$' (default "^Bench(mark)?(?<target>[A-Z]+\\S*)(?<scenario>[A-Z]+\\S*)$")
  -s, --sep string      string separator of a Benchmark string's target and scenario.
                        e.g., we got a benchmark name string 'BenchmarkFibonacci/100times' with separator '/', then the target of it is 'Fibonacci' and the scenario of it is '100times'.
                        
  -v, --version         version for go-benchvisual
```

## TODO

- more custom configs for visualization

## Contribution & Support

Feel free to send a pull request if you consider there's something which can be improved. Also, please open up an issue if you run into a problem when using go-benchvisual or just have a question.
