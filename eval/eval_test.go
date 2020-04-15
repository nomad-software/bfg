package eval

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/nomad-software/bfg/lexer"
)

func BenchmarkEvaluator(b *testing.B) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	program, err := ioutil.ReadFile(path.Join(wd, "../programs/mandelbrot.bf"))
	if err != nil {
		log.Fatalln(err)
	}

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(ioutil.Discard)
	defer output.Flush()

	tokens := lexer.New(program).Tokens

	b.ReportAllocs()
	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		Evaluate(tokens, *input, *output)
	}
}
