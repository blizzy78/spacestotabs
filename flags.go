package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type flags struct {
	tabSize int
	quiet   bool
	dryRun  bool
	paths   []string
}

func parseFlags() flags {
	flags := flags{}

	flag.IntVarP(&flags.tabSize, "size", "s", 0, "Specify an exact tab size to use. 0 switches to auto mode (default.)")
	flag.BoolVarP(&flags.quiet, "quiet", "q", false, "Suppress output.")
	flag.BoolVar(&flags.dryRun, "dry-run", false, "Log only what would happen without actually modifying files.")

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
