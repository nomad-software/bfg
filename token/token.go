package token

// Brainfuck operator types.
const (
	Left  byte   = '<'
	Right byte   = '>'
	Add   byte   = '+'
	Sub   byte   = '-'
	In    byte   = ','
	Out   byte   = '.'
	Open  byte   = '['
	Close byte   = ']'
	All   string = "<>+-,.[]"
)

// Miscellaneous types
const (
	EOF byte = 255
)

// Type represents the operator of a token.
type Type byte

// Token represents a unit of output from the lexer.
type Token struct {
	Type    byte
	Literal string
	Count   int
}
