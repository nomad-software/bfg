package lexer

import (
	"io/ioutil"
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

	program, err := ioutil.ReadFile(path.Join(wd, "../programs/mandelbrot.bf"))
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
		{token.RightType, ">", 1, 1, 0},
		{token.LeftType, "<", 1, 1, 0},
		{token.AddType, "+", 1, 1, 0},
		{token.SubType, "-", 1, 1, 0},
		{token.InType, ",", 1, 1, 0},
		{token.OutType, ".", 1, 1, 0},
		{token.OpenType, "[", 1, 1, 7},
		{token.CloseType, "]", 1, 1, 6},
		{token.EOFType, "", 0, 0, 0},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultipleOperators(t *testing.T) {
	program := []byte("++++++++>+++>>>>+++<+...---<<-<--,,,.")

	tokens := []token.Token{
		{token.AddType, "++++++++", 8, 8, 0},
		{token.RightType, ">", 1, 1, 0},
		{token.AddType, "+++", 3, 3, 0},
		{token.RightType, ">>>>", 4, 4, 0},
		{token.AddType, "+++", 3, 3, 0},
		{token.LeftType, "<", 1, 1, 0},
		{token.AddType, "+", 1, 1, 0},
		{token.OutType, ".", 1, 1, 0},
		{token.OutType, ".", 1, 1, 0},
		{token.OutType, ".", 1, 1, 0},
		{token.SubType, "---", 3, 3, 0},
		{token.LeftType, "<<", 2, 2, 0},
		{token.SubType, "-", 1, 1, 0},
		{token.LeftType, "<", 1, 1, 0},
		{token.SubType, "--", 2, 2, 0},
		{token.InType, ",", 1, 1, 0},
		{token.InType, ",", 1, 1, 0},
		{token.InType, ",", 1, 1, 0},
		{token.OutType, ".", 1, 1, 0},
		{token.EOFType, "", 0, 0, 0},
	}

	assertTokens(t, program, tokens)
}

func TestLexingZeroOptimisation(t *testing.T) {
	program := []byte("++++++++++[-]++++++++++[-]+++++[->+<]")

	tokens := []token.Token{
		{token.AddType, "++++++++++", 10, 10, 0},
		{token.ZeroType, "[-]", 3, 3, 0},
		{token.AddType, "++++++++++", 10, 10, 0},
		{token.ZeroType, "[-]", 3, 3, 0},
		{token.AddType, "+++++", 5, 5, 0},
		{token.OpenType, "[", 1, 1, 10},
		{token.SubType, "-", 1, 1, 0},
		{token.RightType, ">", 1, 1, 0},
		{token.AddType, "+", 1, 1, 0},
		{token.LeftType, "<", 1, 1, 0},
		{token.CloseType, "]", 1, 1, 5},
		{token.EOFType, "", 0, 0, 0},
	}

	assertTokens(t, program, tokens)
}

func assertTokens(t *testing.T, program []byte, tokens []token.Token) {
	output := New(program).Tokens

	for x := 0; x < len(tokens); x++ {
		typeMismatch := output[x].Type != tokens[x].Type
		literalMismatch := output[x].Literal != tokens[x].Literal
		shiftMismatch := output[x].Shift != tokens[x].Shift
		valueMismatch := output[x].Value != tokens[x].Value
		jumpMismatch := output[x].Jump != tokens[x].Jump

		if typeMismatch || literalMismatch || shiftMismatch || valueMismatch || jumpMismatch {
			t.Errorf("Expected: %#v", tokens[x])
			t.Fatalf("Actual  : %#v", output[x])
		}
	}
}
