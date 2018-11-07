package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nomad-software/bfg/lexer"
	"github.com/nomad-software/bfg/token"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("No program file argument")
		return
	}

	program, _ := ioutil.ReadFile(os.Args[1])

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(os.Stdout)
	defer output.Flush()

	stack := make([]byte, 30720)
	cell := 0
	loops := make([]int, 2056)
	start := -1
	skip := 0

	tokens := lexer.New(string(program)).Tokens

	for x := 0; x < len(tokens); x++ {

		switch tokens[x].Type {

		case token.Right:
			cell += tokens[x].Shift

		case token.Left:
			cell -= tokens[x].Shift

		case token.Add:
			stack[cell] += tokens[x].Value

		case token.Sub:
			stack[cell] -= tokens[x].Value

		case token.In:
			stack[cell], _ = input.ReadByte()

		case token.Out:
			output.WriteByte(stack[cell])
			if stack[cell] == '\n' {
				output.Flush()
			}

		case token.Open:
			if stack[cell] == 0 {
				skip++
				for skip > 0 {
					x++
					if tokens[x].Type == token.Open {
						skip++
					} else if tokens[x].Type == token.Close {
						skip--
					}
				}
			} else {
				start++
				loops[start] = x
			}

		case token.Close:
			if stack[cell] == 0 {
				start--
			} else {
				x = loops[start]
			}
		}
	}
}
