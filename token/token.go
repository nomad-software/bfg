package token

import "fmt"

// Brainfuck operators.
const (
	Left  byte = '<'
	Right byte = '>'
	Add   byte = '+'
	Sub   byte = '-'
	In    byte = ','
	Out   byte = '.'
	Open  byte = '['
	Close byte = ']'
	EOF   byte = 255
)

// TokenType is a token type.
type TokenType int

// Token types.
const (
	LeftType TokenType = iota
	RightType
	AddType
	SubType
	InType
	OutType
	OpenType
	CloseType
	EOFType
	ZeroType
	MulAddType
	MulSubType
	ScanRightType
	ScanLeftType
)

// Token represents a unit of output from the lexer.
type Token struct {
	Type  TokenType // The token type.
	Move  int       // An amount to move the stack pointer.
	Value byte      // A delta value to modify a stack cell's value.
	Jump  int       // A matching position of a token.
}

// String implements the stringer interface.
func (t Token) String() string {
	return fmt.Sprintf("type: %d, move: %d, value: %d, jump: %d", t.Type, t.Move, t.Value, t.Jump)
}
