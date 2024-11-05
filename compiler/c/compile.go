package c

import (
	"github.com/nomad-software/bfg/token"
)

// Compile creates an executable and runs it.
func Compile(tokens []token.Token) {
	c := newSource(tokens)
	c.writeFile("/tmp/bfg.c")
	c.compile("/tmp/bfg")
	c.run()
}

func newSource(tokens []token.Token) C {
	var src C

	src.write("#include <stdio.h>")

	src.write("char stack[1024*128] = {0};")
	src.write("int ptr = 0;")

	src.write("int main() {")

	for i, t := range tokens {
		switch t.Type {

		case token.RightType:
			src.write("ptr += %d;", t.Move)

		case token.LeftType:
			src.write("ptr -= %d;", t.Move)

		case token.AddType:
			src.write("stack[ptr] += %d;", t.Value)

		case token.SubType:
			src.write("stack[ptr] -= %d;", t.Value)

		case token.InType:
			src.write("stack[ptr] = getchar();")

		case token.OutType:
			src.write("putchar(stack[ptr]);")

		case token.OpenType:
			src.write("if (stack[ptr] == 0) {")
			src.write("goto close_%d;", i)
			src.write("}")
			src.write("open_%d:;", i)

		case token.CloseType:
			src.write("if (stack[ptr] != 0) {")
			src.write("goto open_%d;", t.Jump)
			src.write("}")
			src.write("close_%d:;", t.Jump)

		case token.MulAddType:
			src.write("if (stack[ptr] != 0) {")
			src.write("stack[ptr + %d] += (stack[ptr] * %d);", t.Move, t.Value)
			src.write("}")

		case token.MulSubType:
			src.write("if (stack[ptr] != 0) {")
			src.write("stack[ptr + %d] -= (stack[ptr] * %d);", t.Move, t.Value)
			src.write("}")

		case token.ZeroType:
			src.write("stack[ptr] = 0;")
		}
	}

	src.write("}")

	return src
}
