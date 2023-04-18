// Package bench parse Golang standard benchmark output
//
//	@update 2023-03-07 12:06:22
package bench

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Kevinello/benchvisual/internal/collections"
	"github.com/charmbracelet/log"
	"github.com/dlclark/regexp2"
)

// Set is a set of benchmark runs
type Set struct {
	Goos    string                   `json:"goos,omitempty"`
	Goarch  string                   `json:"goarch,omitempty"`
	Pkg     string                   `json:"pkg,omitempty"`
	CPU     string                   `json:"cpu,omitempty"`
	Targets map[string]BenchmarkList `json:"targets,omitempty"` // map[target][]Benchmark; group of Benchmark result(Series in visualized result)
}

// Benchmark is an individual run. Note that all metrics in here must be represented as
// a float type, even if Go only emits integer values, so that in checks we can correctly
// evaluate divisions so that results come out as floats instead of being truncated to
// integers.
type Benchmark struct {
	Name     string `json:"name,omitempty"`
	CPUCores int    `json:"cpu_cores,omitempty"`
	Runs     int    `json:"runs,omitempty"` // benchmark times

	// For a Benchmark function BenchXXX10000, its target is XXX, and its Scenario is 10000
	// The Benchmark of different target is compared in each Scenario
	Target   string `json:"target,omitempty"`
	Scenario string `json:"scenario,omitempty"`

	NsPerOp       float64            `json:"ns_per_op,omitempty"`
	Mem           Mem                `json:"mem,omitempty"`            // metrics from '-benchmem'
	CustomMetrics map[string]float64 `json:"custom_metrics,omitempty"` // custom metrics(https://tip.golang.org/pkg/testing/#B.ReportMetric)

	ReachBaseline bool `json:"reach_baseline"` // whether this benchmark reach baseline
}

// BenchmarkList implement sort.Interface
//
//	@author kevineluo
//	@update 2023-03-07 03:21:11
type BenchmarkList []Benchmark

func (b BenchmarkList) Len() int {
	return len(b)
}

func (b BenchmarkList) Less(i, j int) bool {
	return b[i].Scenario < b[j].Scenario
}

func (b BenchmarkList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Mem is memory allocation information about a run
type Mem struct {
	BytesPerOp  float64 `json:"bytes_per_op,omitempty"`
	AllocsPerOp float64 `json:"allocs_per_op,omitempty"`
	MBPerSec    float64 `json:"mb_per_sec,omitempty"`
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
func ParseSet(reader *bufio.Reader, sep string, regex *regexp2.Regexp) (set *Set, err error) {
	set = &Set{
		Targets: make(map[string]BenchmarkList),
	}

	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			err = fmt.Errorf("found EOF before 'PASS' or 'FAIL'(the end of a Benchmark set)")
			return nil, err
		} else if err != nil {
			return nil, err
		}
		line := string(l)

		if strings.HasPrefix(line, "PASS") || strings.HasPrefix(line, "FAIL") {
			// end of one set
			log.Info("Benchmark set parsed")
			break
		} else if os, found := strings.CutPrefix(line, "goos: "); found {
			set.Goos = os
			log.Info("Benchmark metadata", "goos", set.Goos)
		} else if arch, found := strings.CutPrefix(line, "goarch: "); found {
			set.Goarch = arch
			log.Info("Benchmark metadata", "goarch", set.Goarch)
		} else if pkg, found := strings.CutPrefix(line, "pkg: "); found {
			set.Pkg = pkg
			log.Info("Benchmark metadata", "pkg", set.Pkg)
		} else if cpu, found := strings.CutPrefix(line, "cpu: "); found {
			set.CPU = cpu
			log.Info("Benchmark metadata", "cpu", set.CPU)
		} else if strings.HasPrefix(line, "Bench") {
			log.Debug("[ParseSet] Benchmark line", "origin_line", line)
			bench, err := ParseBench(line, sep, regex)
			if err != nil {
				return nil, fmt.Errorf("%w: %q", err, line)
			}
			log.Debug("Benchmark parsed", "name", bench.Name, "runs", bench.Runs, "target", bench.Target, "scenario", bench.Scenario)
			if benchmarks, ok := set.Targets[bench.Target]; ok {
				set.Targets[bench.Target] = append(benchmarks, *bench)
			} else {
				set.Targets[bench.Target] = []Benchmark{*bench}
			}
		}
	}

	return
}

