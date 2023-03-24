// Package cmd provide command for benchvisual
/*
Copyright © 2023 Kevinello kevinello42@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Kevinello/benchvisual/internal/bench"
	"github.com/Kevinello/benchvisual/internal/visual"
	"github.com/charmbracelet/log"
	"github.com/dlclark/regexp2"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	sep       = new(string)
	regexStr  = new(string)
	filePath  = new(string)
	outputDir = new(string)
	jsonMode  = new(bool)
	silent    = new(bool)
	verbose   = new(bool)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "benchvisual [--version] [--help] [-s <separator> | -r <regexp>] [-f <benchmark path>] [-o <output path>] [--json] [--verbose] [--silent]",
	Example: `  go test -bench . | benchvisual -r '^Bench(mark)?(?<target>\\S+)/(?<scenario>\\S+)$'
  benchvisual -s '/' -f "path/to/origin/benchmark/file"`,
	Short: "Parse and visualize Golang standard Benchmark output",
	Long: `Parse and visualize Golang standard Benchmark output.
benchvisual provides pipe mode and file mode, it will work in pipe mode in default, add flag -f to let it work in file mode.
there are two definitions in benchvisual -- "targets" and "scenarios", which the visualization works depends on.
"targets" means targets to generate Benchmark on, and the "scenarios" means specific Benchmark scenarios where the Benchmark run in.
In the visualization progress, we will visualize Benchmark in a concepts mapping below:
	- bench.Set(package) -> page(html)
	- metrics(ns/op...)  -> series of metrics value
	- targets            -> series name(x axis)
	- scenarios          -> dummy values in charts(group name)
benchvisual also provides json output format for your secondary development, use --json to let it output json file.`,
	Version: "0.1.4",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var regex *regexp2.Regexp
		if *sep == "" {
			// only parse regexp when sep is empty
			regex, err = regexp2.Compile(*regexStr, 0)
			if err != nil {
				return
			}
		}
		// check if the output path is exist
		if fileInfo, err := os.Stat(*outputDir); os.IsNotExist(err) {
			return fmt.Errorf("given output directory path not exist: %s", *outputDir)
		} else if err != nil {
			return fmt.Errorf("error when stat given output directory path: %s", *outputDir)
		} else if !fileInfo.IsDir() {
			return fmt.Errorf("given path is not a directory: %s", *outputDir)
		}

		if *silent {
			log.SetLevel(log.FatalLevel)
		} else if *verbose {
			log.SetLevel(log.DebugLevel)
		}

		var reader *bufio.Reader
		if *filePath != "" {
			// file mode
			f, err := os.Open(*filePath)
			if err != nil {
				return err
			}
			reader = bufio.NewReader(f)
		} else {
			// pipe mode
			reader = bufio.NewReader(os.Stdin)
		}

		sets, err := bench.Parse(reader, *sep, regex)
		if err != nil {
			return err
		}
		log.Info("Benchmark parsed success", "set_num", len(sets))

		if *jsonMode {
			// json mode, only export parsed Benchmark in json file
			setsInBytes, err := json.MarshalIndent(sets, "", "    ")
			if err != nil {
				return err
			}
			log.Debug("marshal parsed Benchmark success")
			if *outputDir == "" {
				// only print to stdout when output dir not given
				log.Print("\n" + string(setsInBytes))
			} else {
				return ioutil.WriteFile(filepath.Join(*outputDir, "parsed_benchmark.json"), setsInBytes, os.ModePerm)
			}
		}
		savedPath, err := visual.Visualize(*outputDir, sets)
		if err != nil {
			return err
		}
		log.Info("Benchmark visualized success", "saved paths", savedPath)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(filePath, "file", "f", "", "use file mode instead of pipe mode, Read the original Benchmark output from the given file path")
	rootCmd.Flags().StringVarP(sep, "sep", "s", "", "string separator of a Benchmark string's target and scenario.\ne.g., we got a benchmark name string 'BenchmarkFibonacci/100times' with separator '/', then the target of it is 'Fibonacci' and the scenario of it is '100times'.\n")
	rootCmd.Flags().StringVarP(regexStr, "regex", "r", "^Bench(mark)?(?<target>[A-Z]+\\S*)(?<scenario>[A-Z]+\\S*)$", "regexp expression with two sub groups(target and scenario), written in '.NET-style capture groups'--(?<name>re) or (?'name're).\ne.g., '^Bench(mark)?(?<target>\\S+/\\S+)/(?<scenario>\\S+)$'")
	rootCmd.Flags().StringVarP(outputDir, "output", "o", ".", "directory path to save the output file")
	rootCmd.Flags().BoolVar(jsonMode, "json", false, "only output parsed Benchmark result in json file")
	rootCmd.Flags().BoolVar(silent, "silent", false, "disable log(only show fatal log)")
	rootCmd.Flags().BoolVar(verbose, "verbose", false, "enable debug log")

	rootCmd.MarkFlagsMutuallyExclusive("sep", "regex")
	rootCmd.MarkFlagsMutuallyExclusive("silent", "verbose")
}
