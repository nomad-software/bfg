package eval

import (
	"bufio"

	"github.com/nomad-software/bfg/token"
)

const (
	stackSize = 1 << 16
)

// Evaluate evaluates the program and executes it.
func Evaluate(tokens []token.Token, input bufio.Reader, output bufio.Writer) {

	var stack [stackSize]byte
	ptr := 0

	for x := 0; x < len(tokens); x++ {
		switch tokens[x].Type {

		case token.RightType:
			ptr += tokens[x].Shift

		case token.LeftType:
			ptr -= tokens[x].Shift

		case token.AddType:
			stack[ptr] += tokens[x].Value

		case token.SubType:
			stack[ptr] -= tokens[x].Value

		case token.InType:
			stack[ptr], _ = input.ReadByte()

		case token.OutType:
			output.WriteByte(stack[ptr])
			if stack[ptr] == '\n' {
				output.Flush()
			}

		case token.OpenType:
			if stack[ptr] == 0 {
				x = tokens[x].Jump
			}

		case token.CloseType:
			if stack[ptr] != 0 {
				x = tokens[x].Jump
			}

		case token.ZeroType:
			stack[ptr] = 0
		}
	}
}