// GetScenarios get all unique scenario in a Benchmark set
//
//	@receiver set *Set
//	@return scenarios []string
//	@author kevineluo
//	@update 2023-03-07 04:33:12
func (set *Set) GetScenarios() (scenarios []string) {
	scenarioSet := collections.NewSet[string](0)
	for _, benchmarks := range set.Targets {
		// collect unique scenario
		scenarioSet = scenarioSet.Union(collections.SliceToSet(collections.Map(benchmarks, func(benchmark Benchmark) (scenario string) { return benchmark.Scenario })))
	}
	scenarios = scenarioSet.ToSlice()
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
func ParseBench(line string, sep string, regex *regexp2.Regexp) (bench *Benchmark, err error) {
	bench = new(Benchmark)
	// split out name
	split := collections.Map(strings.Split(line, "\t"), func(s string) string {
		return strings.TrimSpace(s)
	})
	bench.Name, split = popLeft(split)

	// parse cpu core nums
	sepIdx := strings.LastIndex(bench.Name, "-")
	if sepIdx == -1 {
		return nil, fmt.Errorf("[ParseBench] invalid benchmark name, no '-' found")
	}
	bench.CPUCores, err = strconv.Atoi(bench.Name[sepIdx+1:])
	if err != nil {
		return nil, fmt.Errorf("[ParseBench] invalid benchmark name, failed to parse cpu core nums")
	}
	bench.Name = bench.Name[:sepIdx]

	if regex != nil {
		// with regexp
		log.Debug("[ParseBench] in regex mode", "regexp", regex.String())
		match, err := regex.FindStringMatch(bench.Name)
		if err != nil {
			return nil, fmt.Errorf("[ParseBench] error when parse [benchmark name: %s], [regexp: %s], error: %w", bench.Name, regex.String(), err)
		}
		if match == nil {
			return nil, fmt.Errorf("[ParseBench] no match found in [benchmark name: %s], [regexp: %s]", bench.Name, regex.String())
		}
		if group := match.GroupByName("target"); group != nil {
			bench.Target = group.String()
		} else {
			return nil, fmt.Errorf("[ParseBench] group 'target' not found in match result")
		}
		if group := match.GroupByName("scenario"); group != nil {
			bench.Scenario = group.String()
		} else {
			return nil, fmt.Errorf("[ParseBench] group 'scenario' not found in match result")
		}
	} else if sep != "" {
		// with separator
		log.Debug("[ParseBench] in separator mode", "separator", sep)
		var after string
		var found bool
		// Compatible for "Benchmark" and "Bench"
		if after, found = strings.CutPrefix(bench.Name, "Benchmark"); !found {
			if after, found = strings.CutPrefix(bench.Name, "Bench"); !found {
				return nil, fmt.Errorf("[ParseBench] illegal Benchmark name: %s", bench.Name)
			}
		}
		bench.Target, bench.Scenario, found = strings.Cut(after, sep)
		if !found {
			return nil, fmt.Errorf("[ParseBench] given separator[%s] not found in Benchmark name: %s", sep, bench.Name)
		}
	} else {
		return nil, fmt.Errorf("neither given regexp expression(-regex) nor given separator(--sep)")
	}

	// parse runs (doesn't include units)
	tmp, split := popLeft(split)
	if bench.Runs, err = strconv.Atoi(tmp); err != nil {
		return nil, fmt.Errorf("[ParseBench] %s: could not parse run: %w (line: %s)", bench.Name, err, line)
	}

	// parse metrics with units
	for len(split) > 0 {
		tmp, split = popLeft(split)
		valueAndUnits := strings.Split(tmp, " ")
		if len(valueAndUnits) < 2 {
			return nil, fmt.Errorf("[ParseBench] expected two parts in value '%s', got %d", tmp, len(valueAndUnits))
		}

		value, units := valueAndUnits[0], valueAndUnits[1]
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
			return nil, fmt.Errorf("[ParseBench] %s: could not parse %s: %v", bench.Name, units, err)
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
