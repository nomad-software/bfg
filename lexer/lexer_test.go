package lexer

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/nomad-software/bfg/token"
)

var tokens []token.Token

func BenchmarkLexer(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := os.ReadFile(path.Join(wd, "../programs/mandelbrot.bf"))
	if err != nil {
		log.Fatalln(err)
	}
	var t []token.Token

	b.SetBytes(int64(len(program)))
	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		t = New(program).Tokens
	}

	tokens = t
}

func TestLexingSingleOperators(t *testing.T) {
	program := []byte("//this is a comment><+-foo,.[]bar")

	tokens := []token.Token{
		{Type: token.RightType, Move: 1, Value: 1},
		{Type: token.LeftType, Move: 1, Value: 1},
		{Type: token.AddType, Move: 1, Value: 1},
		{Type: token.SubType, Move: 1, Value: 1},
		{Type: token.InType, Move: 1, Value: 1},
		{Type: token.OutType, Move: 1, Value: 1},
		{Type: token.OpenType, Move: 0, Value: 0, Jump: 7},
		{Type: token.CloseType, Move: 0, Value: 0, Jump: 6},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultipleOperators(t *testing.T) {
	program := []byte("++++++++>+++>>>>+++<+...---<<-<--,,,.")

	tokens := []token.Token{
		{Type: token.AddType, Move: 8, Value: 8},
		{Type: token.RightType, Move: 1, Value: 1},
		{Type: token.AddType, Move: 3, Value: 3},
		{Type: token.RightType, Move: 4, Value: 4},
		{Type: token.AddType, Move: 3, Value: 3},
		{Type: token.LeftType, Move: 1, Value: 1},
		{Type: token.AddType, Move: 1, Value: 1},
		{Type: token.OutType, Move: 1, Value: 1},
		{Type: token.OutType, Move: 1, Value: 1},
		{Type: token.OutType, Move: 1, Value: 1},
		{Type: token.SubType, Move: 3, Value: 3},
		{Type: token.LeftType, Move: 2, Value: 2},
		{Type: token.SubType, Move: 1, Value: 1},
		{Type: token.LeftType, Move: 1, Value: 1},
		{Type: token.SubType, Move: 2, Value: 2},
		{Type: token.InType, Move: 1, Value: 1},
		{Type: token.InType, Move: 1, Value: 1},
		{Type: token.InType, Move: 1, Value: 1},
		{Type: token.OutType, Move: 1, Value: 1},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingZeroOptimisation(t *testing.T) {
	program := []byte("++++++++++[-]++++++++++[-]+++++[->+<]")

	tokens := []token.Token{
		{Type: token.AddType, Move: 10, Value: 10},
		{Type: token.ZeroType, Move: 3, Value: 3},
		{Type: token.AddType, Move: 10, Value: 10},
		{Type: token.ZeroType, Move: 3, Value: 3},
		{Type: token.AddType, Move: 5, Value: 5},
		{Type: token.RightMoveAddType, Move: 1},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingRightMoveAddLoopOptimisation(t *testing.T) {
	program := []byte("++++++++++[->+<][->+<<][->>>+<<<][->>>>>+<<<<<]")

	tokens := []token.Token{
		{Type: token.AddType, Move: 10, Value: 10},
		{Type: token.RightMoveAddType, Move: 1},
		{Type: token.OpenType, Move: 0, Value: 0, Jump: 7},
		{Type: token.SubType, Move: 1, Value: 1},
		{Type: token.RightType, Move: 1, Value: 1},
		{Type: token.AddType, Move: 1, Value: 1},
		{Type: token.LeftType, Move: 2, Value: 2},
		{Type: token.CloseType, Move: 0, Value: 0, Jump: 2},
		{Type: token.RightMoveAddType, Move: 3},
		{Type: token.RightMoveAddType, Move: 5},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingLeftMoveAddLoopOptimisation(t *testing.T) {
	program := []byte("++++++++++[-<+>][-<<<<<+>>>>>][-<+>>][-<<<+>>>]")

	tokens := []token.Token{
		{Type: token.AddType, Move: 10, Value: 10},
		{Type: token.LeftMoveAddType, Move: 1},
		{Type: token.LeftMoveAddType, Move: 5},
		{Type: token.OpenType, Move: 0, Value: 0, Jump: 8},
		{Type: token.SubType, Move: 1, Value: 1},
		{Type: token.LeftType, Move: 1, Value: 1},
		{Type: token.AddType, Move: 1, Value: 1},
		{Type: token.RightType, Move: 2, Value: 2},
		{Type: token.CloseType, Move: 0, Value: 0, Jump: 3},
		{Type: token.LeftMoveAddType, Move: 3},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingRightLinearAddLoopOptimisation(t *testing.T) {
	program := []byte("++++++++++[->+<][->+>+>+<<<][->+>+>+>+>+<<<<<]")

	tokens := []token.Token{
		{Type: token.AddType, Move: 10, Value: 10},
		{Type: token.RightMoveAddType, Move: 1},
		{Type: token.RightLinearAddType, Move: 3},
		{Type: token.RightLinearAddType, Move: 5},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingLeftLinearAddLoopOptimisation(t *testing.T) {
	program := []byte("++++++++++[-<+>][-<+<+<+>>>][-<+<+<+<+<+>>>>>]")

	tokens := []token.Token{
		{Type: token.AddType, Move: 10, Value: 10},
		{Type: token.LeftMoveAddType, Move: 1},
		{Type: token.LeftLinearAddType, Move: 3},
		{Type: token.LeftLinearAddType, Move: 5},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func assertTokens(t *testing.T, program []byte, tokens []token.Token) {
	output := New(program).Tokens

	for x := 0; x < len(tokens); x++ {
		typeMismatch := output[x].Type != tokens[x].Type
		moveMismatch := output[x].Move != tokens[x].Move
		valueMismatch := output[x].Value != tokens[x].Value
		jumpMismatch := output[x].Jump != tokens[x].Jump

		if typeMismatch || moveMismatch || valueMismatch || jumpMismatch {
			t.Errorf("Expected: %#v", tokens[x])
			t.Fatalf("Actual  : %#v", output[x])
		}
	}
}
