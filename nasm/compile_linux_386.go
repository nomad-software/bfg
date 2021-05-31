package nasm

import (
	"github.com/nomad-software/bfg/token"
)

// Compile creates an executable for a particular architecture.
func Compile(tokens []token.Token, name string) {
	asm := newAssembly(tokens)
	asm.writeFile("/tmp/bfg.asm")
	asm.compile("elf32")
	asm.link("elf_i386", name)
	asm.run()
}

func newAssembly(tokens []token.Token) nasm {
	var asm nasm

	asm.write("global _start")
	asm.write("section .text")
	asm.write("_start:")
	asm.write("mov edi, stack")

	for i, t := range tokens {
		switch t.Type {

		case token.RightType:
			asm.write("add edi, %d", t.Move)

		case token.LeftType:
			asm.write("sub edi, %d", t.Move)

		case token.AddType:
			asm.write("add byte [edi], %d", t.Value)

		case token.SubType:
			asm.write("sub byte [edi], %d", t.Value)

		case token.InType:
			asm.write("mov eax, 3")
			asm.write("mov ebx, 0")
			asm.write("mov ecx, edi")
			asm.write("mov edx, 1")
			asm.write("int 0x80")

		case token.OutType:
			asm.write("mov eax, 4")
			asm.write("mov ebx, 1")
			asm.write("mov ecx, edi")
			asm.write("mov edx, 1")
			asm.write("int 0x80")

		case token.OpenType:
			asm.write("cmp byte [edi], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)

		case token.CloseType:
			asm.write("cmp byte [edi], 0")
			asm.write("jne open_loop_%d", t.Jump)
			asm.write("close_loop_%d:", t.Jump)

		case token.ZeroType:
			asm.write("mov byte [edi], 0")

		case token.RightMoveAddType:
			asm.write("cmp byte [edi], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("mov byte al, [edi]")
			asm.write("mov byte [edi], 0")
			asm.write("add edi, %d", t.Move)
			asm.write("add byte [edi], al")
			asm.write("sub edi, %d", t.Move)
			asm.write("close_loop_%d:", i)

		case token.LeftMoveAddType:
			asm.write("cmp byte [edi], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("mov byte al, [edi]")
			asm.write("mov byte [edi], 0")
			asm.write("sub edi, %d", t.Move)
			asm.write("add byte [edi], al")
			asm.write("add edi, %d", t.Move)
			asm.write("close_loop_%d:", i)

		case token.RightLinearAddType:
			asm.write("cmp byte [edi], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("mov byte al, [edi]")
			asm.write("mov byte [edi], 0")
			for i := 1; i <= t.Move; i++ {
				asm.write("inc edi")
				asm.write("add byte [edi], al")
			}
			asm.write("sub edi, %d", t.Move)
			asm.write("close_loop_%d:", i)

		case token.LeftLinearAddType:
			asm.write("cmp byte [edi], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("mov byte al, [edi]")
			asm.write("mov byte [edi], 0")
			for i := 1; i <= t.Move; i++ {
				asm.write("dec edi")
				asm.write("add byte [edi], al")
			}
			asm.write("add edi, %d", t.Move)
			asm.write("close_loop_%d:", i)
		}
	}

	asm.write("mov eax, 1")
	asm.write("mov ebx, 0")
	asm.write("int 0x80")

	asm.write("section .bss")
	asm.write("stack: resb 131072")

	return asm
}
