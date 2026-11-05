package parser_test

import (
	"os"
	"testing"

	"github.com/dyxgou/lox/pkg/parser"
	"github.com/dyxgou/lox/pkg/token"
)

var src = readFile("./testdata/fib.lox")
var file token.File

func readFile(filename string) []byte {

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return data
}

func BenchmarkParse(b *testing.B) {
	b.SetBytes(int64(len(src)))
	file.Init("fib", 1, len(src))
	for b.Loop() {
		if _, err := parser.ParseFile(&file, src, 0); err != nil {
			b.Fatalf("benchmark failed due to parse error: %s", err)
		}
	}
}
