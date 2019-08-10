// +build !linux

package cli

import (
	"flag"
	"fmt"
)

// Options contain the command line options passed to the program.
type Options struct {
	File string
	Help bool
}

// ParseOptions parses the command line options.
func ParseOptions() *Options {
	var opt Options

	flag.StringVar(&opt.File, "f", "", "The program file to run.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.Parse()

	return &opt
}

// PrintUsage prints the usage of the program.
func (opt *Options) PrintUsage() {
	var banner = ` _      __
| |__  / _| __ _
| '_ \| |_ / _' |
| |_) |  _| (_| |
|_.__/|_|  \__, |
           |___/

A fast brainfuck interpreter and compiler.
`
	fmt.Println(banner)
	flag.Usage()
}
