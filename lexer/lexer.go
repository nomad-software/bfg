package lexer

import (
	"github.com/nomad-software/bfg/token"
)

// New creates a new instance of the lexer channel.
func New(input []byte) *Lexer {
	l := &Lexer{
		input:  input,
		Tokens: make([]token.Token, 0),
	}
	l.run()
	return l
}

// Lexer is the instance of the lexer.
type Lexer struct {
	input  []byte        // The string being scanned.
	start  int           // Start position of this item.
	pos    int           // Current position in the input.
	width  int           // Width of the last rune read.
	Tokens []token.Token // Lexed tokens
}

type stateFn func(*Lexer) stateFn

func (l *Lexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
}

func (l *Lexer) read() []byte {
	return l.input[l.start:l.pos]
}

func (l *Lexer) unread() []byte {
	return l.input[l.pos:]
}

func (l *Lexer) emit(t token.LexemeType) {
	tok := token.Token{
		Type:    t,
		Literal: string(l.read()),
		Shift:   len(l.read()),
		Value:   byte(len(l.read())),
	}
	l.Tokens = append(l.Tokens, tok)
	l.start = l.pos
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		l.width = 0
		return token.EOF
	}
	return l.unread()[0]
}

func (l *Lexer) advance() (b byte) {
	if l.pos >= len(l.input) {
		l.width = 0
		return token.EOF
	}
	b = l.unread()[0]
	l.pos++
	return b
}

func (l *Lexer) retreat(i int) {
	if l.pos > l.start {
		l.pos -= i
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
	l.start = l.pos
}

func lex(l *Lexer) stateFn {
	for {
		l.skipInvalid()
		l.discard()

		r := l.advance()

		switch r {
		case token.Right:
			fallthrough

		case token.Left:
			fallthrough

		case token.Add:
			fallthrough

		case token.Sub:
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

	l.emit(token.LookupType(b))
	return lex
}

func lexSequence(l *Lexer, seq []byte, t token.LexemeType) stateFn {
	pos := l.pos
	for _, b := range seq {
		if b == l.peek() {
			l.advance()
		} else {
			l.retreat(l.pos - pos)
			return nil
		}
	}

	l.emit(t)
	return lex
}

func lexZero(l *Lexer) stateFn {
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
	s := lexSequence(l, []byte{token.Sub, token.Close}, token.ZeroType)
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
