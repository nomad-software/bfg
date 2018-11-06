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

	ops := parseOperators(program)

	for x := 0; x < len(ops); x++ {

		switch ops[x].token {

		case '>':
			cell += ops[x].count

		case '<':
			cell -= ops[x].count

		case '+':
			stack[cell] += byte(ops[x].count)

		case '-':
			stack[cell] -= byte(ops[x].count)

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
					if ops[x].token == '[' {
						skip++
					} else if ops[x].token == ']' {
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

// Operator is a struct holding a parsed operator and how many times it's used
// consecutively.
type operator struct {
	token byte
	count int
}

// ParseOperators reads the operators from the program and count the number of
// times they are used consecutively.
func parseOperators(program []byte) []operator {
	ops := make([]operator, 1024)
	var current operator

	for x := 0; x < len(program); x++ {
		op := program[x]

		switch op {

		case '>':
			fallthrough
		case '<':
			fallthrough
		case '+':
			fallthrough
		case '-':
			if op == current.token {
				current.count++
				break
			}

			if current.token != 0 {
				ops = append(ops, current)
			}
			current = operator{op, 1}
			break

		case '.':
			fallthrough
		case ',':
			fallthrough
		case '[':
			fallthrough
		case ']':
			ops = append(ops, current)
			current = operator{op, 1}
			break
		}
	}

	ops = append(ops, current)
	return ops
}
