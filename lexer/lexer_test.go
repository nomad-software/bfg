package lexer

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/nomad-software/bfg/token"
)

var tokens []token.Token

func BenchmarkMandelbrot(b *testing.B) {
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

	b.ReportMetric(float64(len(t)), "tokens")

	tokens = t
}

func BenchmarkLong(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := os.ReadFile(path.Join(wd, "../programs/long.bf"))
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

	b.ReportMetric(float64(len(t)), "tokens")

	tokens = t
}

func BenchmarkHanoi(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := os.ReadFile(path.Join(wd, "../programs/hanoi.bf"))
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

	b.ReportMetric(float64(len(t)), "tokens")

	tokens = t
}

func Benchmark99Bottles(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := os.ReadFile(path.Join(wd, "../programs/99bottles.bf"))
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

	b.ReportMetric(float64(len(t)), "tokens")

	tokens = t
}

func BenchmarkSierpinski(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := os.ReadFile(path.Join(wd, "../programs/sierpinski.bf"))
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

	b.ReportMetric(float64(len(t)), "tokens")

	tokens = t
}

func TestLexingSingleOperators(t *testing.T) {
	program := []byte("//this is a comment><+-foo,.[]bar")

	tokens := []token.Token{
		{Type: token.RightType, Move: 1, Value: 0},
		{Type: token.LeftType, Move: 1, Value: 0},
		{Type: token.AddType, Move: 0, Value: 1},
		{Type: token.SubType, Move: 0, Value: 1},
		{Type: token.InType, Move: 0, Value: 0},
		{Type: token.OutType, Move: 0, Value: 0},
		{Type: token.OpenType, Move: 0, Value: 0, Jump: 7},
		{Type: token.CloseType, Move: 0, Value: 0, Jump: 6},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultipleOperators(t *testing.T) {
	program := []byte("++++++++>+++>>>>+++<+...---<<-<--,,,.")

	tokens := []token.Token{
		{Type: token.AddType, Move: 0, Value: 8},
		{Type: token.RightType, Move: 1, Value: 0},
		{Type: token.AddType, Move: 0, Value: 3},
		{Type: token.RightType, Move: 4, Value: 0},
		{Type: token.AddType, Move: 0, Value: 3},
		{Type: token.LeftType, Move: 1, Value: 0},
		{Type: token.AddType, Move: 0, Value: 1},
		{Type: token.OutType, Move: 0, Value: 0},
		{Type: token.OutType, Move: 0, Value: 0},
		{Type: token.OutType, Move: 0, Value: 0},
		{Type: token.SubType, Move: 0, Value: 3},
		{Type: token.LeftType, Move: 2, Value: 0},
		{Type: token.SubType, Move: 0, Value: 1},
		{Type: token.LeftType, Move: 1, Value: 0},
		{Type: token.SubType, Move: 0, Value: 2},
		{Type: token.InType, Move: 0, Value: 0},
		{Type: token.InType, Move: 0, Value: 0},
		{Type: token.InType, Move: 0, Value: 0},
		{Type: token.OutType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingZeroOptimisation(t *testing.T) {
	program := []byte("++++++++++[-]++++++++++[-]-----[-]")

	tokens := []token.Token{
		{Type: token.AddType, Move: 0, Value: 10},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.AddType, Move: 0, Value: 10},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.SubType, Move: 0, Value: 5},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultiplicationOptimisation1(t *testing.T) {
	program := []byte("[->+++>+++++++<<]")

	tokens := []token.Token{
		{Type: token.MulAddType, Move: 1, Value: 3},
		{Type: token.MulAddType, Move: 2, Value: 7},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultiplicationOptimisation2(t *testing.T) {
	program := []byte("[->>>>>>>+<<+<<<+<<]")

	tokens := []token.Token{
		{Type: token.MulAddType, Move: 7, Value: 1},
		{Type: token.MulAddType, Move: 5, Value: 1},
		{Type: token.MulAddType, Move: 2, Value: 1},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultiplicationOptimisation3(t *testing.T) {
	program := []byte("[-<<<+>>>]")

	tokens := []token.Token{
		{Type: token.MulAddType, Move: -3, Value: 1},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultiplicationOptimisation4(t *testing.T) {
	program := []byte("[->--->-------<<]")

	tokens := []token.Token{
		{Type: token.MulSubType, Move: 1, Value: 3},
		{Type: token.MulSubType, Move: 2, Value: 7},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultiplicationOptimisation5(t *testing.T) {
	program := []byte("[->>>>>>>-<<-<<<-<<]")

	tokens := []token.Token{
		{Type: token.MulSubType, Move: 7, Value: 1},
		{Type: token.MulSubType, Move: 5, Value: 1},
		{Type: token.MulSubType, Move: 2, Value: 1},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func TestLexingMultiplicationOptimisation6(t *testing.T) {
	program := []byte("[-<<<->>>]")

	tokens := []token.Token{
		{Type: token.MulSubType, Move: -3, Value: 1},
		{Type: token.ZeroType, Move: 0, Value: 0},
		{Type: token.EOFType},
	}

	assertTokens(t, program, tokens)
}

func assertTokens(t *testing.T, program []byte, tokens []token.Token) {
	t.Helper()

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
