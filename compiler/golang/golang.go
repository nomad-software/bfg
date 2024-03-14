package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Go struct {
	src  strings.Builder
	file string
}

func (a *Go) write(format string, args ...interface{}) {
	a.src.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (a *Go) writeFile(name string) {
	a.file = name
	err := os.WriteFile(a.file, []byte(a.src.String()), 0666)

	if err != nil {
		fmt.Printf("cannot write go file (%s): %s\n", a.file, err)
		os.Exit(1)
	}
}

func (a *Go) run() {
	file, err := filepath.Abs(a.file)
	if err != nil {
		fmt.Printf("cannot get absolute path to go file: %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "run", file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		fmt.Printf("cannot run program: %s - %s\n", a.file, err)
		os.Exit(1)
	}
}
