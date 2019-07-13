package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nomad-software/bfg/eval"
	"github.com/nomad-software/bfg/lexer"
)

// Evaluate the program.
func main() {

	if len(os.Args) <= 1 {
		fmt.Println("No program file argument")
		os.Exit(1)
	}

	program, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Can't read program file. %s\n", err.Error())
		os.Exit(1)
	}

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(os.Stdout)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	eval.Evaluate(tokens, *input, *output)
}
