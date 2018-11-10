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
)

// Token represents a unit of output from the lexer.
type Token struct {
	Type    LexemeType
	Literal string
	Shift   int
	Value   byte
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
