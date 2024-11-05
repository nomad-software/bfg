package c

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type C struct {
	src  strings.Builder
	file string
	exe  string
}

func (a *C) write(format string, args ...interface{}) {
	a.src.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (a *C) writeFile(name string) {
	a.file = name
	err := os.WriteFile(a.file, []byte(a.src.String()), 0666)

	if err != nil {
		fmt.Printf("cannot write c file (%s): %s\n", a.file, err)
		os.Exit(1)
	}
}

func (a *C) compile(exe string) {
	exe, err := filepath.Abs(exe)
	if err != nil {
		fmt.Printf("cannot run gnu c compiler (gcc): %s\n", err)
		os.Exit(1)
	}

	a.exe = exe
	cmd := exec.Command("gcc", "-Ofast", "-o", a.exe, a.file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("cannot run gnu c compiler (gcc): %s\n", err)
		os.Exit(1)
	}
}

func (a *C) run() {
	cmd := exec.Command(a.exe)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Printf("cannot run program: %s - %s\n", a.exe, err)
		os.Exit(1)
	}
}
