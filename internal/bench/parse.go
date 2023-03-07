package bench

import (
	"bufio"
	"io"

	"github.com/dlclark/regexp2"
)

// Parse parse Golang standard benchmark output
//
//	@param reader *bufio.Reader
//	@param sep string sep of a Benchmark string's target and scenario
//	@return []Set Sets of structured benchmark
//	@return error
//	@author kevineluo
//	@update 2023-03-07 01:29:47
func Parse(reader *bufio.Reader, sep string, regex *regexp2.Regexp) ([]Set, error) {
	sets := make([]Set, 0)
	for {
		beginBytes, err := reader.Peek(4)

		beginStr := string(beginBytes)
		if beginStr == "goos" {
			set, err := ParseSet(reader, sep, regex)
			if err != nil {
				return nil, err
			}
			sets = append(sets, *set)
		}

		_, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return sets, nil
}
