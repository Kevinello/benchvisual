// Package visual visualize benchmark sets
//
//	@update 2023-03-07 02:13:12
package visual

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Kevinello/benchvisual/internal/bench"
	"github.com/Kevinello/benchvisual/internal/collections"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	options = []charts.GlobalOpts{
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Type:   "scroll",
			Orient: "vertical",
			Top:    "10%",
			Right:  "0%",
			// Padding: 5,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Benchmark\nTarget",
			SplitLine: &opts.SplitLine{
				Show: true,
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "inside",
			Start: 0,
			End:   100,
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "5%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Save as png",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					Lang:  []string{"data view", "turn off", "refresh"},
				},
				Restore: &opts.ToolBoxFeatureRestore{
					Show:  true,
					Title: "Restore view",
				},
			}},
		),
	}
)

// Visualize visualize benchmark sets and save html to target path
// every set will be visualize as 3+ bar charts for 3+ metrics(include custom metrics),
// and be exported to html files in the given saveDir
//
//	@Concept alignment:
//	bench.Set(package) -> page(html)
//	metrics(ns/op...)  -> series of metrics value
//	targets            -> series name(x axis)
//	scenarios          -> dummy values in charts(group name)
//
//	@param saveDir string
//	@param sets []bench.Set
//	@return savedPaths []string
//	@return err error
//	@author kevineluo
//	@update 2023-03-07 02:14:40
func Visualize(saveDir string, sets []bench.Set) (savedPaths []string, err error) {
	for _, set := range sets {
		// get all scenarios
		scenarios := set.GetScenario()
		// ns/op
		timePerOPChart := charts.NewBar()
		setupBarChart(timePerOPChart, &set, "Time cost per option(ns)", scenarios)
		// bytes/op
		memPerOPChart := charts.NewBar()
		setupBarChart(memPerOPChart, &set, "Alloced bytes per option", scenarios)
		// allocs/op
		allocsPerOPChart := charts.NewBar()
		setupBarChart(allocsPerOPChart, &set, "Alloc times per option", scenarios)
		// MB/s
		memPerSecChart := charts.NewBar()
		setupBarChart(memPerSecChart, &set, "Alloc mem size per sec(MB)", scenarios)
		// TODO: custom metrics
		// customMetricsCharts := make(map[string]*charts.Bar, 0)

		// generate series
		for target, benchmarks := range set.Targets {
			// sort the benchmarks by its scenario first for generating series
			sort.Sort(benchmarks)

			timePerOPSeries := collections.Map(benchmarks, func(benchmark bench.Benchmark) (barData opts.BarData) {
				return opts.BarData{Name: benchmark.Name, Value: benchmark.NsPerOp}
			})
			timePerOPChart.AddSeries(target, timePerOPSeries)
			memPerOPSeries := collections.Map(benchmarks, func(benchmark bench.Benchmark) (barData opts.BarData) {
				return opts.BarData{Name: benchmark.Name, Value: benchmark.Mem.BytesPerOp}
			})
			memPerOPChart.AddSeries(target, memPerOPSeries)
			allocsPerOPSeries := collections.Map(benchmarks, func(benchmark bench.Benchmark) (barData opts.BarData) {
				return opts.BarData{Name: benchmark.Name, Value: benchmark.Mem.AllocsPerOp}
			})
			allocsPerOPChart.AddSeries(target, allocsPerOPSeries)
			memPerSecSeries := collections.Map(benchmarks, func(benchmark bench.Benchmark) (barData opts.BarData) {
				return opts.BarData{Name: benchmark.Name, Value: benchmark.Mem.MBPerSec}
			})
			memPerSecChart.AddSeries(target, memPerSecSeries)
		}

		page := components.NewPage()
		page.AddCharts(timePerOPChart, memPerOPChart, allocsPerOPChart, memPerSecChart)
		f, err := os.Create(filepath.Join(saveDir, strings.ReplaceAll(set.Pkg, "/", "-")+".html"))
		if err != nil {
			return nil, fmt.Errorf("[Visualize] error when create result file: %w", err)
		}
		defer f.Close()
		page.Render(f)
		savedPaths = append(savedPaths, f.Name())
	}
	return
}

func setupBarChart(bar *charts.Bar, set *bench.Set, title string, scenarios []string) {
	bar.SetGlobalOptions(
		append(options,
			charts.WithTitleOpts(opts.Title{
				Title:    title,
				Subtitle: fmt.Sprintf("Package: %s\nOS: %s, ARCH: %s, CPU: %s", set.Pkg, set.Goos, set.Goarch, set.CPU),
				Top:      "0%",
				Left:     "10%",
			}),
		)...,
	).SetSeriesOptions(
		// 0 gap between bars in same scenario
		charts.WithBarChartOpts(opts.BarChart{
			BarGap: "0%",
		}),
	)
	bar.SetXAxis(scenarios)
}
