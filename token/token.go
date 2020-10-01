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
	All = []byte{Left, Right, Add, Sub, In, Out, Open, Close}
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
	RightShiftAddType
	LeftShiftAddType
	RightLinearAddType
	LeftLinearAddType
)

// Token represents a unit of output from the lexer.
// The size of this struct impacts massively on the interpreter's performance.
// The current size (which is optimal) is as follow for the various
// architectures:
// 64bit: 40 bytes
// 32bit: 20 bytes
type Token struct {
	Type  LexemeType // The token lexeme type
	Shift int        // An amount to shift the stack pointer
	Value byte       // A delta value to modify a stack cell's value (packed to word boundry by compiler)
	Jump  int        // A matching position of a lexeme
	_     struct{}   // Prevent unkeyed literals and let the compiler pack it to a word boundry
}
