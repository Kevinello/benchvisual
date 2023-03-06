package bench

import (
	"bufio"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/smartystreets/goconvey/convey"
)

var (
	benchmarkOutputs = []string{
		`goos: darwin
goarch: amd64
pkg: go.bobheadxi.dev/gobenchdata/demo
cpu: Intel AMD Xeon Phenom
BenchmarkFib10/Fib()-12	3293298	330 ns/op
BenchmarkFib100/Fib()-12	329329	329 ns/op
BenchmarkPizzas10/Pizzas()-12	25820055	50.0 ns/op	3.00 pizzas
BenchmarkPizzas100/Pizzas()-12	2582005	50.0 ns/op	3.00 pizzas
PASS`,
		`goos: darwin
goarch: amd64
pkg: go.bobheadxi.dev/gobenchdata/demo
BenchmarkFib10/FibSlow()-12	3033732	358 ns/op	16 B/op	1 allocs/op
BenchmarkFib100/FibSlow()-12	303373	358 ns/op	16 B/op	1 allocs/op
BenchmarkPizzas10/Pizzas()-12	22866814	46.3 ns/op	9.00 pizzas	0 B/op	0 allocs/op
BenchmarkPizzas100/Pizzas()-12	2286681	46.3 ns/op	9.00 pizzas	0 B/op	0 allocs/op
PASS`,
	}
	targetSets = []Set{
		{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Groups: map[string][]Benchmark{
				"others": {
					{
						Name: "BenchmarkFib10/Fib()-12", Runs: 3293298, NsPerOp: 330,
					},
					{
						Name: "BenchmarkFib100/Fib()-12", Runs: 329329, NsPerOp: 329,
					},
					{
						Name: "BenchmarkPizzas10/Pizzas()-12", Runs: 25820055, NsPerOp: 50, CustomMetrics: map[string]float64{"pizzas": 3.00},
					},
					{
						Name: "BenchmarkPizzas100/Pizzas()-12", Runs: 2582005, NsPerOp: 50, CustomMetrics: map[string]float64{"pizzas": 3.00},
					},
				},
			},
		},
		{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Groups: map[string][]Benchmark{
				"others": {
					{
						Name: "BenchmarkFib10/FibSlow()-12", Runs: 3033732, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1},
					},
					{
						Name: "BenchmarkFib100/FibSlow()-12", Runs: 303373, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1},
					},
					{
						Name: "BenchmarkPizzas10/Pizzas()-12", Runs: 22866814, NsPerOp: 46.3, CustomMetrics: map[string]float64{"pizzas": 9.00},
					},
					{
						Name: "BenchmarkPizzas100/Pizzas()-12", Runs: 2286681, NsPerOp: 46.3, CustomMetrics: map[string]float64{"pizzas": 9.00},
					},
				},
			},
		},
	}
	targetGroupedSets = []Set{
		{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Groups: map[string][]Benchmark{
				"^Benchmark\\S+10/.*$": {
					{
						Name: "BenchmarkFib10/Fib()-12", Runs: 3293298, NsPerOp: 330,
					},
					{
						Name: "BenchmarkPizzas10/Pizzas()-12", Runs: 25820055, NsPerOp: 50, CustomMetrics: map[string]float64{"pizzas": 3.00},
					},
				},
				"^Benchmark\\S+100/.*$": {
					{
						Name: "BenchmarkFib100/Fib()-12", Runs: 329329, NsPerOp: 329,
					},
					{
						Name: "BenchmarkPizzas100/Pizzas()-12", Runs: 2582005, NsPerOp: 50, CustomMetrics: map[string]float64{"pizzas": 3.00},
					},
				},
			},
		},
		{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Groups: map[string][]Benchmark{
				"^Benchmark\\S+10/.*$": {
					{
						Name: "BenchmarkFib10/FibSlow()-12", Runs: 3033732, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1},
					},
					{
						Name: "BenchmarkPizzas10/Pizzas()-12", Runs: 22866814, NsPerOp: 46.3, CustomMetrics: map[string]float64{"pizzas": 9.00},
					},
				},
				"^Benchmark\\S+100/.*$": {
					{
						Name: "BenchmarkFib100/FibSlow()-12", Runs: 303373, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1},
					},
					{
						Name: "BenchmarkPizzas100/Pizzas()-12", Runs: 2286681, NsPerOp: 46.3, CustomMetrics: map[string]float64{"pizzas": 9.00},
					},
				},
			},
		},
	}
	groupRegexps = []*regexp2.Regexp{
		regexp2.MustCompile("^Benchmark\\S+10/.*$", 0),
		regexp2.MustCompile("^Benchmark\\S+100/.*$", 0),
	}
)

func TestParseSet(t *testing.T) {
	convey.Convey("Given Golang standard Benchmark outputs", t, func() {
		// Golang standard Benchmark outputs are up there
		convey.Convey("Parse these outputs", func() {
			for idx, output := range benchmarkOutputs {
				set, err := ParseSet(bufio.NewReader(strings.NewReader(output)))
				convey.So(err, convey.ShouldBeNil)
				convey.So(*set, convey.ShouldResemble, targetSets[idx])
			}
		})
		convey.Convey("Parse these outputs with group", func() {
			for idx, output := range benchmarkOutputs {
				set, err := ParseSet(bufio.NewReader(strings.NewReader(output)), groupRegexps...)
				convey.So(err, convey.ShouldBeNil)
				convey.So(*set, convey.ShouldResemble, targetGroupedSets[idx])
			}
		})
	})
}

func TestParse(t *testing.T) {
	convey.Convey("Given combined Golang standard Benchmark output", t, func() {
		combinedOutput := strings.Join(benchmarkOutputs, "\n\n\n")
		convey.Convey("Parse the combined output", func() {
			sets, err := Parse(bufio.NewReader(strings.NewReader(combinedOutput)))
			convey.So(err, convey.ShouldBeNil)
			for idx, set := range sets {
				convey.So(set, convey.ShouldResemble, targetSets[idx])
			}
		})
		convey.Convey("Parse the combined output with group", func() {
			sets, err := Parse(bufio.NewReader(strings.NewReader(combinedOutput)), groupRegexps...)
			convey.So(err, convey.ShouldBeNil)
			for idx, set := range sets {
				convey.So(set, convey.ShouldResemble, targetGroupedSets[idx])
			}
		})
	})
}
