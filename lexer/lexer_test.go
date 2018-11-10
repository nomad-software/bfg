package lexer

import (
	"testing"

	"github.com/nomad-software/bfg/token"
)

var tokens []token.Token

func BenchmarkLexer(b *testing.B) {
	program := []byte(">++++++++[<+++++++++>-]<.>>+>+>++>[-]+<[>[->+<<++++>]<<]>.+++++++..+++.>>+++++++.<<<[[-]<[-]>]<+++++++++++++++.>>.+++.------.--------.>>+.>++++.")
	var t []token.Token

	for x := 0; x < b.N; x++ {
		t = New(program).Tokens
	}

	tokens = t
}

func TestLexingSingleOperators(t *testing.T) {
	program := []byte("//this is a comment><+-foo,.[]bar")

	tests := []token.Token{
		{token.RightType, ">", 1, 1},
		{token.LeftType, "<", 1, 1},
		{token.AddType, "+", 1, 1},
		{token.SubType, "-", 1, 1},
		{token.InType, ",", 1, 1},
		{token.OutType, ".", 1, 1},
		{token.OpenType, "[", 1, 1},
		{token.CloseType, "]", 1, 1},
		{token.EOFType, "", 0, 0},
	}

	tokens := New(program).Tokens

	for x := 0; x < len(tests); x++ {
		if tokens[x].Type != tests[x].Type || tokens[x].Literal != tests[x].Literal || tokens[x].Shift != tests[x].Shift || tokens[x].Value != tests[x].Value {
			fail(t, tests[x], tokens[x])
		}
	}
}

func TestLexingMultipleOperators(t *testing.T) {
	program := []byte("++++++++>+++>>>>+++<+...---<<-<--,,,.")

	tests := []token.Token{
		{token.AddType, "++++++++", 8, 8},
		{token.RightType, ">", 1, 1},
		{token.AddType, "+++", 3, 3},
		{token.RightType, ">>>>", 4, 4},
		{token.AddType, "+++", 3, 3},
		{token.LeftType, "<", 1, 1},
		{token.AddType, "+", 1, 1},
		{token.OutType, ".", 1, 1},
		{token.OutType, ".", 1, 1},
		{token.OutType, ".", 1, 1},
		{token.SubType, "---", 3, 3},
		{token.LeftType, "<<", 2, 2},
		{token.SubType, "-", 1, 1},
		{token.LeftType, "<", 1, 1},
		{token.SubType, "--", 2, 2},
		{token.InType, ",", 1, 1},
		{token.InType, ",", 1, 1},
		{token.InType, ",", 1, 1},
		{token.OutType, ".", 1, 1},
		{token.EOFType, "", 0, 0},
	}

	tokens := New(program).Tokens

	for x := 0; x < len(tests); x++ {
		if tokens[x].Type != tests[x].Type || tokens[x].Literal != tests[x].Literal || tokens[x].Shift != tests[x].Shift || tokens[x].Value != tests[x].Value {
			fail(t, tests[x], tokens[x])
		}
	}
}

func fail(t *testing.T, a token.Token, b token.Token) {
	t.Errorf("Expected: %#v", a)
	t.Fatalf("Actual  : %#v", b)
}
