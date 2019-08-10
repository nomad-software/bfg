package nasm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nomad-software/bfg/token"
)

// Compile creates an executable for a particular architecture.
func Compile(tokens []token.Token, name string) {
	program := generateAssembly(tokens)
	file := writeFile(program)
	object := compileNasm(file)

	link(object, name)
	run(name)
}

func generateAssembly(tokens []token.Token) string {
	var asm assembly

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
	asm.write("stack: resb 32768")

	return asm.String()
}

func compileNasm(file string) string {
	obj := "/tmp/bfg.o"
	cmd := exec.Command("nasm", "-f", "elf32", "-o", obj, file)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Can't run the Netwide Assembler (nasm). %s\n", err.Error())
		os.Exit(1)
	}

	return obj
}

func link(object string, exe string) {
	cmd := exec.Command("ld", "-m", "elf_i386", "/tmp/bfg.o", "-o", exe)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Can't run the GNU linker (ld). %s\n", err.Error())
		os.Exit(1)
	}
}
