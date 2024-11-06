package cli

import (
	"flag"
	"fmt"
	"os"
)

// Options contain the command line options passed to the program.
type Options struct {
	File string
	Exe  string
	Go   bool
	Nasm bool
	C    bool
	Help bool
}

// ParseOptions parses the command line options.
func ParseOptions() *Options {
	var opt Options

	flag.StringVar(&opt.File, "f", "", "The program file to run.")
	flag.BoolVar(&opt.Go, "g", false, "Use the go compiler")
	flag.BoolVar(&opt.Nasm, "n", false, "Use the nasm compiler")
	flag.BoolVar(&opt.C, "c", false, "Use the c compiler")
	flag.BoolVar(&opt.Help, "h", false, "Show help.")
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

1. Cross-platform interpreter  slowest  (default)
2. Go compiler                 medium
3. Nasm compiler               fast
4. C compiler                  fastest

Compilers can be installed on Linux with the following commands:
$ sudo apt install build-essential
$ sudo apt install nasm
`
	fmt.Println(banner)
	flag.Usage()

	os.Exit(0)
}
