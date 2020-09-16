package main

import (
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func Test_CountSpaces(t *testing.T) {
	is := is.New(t)

	is.Equal(countSpaces("foo"), 0)
	is.Equal(countSpaces("foo   "), 0)
	is.Equal(countSpaces("   "), 3)
	is.Equal(countSpaces("   foo   "), 3)
	is.Equal(countSpaces("     foo"), 5)
}

func Test_CalcTabSize(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		lines    []string
		expected int
	}{
		{
			[]string{
				"    foo",
				"bar",
				"        baz",
			},
			4,
		},
		{
			[]string{
				"   foo",
				"bar",
				"        baz",
				"    qux",
			},
			4,
		},
		{
			[]string{
				"   foo",
				"bar",
				"        baz",
				"    qux",
			},
			4,
		},
		{
			[]string{
				"    /**",
				"     * foo",
				"     */",
				"    void foo() {}",
				"    void bar() {}",
				"    void baz() {}",
			},
			4,
		},
		{
			[]string{
				"foo",
				"bar",
			},
			-1,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			lines := []line{}
			for _, s := range test.lines {
				lines = append(lines, line{
					text:   s,
					spaces: countSpaces(s),
				})
			}

			is.Equal(calcTabSize(lines), test.expected)
		})
	}
}

func Test_Convert(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		lines             []string
		tabSize           int
		expected          []string
		expectedRemaining int
	}{
		{
			[]string{
				"    foo",
				"bar",
				"        baz",
			},
			4,
			[]string{
				"\tfoo",
				"bar",
				"\t\tbaz",
			},
			0,
		},
		{
			[]string{
				"   foo",
				"bar",
				"        baz",
				"    qux",
			},
			4,
			[]string{
				"   foo",
				"bar",
				"\t\tbaz",
				"\tqux",
			},
			1,
		},
		{
			[]string{
				"    /**",
				"     * foo",
				"     */",
				"    void foo() {}",
				"    void bar() {}",
				"    void baz() {}",
			},
			4,
			[]string{
				"\t/**",
				"\t * foo",
				"\t */",
				"\tvoid foo() {}",
				"\tvoid bar() {}",
				"\tvoid baz() {}",
			},
			2,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			lines := []line{}
			for _, s := range test.lines {
				lines = append(lines, line{
					text:   s,
					spaces: countSpaces(s),
				})
			}

			convertedLines, remainingLines := convert(lines, test.tabSize)
			is.Equal(convertedLines, test.expected)
			is.Equal(remainingLines, test.expectedRemaining)
		})
	}
}
