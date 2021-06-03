package nasm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type nasm struct {
	src  strings.Builder
	file string
	obj  string
	exe  string
}

func (a *nasm) write(format string, args ...interface{}) {
	a.src.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (a *nasm) writeFile(name string) {
	a.file = name
	err := os.WriteFile(a.file, []byte(a.src.String()), 0666)

	if err != nil {
		fmt.Printf("Can't write assembly file (%s). %s\n", a.file, err.Error())
		os.Exit(1)
	}
}

func (a *nasm) compile(arch string) {
	a.obj = "/tmp/bfg.o"
	cmd := exec.Command("nasm", "-f", arch, "-o", a.obj, a.file)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Can't run the Netwide Assembler (nasm). %s\n", err.Error())
		os.Exit(1)
	}
}

func (a *nasm) link(arch string, exe string) {
	exe, err := filepath.Abs(exe)
	if err != nil {
		fmt.Printf("Can't run the GNU linker (ld). %s\n", err.Error())
		os.Exit(1)
	}

	a.exe = exe
	cmd := exec.Command("ld", "-m", arch, a.obj, "-o", a.exe)

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Can't run the GNU linker (ld). %s\n", err.Error())
		os.Exit(1)
	}
}

func (a *nasm) run() {
	cmd := exec.Command(a.exe)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Can't run %s. %s\n", a.exe, err.Error())
		os.Exit(1)
	}
}
