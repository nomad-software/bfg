package cli

import (
	"flag"
	"fmt"
)

const (
	defaultExeName = "/tmp/bfg"
)

// Options contain the command line options passed to the program.
type Options struct {
	File      string
	Exe       string
	Interpret bool
	Help      bool
}

// ParseOptions parses the command line options.
func ParseOptions() *Options {
	var opt Options

	flag.StringVar(&opt.File, "f", "", "The program file to run.")
	flag.StringVar(&opt.Exe, "o", defaultExeName, "The name of the compiled executable.")
	flag.BoolVar(&opt.Interpret, "i", false, "Use to run the cross-platform interpreter instead of the compiler")
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
Nasm is required for compilation. This can be installed on Linux like this:

	sudo apt install nasm

Compilation is the default. To use the interpreter pass the relevant switch.
`
	fmt.Println(banner)
	flag.Usage()
}
