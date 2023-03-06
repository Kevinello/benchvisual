// Package bench parse Golang standard benchmark output
//
//	@update 2023-03-07 12:06:22
package bench

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

// Set is a set of benchmark runs
type Set struct {
	Goos   string
	Goarch string
	Pkg    string
	Groups map[string][]Benchmark `json:"groups,omitempty"` // map[regexp][]Benchmark; group of Benchmark result(for visualizing result in group)
}

// Benchmark is an individual run. Note that all metrics in here must be represented as
// a float type, even if Go only emits integer values, so that in checks we can correctly
// evaluate divisions so that results come out as floats instead of being truncated to
// integers.
type Benchmark struct {
	Name string `json:"name,omitempty"`
	Runs int    `json:"runs,omitempty"` // benchmark times

	NsPerOp       float64            `json:"ns_per_op,omitempty"`
	Mem           Mem                `json:"mem,omitempty"` // metrics from '-benchmem'
	CustomMetrics map[string]float64 `json:",omitempty"`    // https://tip.golang.org/pkg/testing/#B.ReportMetric
}

// Mem is memory allocation information about a run
type Mem struct {
	BytesPerOp  float64
	AllocsPerOp float64
	MBPerSec    float64
}

// ParseSet Parse one set of benchmark output
//
//	@param reader LineReader
//	@param firstLine string
//	@param groupRegexps ...*regexp2.Regexp regexps for identify group of benchmark
//	@return set *Set
//	@return err error
//	@author kevineluo
//	@update 2023-03-07 12:16:30
func ParseSet(reader *bufio.Reader, groupRegexps ...*regexp2.Regexp) (set *Set, err error) {
	set = &Set{
		Groups: make(map[string][]Benchmark),
	}

	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line := string(l)

		if strings.HasPrefix(line, "PASS") || strings.HasPrefix(line, "FAIL") {
			// end of one set
			break
		} else if os, found := strings.CutPrefix(line, "goos: "); found {
			set.Goos = os
		} else if arch, found := strings.CutPrefix(line, "goarch: "); found {
			set.Goarch = arch
		} else if pkg, found := strings.CutPrefix(line, "pkg: "); found {
			set.Pkg = pkg
		} else if strings.HasPrefix(line, "Benchmark") {
			bench, err := ParseBench(line)
			if err != nil {
				return nil, fmt.Errorf("%w: %q", err, line)
			}
			// default group
			group := "others"
			for _, groupRegexp := range groupRegexps {
				if match, _ := groupRegexp.MatchString(bench.Name); match {
					group = groupRegexp.String()
					break
				}
			}
			if benchmarks, ok := set.Groups[group]; ok {
				set.Groups[group] = append(benchmarks, *bench)
			} else {
				set.Groups[group] = []Benchmark{*bench}
			}
		}
	}

	return
}

// ParseBench parses a single line from a benchmark.
//
// Benchmarks take the following format:
//
//	BenchmarkXXX	300000	5160 ns/op	5408 B/op	69 allocs/op
//	@param line string
//	@return bench *Benchmark
//	@return err error
//	@author kevineluo
//	@update 2023-03-07 12:11:18
func ParseBench(line string) (bench *Benchmark, err error) {
	bench = new(Benchmark)
	// split out name
	split := strings.Split(line, "\t")
	bench.Name, split = popLeft(split)

	// runs - doesn't include units
	tmp, split := popLeft(split)
	if bench.Runs, err = strconv.Atoi(tmp); err != nil {
		return nil, fmt.Errorf("%s: could not parse run: %w (line: %s)", bench.Name, err, line)
	}

	// parse metrics with units
	for len(split) > 0 {
		tmp, split = popLeft(split)
		valueAndUnits := strings.Split(tmp, " ")
		if len(valueAndUnits) < 2 {
			return nil, fmt.Errorf("expected two parts in value '%s', got %d", tmp, len(valueAndUnits))
		}

		var value, units = valueAndUnits[0], valueAndUnits[1]
		switch units {
		case "ns/op":
			bench.NsPerOp, err = strconv.ParseFloat(value, 64)
		case "B/op":
			bench.Mem.BytesPerOp, err = strconv.ParseFloat(value, 64)
		case "allocs/op":
			bench.Mem.AllocsPerOp, err = strconv.ParseFloat(value, 64)
		case "MB/s":
			bench.Mem.MBPerSec, err = strconv.ParseFloat(value, 64)
		default:
			if bench.CustomMetrics == nil {
				bench.CustomMetrics = make(map[string]float64)
			}
			bench.CustomMetrics[units], err = strconv.ParseFloat(value, 64)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: could not parse %s: %v", bench.Name, units, err)
		}
	}

	return
}

// popLeft pops the first string off a slice
//
//	@param s []string
//	@return string
//	@return []string
//	@author kevineluo
//	@update 2023-03-06 08:01:06
func popLeft(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return strings.Trim(s[0], " "), s[1:]
}
