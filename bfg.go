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
		return
	}

	program, _ := ioutil.ReadFile(os.Args[1])

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(os.Stdout)
	defer output.Flush()

	stack := make([]byte, 30720)
	cell := 0
	loops := make([]int, 2056)
	loop := -1
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
			output.WriteByte(stack[cell])
			if stack[cell] == '\n' {
				output.Flush()
			}

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
				loop++
				loops[loop] = x
			}

		case ']':
			if stack[cell] == 0 {
				loop--
			} else {
				x = loops[loop]
			}
		}
	}
}
