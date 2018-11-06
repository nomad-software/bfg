package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// Brainfuck operators.
const (
	right = '>'
	left  = '<'
	add   = '+'
	sub   = '-'
	in    = ','
	out   = '.'
	open  = '['
	close = ']'
)

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

		case right:
			fallthrough
		case left:
			fallthrough
		case add:
			fallthrough
		case sub:
			if op == current.token {
				current.count++
				break
			}

			if current.token != 0 {
				ops = append(ops, current)
			}
			current = operator{op, 1}
			break

		case in:
			fallthrough
		case out:
			fallthrough
		case open:
			fallthrough
		case close:
			ops = append(ops, current)
			current = operator{op, 1}
			break
		}
	}

	ops = append(ops, current)
	return ops
}

// Main entry point.
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
	start := -1
	skip := 0

	ops := parseOperators(program)

	for x := 0; x < len(ops); x++ {

		switch ops[x].token {

		case right:
			cell += ops[x].count

		case left:
			cell -= ops[x].count

		case add:
			stack[cell] += byte(ops[x].count)

		case sub:
			stack[cell] -= byte(ops[x].count)

		case in:
			stack[cell], _ = input.ReadByte()

		case out:
			output.WriteByte(stack[cell])
			if stack[cell] == '\n' {
				output.Flush()
			}

		case open:
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
				start++
				loops[start] = x
			}

		case close:
			if stack[cell] == 0 {
				start--
			} else {
				x = loops[start]
			}
		}
	}
}
