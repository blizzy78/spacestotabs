package main

import (
	"sort"
	"strings"
)

func countSpaces(s string) int {
	var c int
	for _, r := range s {
		if r != ' ' {
			break
		}
		c++
	}
	return c
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
