package nasm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Nasm struct {
	src  strings.Builder
	file string
	obj  string
	exe  string
}

func (a *Nasm) write(format string, args ...interface{}) {
	a.src.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (a *Nasm) writeFile(name string) {
	a.file = name
	err := os.WriteFile(a.file, []byte(a.src.String()), 0666)

	if err != nil {
		fmt.Printf("cannot write assembly file (%s): %s\n", a.file, err)
		os.Exit(1)
	}
}

func (a *Nasm) compile(arch string) {
	a.obj = "/tmp/bfg.o"
	cmd := exec.Command("nasm", "-Ox", "-f", arch, "-o", a.obj, a.file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("cannot run the netwide assembler (nasm): %s\n", err)
		os.Exit(1)
	}
}

func (a *Nasm) link(arch string, exe string) {
	exe, err := filepath.Abs(exe)
	if err != nil {
		fmt.Printf("cannot run the gnu linker (ld): %s\n", err)
		os.Exit(1)
	}

	a.exe = exe
	cmd := exec.Command("ld", "-m", arch, a.obj, "-o", a.exe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("cannot run the gnu linker (ld): %s\n", err)
		os.Exit(1)
	}
}

func (a *Nasm) run() {
	cmd := exec.Command(a.exe)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Printf("cannot run program: %s - %s\n", a.exe, err)
		os.Exit(1)
	}
}
