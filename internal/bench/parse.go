package bench

import (
	"bufio"
	"io"

	"github.com/dlclark/regexp2"
)

// Parse parse Golang standard benchmark output
//
//	@param reader LineReader
//	@return []Set Sets of structured benchmark
//	@return error
//	@author kevineluo
//	@update 2023-03-07 12:13:25
func Parse(reader *bufio.Reader, groupRegexps ...*regexp2.Regexp) ([]Set, error) {
	sets := make([]Set, 0)
	for {
		beginBytes, err := reader.Peek(16)

		if string(beginBytes) == "goos" {
			set, err := ParseSet(reader, groupRegexps...)
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
