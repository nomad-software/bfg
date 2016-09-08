package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("No program file argument")
		os.Exit(1)
	}

	program, _ := ioutil.ReadFile(os.Args[1])

	input := bufio.NewReader(os.Stdin)
	stack := make([]byte, 30720)
	loops := make([]int, 0, 1024)
	cell := 0
	skip := 0

	for x := 0; x < len(program); x++ {

		switch program[x] {

		case '>':
			cell++

		case '<':
			cell--

		case '+':
			stack[cell]++

		case '-':
			stack[cell]--

		case '.':
			fmt.Printf("%c", stack[cell])

		case ',':
			stack[cell], _ = input.ReadByte()

		case '[':
			if stack[cell] == 0 {
				skip++
				for skip > 0 {
					x++
					if program[x] == '[' {
						skip++
					} else if program[x] == ']' {
						skip--
					}
				}
			} else {
				loops = append(loops, x)
			}

		case ']':
			if stack[cell] == 0 {
				loops = loops[:len(loops)-1]
			} else {
				x = loops[len(loops)-1]
			}
		}
	}
}
