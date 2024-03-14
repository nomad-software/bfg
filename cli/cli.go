package cli

import (
	"flag"
	"fmt"
	"os"
)

// Options contain the command line options passed to the program.
type Options struct {
	File      string
	Exe       string
	Interpret bool
	Nasm      bool
	Help      bool
}

// ParseOptions parses the command line options.
func ParseOptions() *Options {
	var opt Options

	flag.StringVar(&opt.File, "f", "", "The program file to run.")
	flag.BoolVar(&opt.Interpret, "i", false, "Use the cross-platform interpreter")
	flag.BoolVar(&opt.Nasm, "n", false, "Use the nasm compiler")
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

Modes:

1. Dynamic interpreter      slowest
2. Go compiler (default)    medium
3. Nasm compiler            fastest

Nasm can be installed on Linux with the following command:
$ sudo apt install nasm
`
	fmt.Println(banner)
	flag.Usage()

	os.Exit(0)
}
