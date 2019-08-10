package nasm

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type assembly struct {
	src strings.Builder
}

func (a *assembly) write(format string, args ...interface{}) {
	a.src.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (a *assembly) String() string {
	return a.src.String()
}

func writeFile(program string) string {
	file := "/tmp/bfg.asm"

	err := ioutil.WriteFile(file, []byte(program), 0666)
	if err != nil {
		fmt.Printf("Can't write assembly file (%s). %s\n", file, err.Error())
		os.Exit(1)
	}

	return file
}

func run(exe string) {
	exe, err := filepath.Abs(exe)
	if err != nil {
		fmt.Printf("Can't run %s. %s\n", exe, err.Error())
		os.Exit(1)
	}

	cmd := exec.Command(exe)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Can't run %s. %s\n", exe, err.Error())
		os.Exit(1)
	}
}
