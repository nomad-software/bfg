package evaluator

import (
	"bufio"

	"github.com/nomad-software/bfg/token"
)

const (
	stackSize = 1024 * 128
)

// Evaluate executes the program as it reads it.
func Evaluate(tokens []token.Token, input *bufio.Reader, output *bufio.Writer) {

	stack := [stackSize]byte{}
	ptr := 0

	for x := 0; x < len(tokens); x++ {
		switch tokens[x].Type {

		case token.RightType:
			ptr += tokens[x].Move

		case token.LeftType:
			ptr -= tokens[x].Move

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

		case token.MulAddType:
			if stack[ptr] != 0 {
				stack[ptr+tokens[x].Move] += (stack[ptr] * tokens[x].Value)
			}

		case token.MulSubType:
			if stack[ptr] != 0 {
				stack[ptr+tokens[x].Move] -= (stack[ptr] * tokens[x].Value)
			}

		case token.ScanRightType:
			for stack[ptr] != 0 {
				ptr += tokens[x].Move
			}

		case token.ScanLeftType:
			for stack[ptr] != 0 {
				ptr -= tokens[x].Move
			}

		case token.ZeroType:
			stack[ptr] = 0

		case token.EOFType:
			output.Flush()
		}
	}
}
