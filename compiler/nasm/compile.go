package nasm

import (
	"github.com/nomad-software/bfg/token"
)

// Compile creates an executable and runs it.
func Compile(tokens []token.Token) {
	asm := newSource(tokens)
	asm.writeFile("/tmp/bfg.asm")
	asm.compile("elf64")
	asm.link("elf_x86_64", "/tmp/bfg")
	asm.run()
}

func newSource(tokens []token.Token) Nasm {
	var asm Nasm

	asm.write("global _start")
	asm.write("section .text")
	asm.write("_start:")
	asm.write("mov r8, stack")

	for i, t := range tokens {
		switch t.Type {

		case token.RightType:
			asm.write("add r8, %d", t.Move)

		case token.LeftType:
			asm.write("sub r8, %d", t.Move)

		case token.AddType:
			asm.write("add byte [r8], %d", t.Value)

		case token.SubType:
			asm.write("sub byte [r8], %d", t.Value)

		case token.InType:
			asm.write("mov rax, 0")
			asm.write("mov rdi, 0")
			asm.write("mov rsi, r8")
			asm.write("mov rdx, 1")
			asm.write("syscall")

		case token.OutType:
			asm.write("mov rax, 1")
			asm.write("mov rdi, 1")
			asm.write("mov rsi, r8")
			asm.write("mov rdx, 1")
			asm.write("syscall")

		case token.OpenType:
			asm.write("cmp byte [r8], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)

		case token.CloseType:
			asm.write("cmp byte [r8], 0")
			asm.write("jne open_loop_%d", t.Jump)
			asm.write("close_loop_%d:", t.Jump)

		case token.MulAddType:
			asm.write("cmp byte [r8], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("mov byte al, [r8]")
			asm.write("mov byte bl, %d", t.Value)
			asm.write("mul byte bl")
			asm.write("add byte [r8+%d], al", t.Move)
			asm.write("close_loop_%d:", i)

		case token.MulSubType:
			asm.write("cmp byte [r8], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("mov byte al, [r8]")
			asm.write("mov byte bl, %d", t.Value)
			asm.write("mul byte bl")
			asm.write("sub byte [r8+%d], al", t.Move)
			asm.write("close_loop_%d:", i)

		case token.ScanRightType:
			asm.write("cmp byte [r8], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("add r8, %d", t.Move)
			asm.write("cmp byte [r8], 0")
			asm.write("jne open_loop_%d", i)
			asm.write("close_loop_%d:", i)

		case token.ScanLeftType:
			asm.write("cmp byte [r8], 0")
			asm.write("je close_loop_%d", i)
			asm.write("open_loop_%d:", i)
			asm.write("sub r8, %d", t.Move)
			asm.write("cmp byte [r8], 0")
			asm.write("jne open_loop_%d", i)
			asm.write("close_loop_%d:", i)

		case token.ZeroType:
			asm.write("mov byte [r8], 0")
		}
	}

	asm.write("mov rax, 60")
	asm.write("mov rdi, 0")
	asm.write("syscall")

	asm.write("section .bss")
	asm.write("stack: resb %d", 1024*128) // 128K

	return asm
}
