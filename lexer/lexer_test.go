package lexer

import (
	"testing"

	"github.com/nomad-software/bfg/token"
)

func TestLexingSingleOperators(t *testing.T) {
	program := `//this is a comment><+-foo,.[]bar`

	tests := []token.Token{
		{token.Right, ">", 1},
		{token.Left, "<", 1},
		{token.Add, "+", 1},
		{token.Sub, "-", 1},
		{token.In, ",", 1},
		{token.Out, ".", 1},
		{token.Open, "[", 1},
		{token.Close, "]", 1},
		{token.EOF, "", 1},
	}

	tokens := New(program).Tokens

	for x := 0; x < len(tests); x++ {
		if tokens[x].Type != tests[x].Type || tokens[x].Literal != tests[x].Literal {
			fail(t, tests[x], tokens[x])
		}
	}
}

func TestLexingMultipleOperators(t *testing.T) {
	program := `++++++++>+++>>>>+++<+...---<<-<--,,,.`

	tests := []token.Token{
		{token.Add, "++++++++", 8},
		{token.Right, ">", 1},
		{token.Add, "+++", 3},
		{token.Right, ">>>>", 4},
		{token.Add, "+++", 3},
		{token.Left, "<", 1},
		{token.Add, "+", 1},
		{token.Out, ".", 1},
		{token.Out, ".", 1},
		{token.Out, ".", 1},
		{token.Sub, "---", 3},
		{token.Left, "<<", 2},
		{token.Sub, "-", 1},
		{token.Left, "<", 1},
		{token.Sub, "--", 2},
		{token.In, ",", 1},
		{token.In, ",", 1},
		{token.In, ",", 1},
		{token.Out, ".", 1},
		{token.EOF, "", 1},
	}

	tokens := New(program).Tokens

	for x := 0; x < len(tests); x++ {
		if tokens[x].Type != tests[x].Type || tokens[x].Literal != tests[x].Literal {
			fail(t, tests[x], tokens[x])
		}
	}
}

func fail(t *testing.T, a token.Token, b token.Token) {
	t.Errorf("Expected Type: %q Literal: %q Count: %d", a.Type, a.Literal, a.Count)
	t.Fatalf("Actual   Type: %q Literal: %q Count: %d", b.Type, b.Literal, b.Count)
}
