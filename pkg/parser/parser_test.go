package parser_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dyxgou/lox/pkg/parser"
	"github.com/dyxgou/lox/pkg/token"
)

const filePackage = "package main;\n"

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{
			name: "func main",
			in:   "func main() {}",
		},
		{
			name: "params",
			in:   "func params(param1 type1) {}",
		},
		{
			name: "params2",
			in:   "func params(param1 type1, params2 type2) {}",
		},
		{
			name: "ret1",
			in:   "func params() type1 {}",
		},
		{
			name: "ret2",
			in:   "func params() (type1, type2) {}",
		},
	}

	var file token.File

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(filePackage)

			if _, err := buf.WriteString(tt.in); err != nil {
				t.Fatalf("WritingBuffer: %s", err)
			}

			file.Init(tt.name, 1, buf.Len())

			f, err := parser.ParseFile(&file, buf.Bytes())
			if err != nil {
				t.Fatalf("ParseFile: %s", err)
			}

			fmt.Println(f)
		})
	}
}
