// +build !linux

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nomad-software/bfg/cli"
	"github.com/nomad-software/bfg/eval"
	"github.com/nomad-software/bfg/lexer"
)

// Evaluate the program.
func main() {
	opt := cli.ParseOptions()

	if opt.Help {
		opt.PrintUsage()
	}

	program, err := ioutil.ReadFile(opt.File)
	if err != nil {
		fmt.Printf("Can't read program file. %s\n", err.Error())
		os.Exit(1)
	}

	tokens := lexer.New(program).Tokens

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(os.Stdout)
	defer output.Flush()

	eval.Evaluate(tokens, input, output)
}
