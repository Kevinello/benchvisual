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
BenchmarkFib10	3293298	330 ns/op
BenchmarkFib100	329329	329 ns/op
BenchmarkPizzas10	25820055	50.0 ns/op	3.00 pizzas
BenchmarkPizzas100	2582005	50.0 ns/op	3.00 pizzas
PASS`,
		`goos: darwin
goarch: amd64
pkg: go.bobheadxi.dev/gobenchdata/demo
BenchmarkFib/10	3033732	358 ns/op	16 B/op	1 allocs/op
BenchmarkFib/100	303373	358 ns/op	16 B/op	1 allocs/op
BenchmarkPizzas/10	22866814	46.3 ns/op	9.00 pizzas	0 B/op	0 allocs/op
BenchmarkPizzas/100	2286681	46.3 ns/op	9.00 pizzas	0 B/op	0 allocs/op
PASS`,
	}
	targetSets = []Set{
		{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			CPU:    "Intel AMD Xeon Phenom",
			Targets: map[string]BenchmarkList{
				"Fib": {
					{
						Name: "BenchmarkFib10", Runs: 3293298, NsPerOp: 330, Target: "Fib", Scenario: "10",
					},
					{
						Name: "BenchmarkFib100", Runs: 329329, NsPerOp: 329, Target: "Fib", Scenario: "100",
					},
				},
				"Pizzas": {
					{
						Name: "BenchmarkPizzas10", Runs: 25820055, NsPerOp: 50, CustomMetrics: map[string]float64{"pizzas": 3.00}, Target: "Pizzas", Scenario: "10",
					},
					{
						Name: "BenchmarkPizzas100", Runs: 2582005, NsPerOp: 50, CustomMetrics: map[string]float64{"pizzas": 3.00}, Target: "Pizzas", Scenario: "100",
					},
				},
			},
		},
		{
			Goos:   "darwin",
			Goarch: "amd64",
			Pkg:    "go.bobheadxi.dev/gobenchdata/demo",
			Targets: map[string]BenchmarkList{
				"Fib": {
					{
						Name: "BenchmarkFib/10", Runs: 3033732, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1}, Target: "Fib", Scenario: "10",
					},
					{
						Name: "BenchmarkFib/100", Runs: 303373, NsPerOp: 358, Mem: Mem{BytesPerOp: 16, AllocsPerOp: 1}, Target: "Fib", Scenario: "100",
					},
				},
				"Pizzas": {
					{
						Name: "BenchmarkPizzas/10", Runs: 22866814, NsPerOp: 46.3, CustomMetrics: map[string]float64{"pizzas": 9.00}, Target: "Pizzas", Scenario: "10",
					},
					{
						Name: "BenchmarkPizzas/100", Runs: 2286681, NsPerOp: 46.3, CustomMetrics: map[string]float64{"pizzas": 9.00}, Target: "Pizzas", Scenario: "100",
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
			set, err := ParseSet(bufio.NewReader(strings.NewReader(benchmarkOutputs[0])), "", regexp2.MustCompile("^Bench(mark)?(?<target>[A-Z][a-z]*)(?<scenario>[^A-Z]\\S+)", 0))
			convey.So(err, convey.ShouldBeNil)
			convey.So(*set, convey.ShouldResemble, targetSets[0])

			set, err = ParseSet(bufio.NewReader(strings.NewReader(benchmarkOutputs[1])), "/", nil)
			convey.So(err, convey.ShouldBeNil)
			convey.So(*set, convey.ShouldResemble, targetSets[1])
		})
	})
}

func TestParse(t *testing.T) {
	convey.Convey("Given combined Golang standard Benchmark output", t, func() {
		combinedOutput := benchmarkOutputs[1] + "\n\n\n" + benchmarkOutputs[1]
		convey.Convey("Parse the combined output", func() {
			sets, err := Parse(bufio.NewReader(strings.NewReader(combinedOutput)), "/", nil)
			convey.So(err, convey.ShouldBeNil)
			for _, set := range sets {
				convey.So(set, convey.ShouldResemble, targetSets[1])
			}
		})
	})
}
