package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nomad-software/bfg/cli"
	"github.com/nomad-software/bfg/compiler/c"
	"github.com/nomad-software/bfg/compiler/golang"
	"github.com/nomad-software/bfg/compiler/nasm"
	"github.com/nomad-software/bfg/evaluator"
	"github.com/nomad-software/bfg/lexer"
)

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

	if opt.Go {
		golang.Compile(tokens)

	} else if opt.Nasm {
		nasm.Compile(tokens)

	} else if opt.C {
		c.Compile(tokens)

	} else {
		input := bufio.NewReader(os.Stdin)
		output := bufio.NewWriter(os.Stdout)
		defer output.Flush()
		evaluator.Evaluate(tokens, input, output)
	}
}
