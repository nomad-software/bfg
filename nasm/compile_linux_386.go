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

	// This write seems unnecessary but it's required or else subsequent reads from
	// stdin fail. I need to investigate why.
	asm.write("mov eax, 4")
	asm.write("mov ebx, 1")
	asm.write("mov ecx, edi")
	asm.write("mov edx, 1")
	asm.write("int 0x80")

	for i, t := range tokens {
		switch t.Type {

		case token.RightType:
			asm.write("add edi, %d", t.Shift)

		case token.LeftType:
			asm.write("sub edi, %d", t.Shift)

		case token.AddType:
			asm.write("add byte [edi], %d", t.Value)

		case token.SubType:
			asm.write("sub byte [edi], %d", t.Value)

		case token.InType:
			asm.write("mov eax, 3")
			asm.write("mov ebx, 0")
			asm.write("mov ecx, edi")
			asm.write("int 0x80")

		case token.OutType:
			asm.write("mov eax, 4")
			asm.write("mov ebx, 1")
			asm.write("mov ecx, edi")
			asm.write("mov edx, 1")
			asm.write("int 0x80")

		case token.OpenType:
			asm.write("cmp byte [edi], 0")
			asm.write("je close_loop_%d", t.Jump)
			asm.write("open_loop_%d:", i)

		case token.CloseType:
			asm.write("cmp byte [edi], 0")
			asm.write("jne open_loop_%d", t.Jump)
			asm.write("close_loop_%d:", i)

		case token.ZeroType:
			asm.write("mov byte [edi], 0")
		}
	}

	asm.write("mov eax, 1")
	asm.write("mov ebx, 0")
	asm.write("int 0x80")

	asm.write("section .bss")
	asm.write("stack: resb 131072")

	return asm
}
