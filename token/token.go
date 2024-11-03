package token

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

var (
	// All exported operators.
	All = []byte{Add, Sub, Right, Left, Open, Close, Out, In}
)

// LexemeType is a lexeme type.
type LexemeType int

// Lexeme types.
const (
	LeftType LexemeType = iota
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
)

// Token represents a unit of output from the lexer.
// The size of this struct impacts massively on the interpreter's performance.
// The current size (which is optimal) is as follow for the various
// architectures:
// 64bit: 40 bytes
// 32bit: 20 bytes
type Token struct {
	Type  LexemeType // The token lexeme type
	Move  int        // An amount to move the stack pointer
	Value byte       // A delta value to modify a stack cell's value (packed to word boundry by compiler)
	Jump  int        // A matching position of a lexeme
	_     struct{}   // Prevent unkeyed literals and let the compiler pack it to a word boundry
}
