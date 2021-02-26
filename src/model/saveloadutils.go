package model

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// useFile opens a file, executes the specified task and cleans up.
// 'newF' specifies whether or not to create a new file.
func useFile(path string, newF bool, task func(f *os.File) error) error {
	var _f *os.File
	var e error

	// # Optionally create file.
	if newF {
		_f, e = os.Create(path)

	} else {
		_f, e = os.Open(path)
	}

	if e != nil {
		return e
	}

	defer func() {
		if err := _f.Close(); err != nil {
			panic(err)
		}
	}()
	return task(_f)
}
func readFileLineCount(path string) (int, error) {
	count := 0
	err := useFile(path, false, func(f *os.File) error {
		buf := make([]byte, 32*1024)
		lineSep := []byte{'\n'}

		for {
			c, err := f.Read(buf)
			count += bytes.Count(buf[:c], lineSep)

			switch {
			case err == io.EOF:
				return nil

			case err != nil:
				return nil
			}
		}

	})
	return count, err
}

// readFileRows iterates over all rows in a file and executes
// a specified task func for each row. Iteration is stopped
// early if the task func returns false.
func readFileRows(path string, task func(s string) bool) error {
	return useFile(path, false, func(f *os.File) error {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if cont := task(scanner.Text()); !cont {
				return nil
			}
		}
		return nil
	})
}

// strToKeyVals tries to convert a string into a Model (type)
// record (key and val).
func strToKeyVals(s string) (string, []float64, error) {
	// # No point in continuing if the format is incorrect.
	data := strings.Split(s, " ")
	if len(data) == 0 {
		return "", []float64{}, errors.New("could not parse row")
	}

	vec := make([]float64, len(data)-1)
	// # Ignore first column since that is the key.
	for i := 1; i < len(data); i++ {
		f, err := strconv.ParseFloat(data[i], 64)
		if err != nil {
			return "", []float64{}, errors.New("could not parse float in row")
		}
		vec[i-1] = f
	}
	return data[0], vec, nil
}

// fileToModel tries to parse a a model file and put all data correctly
// into a new Model (type). Args in order: path to file; whether or not
// to print progress; whether or not to stop if any issue occurs (such
// as bad text-file/model format).
func fileToModel(path string, verbose, strict bool) (*Model, error) {
	// # Implicit linecount because that seems to increase performance,
	// # likely because it reduces the resizing of the map.
	lineCount, err := readFileLineCount(path)
	if err != nil {
		return nil, err
	}
	m := Model{Data: make(map[string][]float64, lineCount)}
	i := 0

	err = readFileRows(path, func(s string) bool {
		key, val, err := strToKeyVals(s)
		if err != nil && strict {
			return false
		}

		m.Data[key] = val
		m.Dim = len(val)

		if verbose {
			i++
			percent := 100. / float32(lineCount) * float32(i)
			fmt.Printf("\r progresss: %.4f%% (%d lines)", percent, i)
		}
		return true
	})
	// # A bit redundant but it is a way creating a new line
	// # (after the print in the prev block, which has \r).
	if verbose {
		fmt.Printf("\ndone with import of '%s'\n", path)
	}
	return &m, err
}

// keyValsToStr tries to convert a single Model (type) record
// into a string which can be saved.
func keyValsToStr(key string, vals []float64) string {
	s := key
	for i := 0; i < len(vals); i++ {
		s += fmt.Sprintf(" %f", vals[i])
	}
	return s
}

// modelToFile attempts to save a Model to a text file.
func modelToFile(m *Model, path string) error {
	return useFile(path, true, func(f *os.File) error {
		for k, v := range m.Data {
			s := keyValsToStr(k, v)
			f.WriteString(s + "\n")
		}
		return nil
	})
}
