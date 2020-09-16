package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/scanner"

	flag "github.com/spf13/pflag"
)

type flags struct {
	tabSize int
	quiet   bool
	paths   []string
}

type line struct {
	spaces int
	text   string
}

func main() {
	flags := parseFlags()

	errors := false

	for _, p := range flags.paths {
		err := convertFile(p, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error converting file: %v\n", p, err)
			errors = true
		}
	}

	if errors {
		os.Exit(1)
	}
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

	err = writeFile(path, newLines)
	if err != nil {
		return err
	}

	sizeText := strconv.Itoa(size)
	if flags.tabSize <= 0 {
		sizeText += " (auto)"
	}

	log(flags, "%s: converted to tab size %s, %d lines with spaces remaining\n", path, sizeText, linesRemaining)

	return nil
}

func parseFlags() flags {
	flags := flags{}

	flag.IntVarP(&flags.tabSize, "size", "s", 0, "Specify an exact tab size to use. 0 switches to auto mode (default.)")
	flag.BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output.")

	flag.CommandLine.Init("", flag.ContinueOnError)
	err := flag.CommandLine.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	flags.paths = flag.Args()
	if len(flags.paths) <= 0 {
		fmt.Fprintln(os.Stderr, "No files given, exiting.")
		os.Exit(1)
	}

	return flags
}

func log(flags flags, p string, args ...interface{}) {
	if !flags.quiet {
		fmt.Printf(p, args...)
	}
}

func writeFile(path string, lines []string) error {
	b := strings.Builder{}
	for _, l := range lines {
		b.WriteString(l)
	}

	return ioutil.WriteFile(path, []byte(b.String()), 0644)
}

func convert(lines []line, tabSize int) ([]string, int) {
	newLines := []string{}

	linesRemaining := 0

	for _, l := range lines {
		n := l.spaces / tabSize
		newLines = append(newLines, strings.Repeat("\t", n)+string([]rune(l.text)[n*tabSize:]))

		if l.spaces > n*tabSize {
			linesRemaining++
		}
	}

	return newLines, linesRemaining
}

func calcTabSize(lines []line) int {
	spacesFound := false
	for _, l := range lines {
		if l.spaces > 0 {
			spacesFound = true
			break
		}
	}
	if !spacesFound {
		return -1
	}

	type tabSize struct {
		size  int
		lines int
	}

	sizes := []tabSize{}
	for s := 10; s >= 2; s-- {
		c := 0
		for _, l := range lines {
			if l.spaces%s == 0 {
				c++
			}
		}

		sizes = append(sizes, tabSize{
			size:  s,
			lines: c,
		})
	}

	sort.SliceStable(sizes, func(a, b int) bool {
		s1 := sizes[a]
		s2 := sizes[b]
		return s1.lines > s2.lines
	})

	return sizes[0].size
}

func readFile(path string) ([]line, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	lines := []line{}

	s := newScanner(file)
	for s.Scan() {
		t := s.Text()

		var sp int
		sp, err = countSpaces(t)
		if err != nil {
			return nil, err
		}

		lines = append(lines, line{
			spaces: sp,
			text:   t,
		})
	}

	return lines, s.Err()
}

func countSpaces(s string) (int, error) {
	c := 0

	var scanErr error
	sc := scanner.Scanner{
		Error: func(s *scanner.Scanner, msg string) {
			if scanErr == nil {
				scanErr = errors.New(msg)
			}
		},
	}

	sc.Init(strings.NewReader(s))

	for {
		r := sc.Next()
		if scanErr != nil {
			return 0, scanErr
		}

		if r == scanner.EOF {
			break
		}

		if r != ' ' {
			break
		}

		c++
	}

	return c, nil
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
