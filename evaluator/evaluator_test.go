package evaluator

import (
	"bufio"
	"io"
	"log"
	"os"
	"path"
	"testing"

	"github.com/nomad-software/bfg/lexer"
)

func BenchmarkMandelbrot(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := os.ReadFile(path.Join(wd, "../programs/mandelbrot.bf"))
	if err != nil {
		log.Fatalln(err)
	}

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(io.Discard)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		Evaluate(tokens, input, output)
	}

	b.ReportMetric(float64(len(tokens)), "tokens")
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

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(io.Discard)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		Evaluate(tokens, input, output)
	}

	b.ReportMetric(float64(len(tokens)), "tokens")
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

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(io.Discard)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		Evaluate(tokens, input, output)
	}

	b.ReportMetric(float64(len(tokens)), "tokens")
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

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(io.Discard)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		Evaluate(tokens, input, output)
	}

	b.ReportMetric(float64(len(tokens)), "tokens")
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

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(io.Discard)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		Evaluate(tokens, input, output)
	}

	b.ReportMetric(float64(len(tokens)), "tokens")
}
