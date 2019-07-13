package eval

import (
	"bufio"

	"github.com/nomad-software/bfg/token"
)

// Evaluate evaluates the program and executes it.
func Evaluate(tokens []token.Token, input bufio.Reader, output bufio.Writer) {

	stack := make([]byte, 30720)
	cell := 0
	loops := make([]int, 2056)
	start := -1
	skip := 0

	for x := 0; x < len(tokens); x++ {

		switch tokens[x].Type {

		case token.RightType:
			cell += tokens[x].Shift

		case token.LeftType:
			cell -= tokens[x].Shift

		case token.AddType:
			stack[cell] += tokens[x].Value

		case token.SubType:
			stack[cell] -= tokens[x].Value

		case token.InType:
			stack[cell], _ = input.ReadByte()

		case token.OutType:
			output.WriteByte(stack[cell])
			if stack[cell] == '\n' {
				output.Flush()
			}

		case token.OpenType:
			if stack[cell] == 0 {
				skip++
				for skip > 0 {
					x++
					if tokens[x].Type == token.OpenType {
						skip++
					} else if tokens[x].Type == token.CloseType {
						skip--
					}
				}
			} else {
				start++
				loops[start] = x
			}

		case token.CloseType:
			if stack[cell] == 0 {
				start--
			} else {
				x = loops[start]
			}

		case token.ZeroType:
			stack[cell] = 0
		}
	}
}
