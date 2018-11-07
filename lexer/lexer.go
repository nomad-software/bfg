package lexer

import (
	"bytes"
	"unicode/utf8"

	"github.com/nomad-software/bfg/token"
)

const (
	// EOF represents the end of the input/file.
	EOF = rune(token.EOF)
)

// New creates a new instance of the lexer channel.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		Tokens: make([]token.Token, 0),
	}
	l.run()
	return l
}

// Lexer is the instance of the lexer.
type Lexer struct {
	input  string        // The string being scanned.
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

func (l *Lexer) read() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) unread() string {
	return l.input[l.pos:]
}

func (l *Lexer) emit(typ byte) {
	tok := token.Token{
		Type:    typ,
		Literal: l.read(),
		Shift:   len(l.read()),
		Value:   byte(len(l.read())),
	}
	l.Tokens = append(l.Tokens, tok)
	l.start = l.pos
}

func (l *Lexer) peek() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, _ = utf8.DecodeRuneInString(l.unread())
	return r
}

func (l *Lexer) advance() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.unread())
	l.pos += l.width
	return r
}

func (l *Lexer) retreat() {
	if l.pos > l.start {
		_, l.width = utf8.DecodeLastRuneInString(l.read())
		l.pos -= l.width
	}
}

func (l *Lexer) skipInvalid() {
	for {
		r := l.peek()
		if r == EOF || bytes.ContainsRune(token.All, r) {
			return
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

		switch byte(r) {
		case token.Right:
			fallthrough

		case token.Left:
			fallthrough

		case token.Add:
			fallthrough

		case token.Sub:
			return lexRepeating

		case token.In:
			l.emit(token.In)

		case token.Out:
			l.emit(token.Out)

		case token.Open:
			l.emit(token.Open)

		case token.Close:
			l.emit(token.Close)

		case token.EOF:
			return lexEOF
		}
	}
}

func lexRepeating(l *Lexer) stateFn {
	r, _ := utf8.DecodeRuneInString(l.read())
	for r == l.peek() {
		l.advance()
	}

	l.emit(byte(r))
	return lex
}

func lexEOF(l *Lexer) stateFn {
	l.emit(token.EOF)
	return nil
}
