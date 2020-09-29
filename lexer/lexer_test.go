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
		{Type: token.RightType, Shift: 1, Value: 1},
		{Type: token.LeftType, Shift: 1, Value: 1},
		{Type: token.AddType, Shift: 1, Value: 1},
		{Type: token.SubType, Shift: 1, Value: 1},
		{Type: token.InType, Shift: 1, Value: 1},
		{Type: token.OutType, Shift: 1, Value: 1},
		{Type: token.OpenType, Shift: 1, Value: 1, Jump: 7},
		{Type: token.CloseType, Shift: 1, Value: 1, Jump: 6},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultipleOperators(t *testing.T) {
	program := []byte("++++++++>+++>>>>+++<+...---<<-<--,,,.")

	tokens := []token.Token{
		{Type: token.AddType, Shift: 8, Value: 8},
		{Type: token.RightType, Shift: 1, Value: 1},
		{Type: token.AddType, Shift: 3, Value: 3},
		{Type: token.RightType, Shift: 4, Value: 4},
		{Type: token.AddType, Shift: 3, Value: 3},
		{Type: token.LeftType, Shift: 1, Value: 1},
		{Type: token.AddType, Shift: 1, Value: 1},
		{Type: token.OutType, Shift: 1, Value: 1},
		{Type: token.OutType, Shift: 1, Value: 1},
		{Type: token.OutType, Shift: 1, Value: 1},
		{Type: token.SubType, Shift: 3, Value: 3},
		{Type: token.LeftType, Shift: 2, Value: 2},
		{Type: token.SubType, Shift: 1, Value: 1},
		{Type: token.LeftType, Shift: 1, Value: 1},
		{Type: token.SubType, Shift: 2, Value: 2},
		{Type: token.InType, Shift: 1, Value: 1},
		{Type: token.InType, Shift: 1, Value: 1},
		{Type: token.InType, Shift: 1, Value: 1},
		{Type: token.OutType, Shift: 1, Value: 1},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingZeroOptimisation(t *testing.T) {
	program := []byte("++++++++++[-]++++++++++[-]+++++[->+<]")

	tokens := []token.Token{
		{Type: token.AddType, Shift: 10, Value: 10},
		{Type: token.ZeroType, Shift: 3, Value: 3},
		{Type: token.AddType, Shift: 10, Value: 10},
		{Type: token.ZeroType, Shift: 3, Value: 3},
		{Type: token.AddType, Shift: 5, Value: 5},
		{Type: token.CopyType, Shift: 1},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingCopyOptimisation(t *testing.T) {
	program := []byte("++++++++++[->+>+>+<<<][->+>+>+>+>+<<<<<]")

	tokens := []token.Token{
		{Type: token.AddType, Shift: 10, Value: 10},
		{Type: token.CopyType, Shift: 3},
		{Type: token.CopyType, Shift: 5},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func assertTokens(t *testing.T, program []byte, tokens []token.Token) {
	output := New(program).Tokens

	for x := 0; x < len(tokens); x++ {
		typeMismatch := output[x].Type != tokens[x].Type
		shiftMismatch := output[x].Shift != tokens[x].Shift
		valueMismatch := output[x].Value != tokens[x].Value
		jumpMismatch := output[x].Jump != tokens[x].Jump

		if typeMismatch || shiftMismatch || valueMismatch || jumpMismatch {
			t.Errorf("Expected: %#v", tokens[x])
			t.Fatalf("Actual  : %#v", output[x])
		}
	}
}
