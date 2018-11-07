package token

// Brainfuck operator types.
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

var (
	// All exported operators.
	All = []byte{Left, Right, Add, Sub, In, Out, Open, Close}
)

// Type represents the operator of a token.
type Type byte

// Token represents a unit of output from the lexer.
type Token struct {
	Type    byte
	Literal string
	Shift   int
	Value   byte
}
