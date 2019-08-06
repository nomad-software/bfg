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
)

// Token represents a unit of output from the lexer.
type Token struct {
	Type    LexemeType // The token lexeme type
	Literal string     // The lexed source code
	Shift   int        // An amount to shift a stack pointer
	Value   byte       // A delta value to modify a stack cell's value
	Jump    int        // A matching position of a lexeme
}

var types = map[byte]LexemeType{
	Left:  LeftType,
	Right: RightType,
	Add:   AddType,
	Sub:   SubType,
	In:    InType,
	Out:   OpenType,
	Open:  OpenType,
	Close: CloseType,
	EOF:   EOFType,
}

// LookupType returns the token type for the passed operator.
func LookupType(op byte) LexemeType {
	if tok, ok := types[op]; ok {
		return tok
	}
	panic("Lexeme type not registerd.")
}
