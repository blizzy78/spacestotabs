package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type line struct {
	spaces int
	text   string
}

func main() {
	flags := parseFlags()
	ok := convertFiles(flags.paths, flags)
	if !ok {
		os.Exit(1)
	}
}

func log(flags flags, p string, args ...interface{}) {
	if !flags.quiet {
		fmt.Printf(p, args...)
	}
}

func convertFiles(paths []string, flags flags) bool {
	ok := true
	for _, p := range paths {
		err := convertFile(p, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error converting file: %v\n", p, err)
			ok = false
		}
	}
	return ok
}

func convertFile(path string, flags flags) error {
	lines, err := readFile(path)
	if err != nil {
		return err
	}

	size := calcTabSize(lines)
	if size <= 0 {
		log(flags, "%s: no spaces found\n", path)
		return nil
	}

	if flags.tabSize > 0 {
		size = flags.tabSize
	}

	newLines, linesRemaining := convert(lines, size)

	if !flags.dryRun {
		err = writeFile(path, newLines)
		if err != nil {
			return err
		}
	}

	sizeText := strconv.Itoa(size)
	if flags.tabSize <= 0 {
		sizeText += " (auto)"
	}

	log(flags, "%s: converted to tab size %s, %d lines with spaces remaining\n", path, sizeText, linesRemaining)

	return nil
}

func readFile(path string) ([]line, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	return readLines(file)
}

func readLines(r io.Reader) ([]line, error) {
	lines := []line{}

	s := newScanner(r)
	for s.Scan() {
		t := s.Text()
		sp := countSpaces(t)
		lines = append(lines, line{
			spaces: sp,
			text:   t,
		})
	}

	return lines, s.Err()
}

func writeFile(path string, lines []string) error {
	b := strings.Builder{}
	for _, l := range lines {
		b.WriteString(l)
	}

	return ioutil.WriteFile(path, []byte(b.String()), 0644)
}

func newScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(scanLines)
	return s
}

// scanLines works like bufio.ScanLines, but does not drop trailing end-of-line markers.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0 : i+1], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
