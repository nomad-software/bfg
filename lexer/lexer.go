package lexer

import (
	"github.com/nomad-software/bfg/token"
)

// Lexer is the instance of the lexer.
type Lexer struct {
	input  []byte        // The string being scanned.
	start  int           // Start position of this item.
	cur    int           // Current position in the input.
	mrk    int           // Position of saved mark.
	loops  []int         // Indexes of open lexeme types
	Tokens []token.Token // Lexed tokens
}

type stateFn func(*Lexer) stateFn

// New creates a new instance of the lexer.
func New(input []byte) *Lexer {
	l := &Lexer{
		input:  sanitise(input),
		Tokens: make([]token.Token, 0, 4096),
		loops:  make([]int, 0, 16),
	}
	l.run()
	return l
}

// Sanitise removes all invalid operators from the input.
func sanitise(input []byte) []byte {
	result := make([]byte, 0, len(input))

	for _, op := range input {
		if op == token.Add || op == token.Sub || op == token.Right || op == token.Left || op == token.Open || op == token.Close || op == token.Out || op == token.In {
			result = append(result, op)
		}
	}

	return result
}

// Run runs the lexer.
func (l *Lexer) run() {
	for f := lex; f != nil; {
		f = f(l)
	}
}

// Red returns a slice of bytes from the input that have read so far.
func (l *Lexer) red() []byte {
	return l.input[l.start:l.cur]
}

// Redbyte returns a the first byte from the input read so far.
func (l *Lexer) redbyte() byte {
	return l.red()[0]
}

// Unred returns a slice of bytes from the input that have been unread.
func (l *Lexer) unred() []byte {
	return l.input[l.cur:]
}

// Unredbyte returns a the first byte from the input that's so far unread.
func (l *Lexer) unredbyte() byte {
	return l.unred()[0]
}

// Peek returns the next unread byte from the input.
func (l *Lexer) peek() byte {
	if l.cur >= len(l.input) {
		return token.EOF
	}
	return l.unredbyte()
}

// Advance reads a byte from the input and returns it.
func (l *Lexer) advance() byte {
	if l.cur >= len(l.input) {
		return token.EOF
	}
	b := l.unredbyte()
	l.cur++
	return b
}

// Discard discards all bytes read so far.
func (l *Lexer) discard() {
	l.start = l.cur
}

// Mark saves a point to return to in the input.
func (l *Lexer) mark() {
	l.mrk = l.cur
}

// Reset returns to the saved mark.
func (l *Lexer) reset() {
	l.cur = l.mrk
	l.mrk = l.cur
}

// Emit emits a token.
func (l *Lexer) emit(t token.LexemeType, move int, value byte, jump int) {
	tok := token.Token{
		Type:  t,
		Move:  move,
		Value: value,
		Jump:  jump,
	}

	if t == token.CloseType {
		l.Tokens[jump].Jump = len(l.Tokens)
	}

	l.Tokens = append(l.Tokens, tok)
}

func lex(l *Lexer) stateFn {
	for {
		b := l.advance()

		switch b {
		case token.Right, token.Left, token.Add, token.Sub:
			return lexRepeating

		case token.In:
			return lexIn

		case token.Out:
			return lexOut

		case token.Open:
			return lexOpen

		case token.Close:
			return lexClose

		case token.EOF:
			return lexEOF
		}
	}
}

func lexRepeating(l *Lexer) stateFn {
	b := l.redbyte()

	move := 1
	value := byte(1)

	for b == l.peek() {
		move++
		value++
		l.advance()
	}

	switch b {
	case token.Right:
		l.emit(token.RightType, move, 0, 0)

	case token.Left:
		l.emit(token.LeftType, move, 0, 0)

	case token.Add:
		l.emit(token.AddType, 0, value, 0)

	case token.Sub:
		l.emit(token.SubType, 0, value, 0)
	}

	l.discard()

	return lex
}

func lexIn(l *Lexer) stateFn {
	l.emit(token.InType, 0, 0, 0)
	l.discard()
	return lex
}

func lexOut(l *Lexer) stateFn {
	l.emit(token.OutType, 0, 0, 0)
	l.discard()
	return lex
}

func lexOpen(l *Lexer) stateFn {
	if s := lexZeroLoop(l); s != nil {
		return s
	}

	if s := lexMulLoop(l); s != nil {
		return s
	}

	l.loops = append(l.loops, len(l.Tokens))
	l.emit(token.OpenType, 0, 0, 0)
	l.discard()

	return lex
}

func lexZeroLoop(l *Lexer) stateFn {
	if l.peek() != token.Sub {
		return nil
	}

	l.mark()
	l.advance()

	if l.peek() == token.Close {
		l.advance()
		l.emit(token.ZeroType, 0, 0, 0)
		l.discard()
		return lex
	}

	l.reset()
	return nil
}

func lexMulLoop(l *Lexer) stateFn {
	if l.peek() != token.Sub {
		return nil
	}

	l.mark()

	for {
		b := l.advance()

		if b == token.Open || b == token.In || b == token.Out {
			l.reset()
			return nil

		} else if b == token.Close {
			l.reset()
			break
		}
	}

	l.advance()

	move := 0

	for {
		b := l.peek()

		if b == token.Left || b == token.Right {
			if b == token.Left {
				move--
			} else {
				move++
			}

			l.advance()

			if l.peek() == token.Add {
				l.advance()
				value := byte(1)
				for l.peek() == token.Add {
					l.advance()
					value++
				}
				l.emit(token.MulAddType, move, value, 0)
			}

			if l.peek() == token.Sub {
				l.advance()
				value := byte(1)
				for l.peek() == token.Sub {
					l.advance()
					value++
				}
				l.emit(token.MulSubType, move, value, 0)
			}
		}

		if b == token.Close {
			l.advance()
			l.emit(token.ZeroType, 0, 0, 0)
			l.discard()
			return lex
		}
	}
}

func lexClose(l *Lexer) stateFn {
	jump := l.loops[len(l.loops)-1]
	l.loops = l.loops[:len(l.loops)-1]

	l.emit(token.CloseType, 0, 0, jump)
	l.discard()

	return lex
}

func lexEOF(l *Lexer) stateFn {
	l.emit(token.EOFType, 0, 0, 0)
	l.discard()
	return nil
}
