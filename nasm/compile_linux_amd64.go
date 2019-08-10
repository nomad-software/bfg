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
	asm.write("mov rsp, stack")

	// This write seems unnecessary but it's required or else subsequent reads from
	// stdin fail. I need to investigate why.
	asm.write("mov rax, 1")
	asm.write("mov rdi, 1")
	asm.write("mov rsi, rsp")
	asm.write("mov rdx, 1")
	asm.write("syscall")

	for i, t := range tokens {
		switch t.Type {

		case token.RightType:
			asm.write("add rsp, %d", t.Shift)

		case token.LeftType:
			asm.write("sub rsp, %d", t.Shift)

		case token.AddType:
			asm.write("add byte [rsp], %d", t.Value)

		case token.SubType:
			asm.write("sub byte [rsp], %d", t.Value)

		case token.InType:
			asm.write("mov rax, 0")
			asm.write("mov rdi, 0")
			asm.write("mov rsi, rsp")
			asm.write("syscall")

		case token.OutType:
			asm.write("mov rax, 1")
			asm.write("mov rdi, 1")
			asm.write("mov rsi, rsp")
			asm.write("mov rdx, 1")
			asm.write("syscall")

		case token.OpenType:
			asm.write("cmp byte [rsp], 0")
			asm.write("je close_loop_%d", t.Jump)
			asm.write("open_loop_%d:", i)

		case token.CloseType:
			asm.write("cmp byte [rsp], 0")
			asm.write("jne open_loop_%d", t.Jump)
			asm.write("close_loop_%d:", i)

		case token.ZeroType:
			asm.write("mov byte [rsp], 0")
		}
	}

	asm.write("mov rax, 60")
	asm.write("mov rdi, 0")
	asm.write("syscall")

	asm.write("section .bss")
	asm.write("stack: resb 32768")

	return asm.String()
}

func compileNasm(file string) string {
	obj := "/tmp/bfg.o"
	cmd := exec.Command("nasm", "-f", "elf64", "-o", obj, file)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Can't run the Netwide Assembler (nasm). %s\n", err.Error())
		os.Exit(1)
	}

	return obj
}

func link(object string, exe string) {
	cmd := exec.Command("ld", "-m", "elf_x86_64", "/tmp/bfg.o", "-o", exe)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Can't run the GNU linker (ld). %s\n", err.Error())
		os.Exit(1)
	}
}
