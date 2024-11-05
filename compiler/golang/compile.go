package golang

import (
	"github.com/nomad-software/bfg/token"
)

// Compile creates an executable and runs it.
func Compile(tokens []token.Token) {
	asm := newSource(tokens)
	asm.writeFile("/tmp/bfg.go")
	asm.run()
}

func newSource(tokens []token.Token) Go {
	var src Go

	src.write("package main")

	src.write("import (")
	src.write(`"bufio"`)
	src.write(`"os"`)
	src.write(")")

	src.write("var (")
	src.write("stack = [1024*128]byte{}")
	src.write("ptr = 0")
	src.write("input = bufio.NewReader(os.Stdin)")
	src.write("output = bufio.NewWriter(os.Stdout)")
	src.write(")")

	src.write("func main() {")

	for i, t := range tokens {
		switch t.Type {

		case token.RightType:
			src.write("ptr += %d", t.Move)

		case token.LeftType:
			src.write("ptr -= %d", t.Move)

		case token.AddType:
			src.write("stack[ptr] += %d", t.Value)

		case token.SubType:
			src.write("stack[ptr] -= %d", t.Value)

		case token.InType:
			src.write("stack[ptr], _ = input.ReadByte()")

		case token.OutType:
			src.write("output.WriteByte(stack[ptr])")
			src.write("if stack[ptr] == '\\n' {")
			src.write("output.Flush()")
			src.write("}")

		case token.OpenType:
			src.write("if stack[ptr] == 0 {")
			src.write("goto close_%d", i)
			src.write("}")
			src.write("open_%d:", i)

		case token.CloseType:
			src.write("if stack[ptr] != 0 {")
			src.write("goto open_%d", t.Jump)
			src.write("}")
			src.write("close_%d:", t.Jump)

		case token.MulAddType:
			src.write("if stack[ptr] != 0 {")
			src.write("stack[ptr + %d] += (stack[ptr] * %d)", t.Move, t.Value)
			src.write("}")

		case token.MulSubType:
			src.write("if stack[ptr] != 0 {")
			src.write("stack[ptr + %d] -= (stack[ptr] * %d)", t.Move, t.Value)
			src.write("}")

		case token.ZeroType:
			src.write("stack[ptr] = 0")
		}
	}

	src.write("output.Flush()")
	src.write("}")

	return src
}
