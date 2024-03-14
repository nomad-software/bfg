package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nomad-software/bfg/cli"
	"github.com/nomad-software/bfg/compiler/golang"
	"github.com/nomad-software/bfg/compiler/nasm"
	"github.com/nomad-software/bfg/interpreter/eval"
	"github.com/nomad-software/bfg/lexer"
)

// Evaluate the program.
func main() {
	opt := cli.ParseOptions()

	if opt.Help {
		opt.PrintUsage()
		return
	}

	program, err := os.ReadFile(opt.File)
	if err != nil {
		fmt.Printf("Can't read program file. %s\n", err.Error())
		os.Exit(1)
	}

	tokens := lexer.New(program).Tokens

	if opt.Interpret {
		input := bufio.NewReader(os.Stdin)
		output := bufio.NewWriter(os.Stdout)
		defer output.Flush()
		eval.Evaluate(tokens, input, output)

	} else if opt.Nasm {
		nasm.Compile(tokens)

	} else {
		golang.Compile(tokens)
	}
}
