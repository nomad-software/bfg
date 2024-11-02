package lexer

import (
	"github.com/nomad-software/bfg/token"
)

// Lexer is the instance of the lexer.
type Lexer struct {
	input  []byte        // The string being scanned.
	start  int           // Start position of this item.
	end    int           // Current position in the input.
	loops  []int         // Indexes of open lexeme types
	Tokens []token.Token // Lexed tokens
}

type stateFn func(*Lexer) stateFn

// New creates a new instance of the lexer.
func New(input []byte) *Lexer {
	l := &Lexer{
		input:  input,
		Tokens: make([]token.Token, 0, 4096),
		loops:  make([]int, 0, 16),
	}
	l.run()
	return l
}

func (l *Lexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
}

func (l *Lexer) read() []byte {
	return l.input[l.start:l.end]
}

func (l *Lexer) unread() []byte {
	return l.input[l.end:]
}

func (l *Lexer) emit(t token.LexemeType) {
	tok := token.Token{
		Type: t,
	}

	switch t {
	case token.OpenType:
		l.loops = append(l.loops, len(l.Tokens))

	case token.CloseType:
		tok.Jump = l.loops[len(l.loops)-1]
		l.loops = l.loops[:len(l.loops)-1]
		l.Tokens[tok.Jump].Jump = len(l.Tokens)

	case token.LeftType, token.RightType, token.AddType, token.SubType:
		tok.Move = len(l.read())
		tok.Value = byte(len(l.read()))
	}

	l.Tokens = append(l.Tokens, tok)
	l.start = l.end
}

func (l *Lexer) peek() byte {
	if l.end >= len(l.input) {
		return token.EOF
	}
	return l.unread()[0]
}

func (l *Lexer) advance() (b byte) {
	if l.end >= len(l.input) {
		return token.EOF
	}
	b = l.unread()[0]
	l.end++
	return b
}

func (l *Lexer) retreat(i int) {
	if l.end > l.start {
		l.end -= i
	}
}

func (l *Lexer) skipInvalid() {
	for {
		b := l.peek()
		if b == token.EOF {
			return
		}
		for _, op := range token.All {
			if b == op {
				return
			}
		}
		l.advance()
	}
}

func (l *Lexer) discard() {
	l.start = l.end
}

func lex(l *Lexer) stateFn {
	for {
		l.skipInvalid()
		l.discard()

		b := l.advance()

		switch b {
		case token.Right, token.Left, token.Add, token.Sub:
			return lexRepeating

		case token.In:
			l.emit(token.InType)

		case token.Out:
			l.emit(token.OutType)

		case token.Open:
			return lexOpen

		case token.Close:
			l.emit(token.CloseType)

		case token.EOF:
			return lexEOF
		}
	}
}

func lexRepeating(l *Lexer) stateFn {
	b := l.read()[0]

	for b == l.peek() {
		l.advance()
	}

	switch b {
	case token.Right:
		l.emit(token.RightType)

	case token.Left:
		l.emit(token.LeftType)

	case token.Add:
		l.emit(token.AddType)

	case token.Sub:
		l.emit(token.SubType)
	}

	return lex
}

func lexZeroLoop(l *Lexer) stateFn {
	if l.peek() == token.Sub {
		l.advance()
		if l.peek() == token.Close {
			l.advance()
			l.emit(token.ZeroType)
			return lex
		}
		l.retreat(1)
	}
	return nil
}

func lexOpen(l *Lexer) stateFn {
	s := lexZeroLoop(l)
	if s != nil {
		return s
	}
	l.emit(token.OpenType)
	return lex
}

func lexEOF(l *Lexer) stateFn {
	l.emit(token.EOFType)
	return nil
}
