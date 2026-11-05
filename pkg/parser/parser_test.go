package parser_test

import (
	"fmt"
	"iter"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dyxgou/lox/pkg/ast"
	"github.com/dyxgou/lox/pkg/parser"
	"github.com/dyxgou/lox/pkg/token"
)

func readFiles(dir string) (iter.Seq2[string, []byte], error) {
	p := path.Clean(path.Join(dir, "/", "*.lox"))
	fs, err := filepath.Glob(p)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024*3)

	it := func(yield func(string, []byte) bool) {
		for _, file := range fs {
			f, err := os.Open(file)

			if err != nil {
				log.Fatalf("opening file: %s", err)
			}

			n, err := f.Read(buf)
			if err != nil {
				log.Fatalf("reading file error: %s", err)
			}

			if !yield(f.Name(), buf[:n]) {
				return
			}
		}
	}

	return it, nil
}

func TestParseSource(t *testing.T) {
	files, err := readFiles("./testdata/")
	if err != nil {
		t.Fatalf("reading files: %s", err)
	}

	var file token.File

	for name, src := range files {
		name := path.Base(name)

		t.Run(name, func(t *testing.T) {
			file.Init(name, 1, len(src))
			f, err := parser.ParseFile(&file, src, parser.SkipPackage|parser.TraceMode)

			if f != nil {
				t.Logf("file: %s", parser.Print(f))
			}

			if err != nil {
				t.Fatalf("parsing file: %s", err)
			}
		})
	}
}

func TestParseBinaryExpr(t *testing.T) {
	tests := []struct {
		want *ast.BinaryExpr
		in   string
	}{
		{
			in: "1 + 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.ADD,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "1 + 2 + 3",
			want: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: "1", Kind: token.INT},
					Op: token.ADD,
					Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
				},
				Op: token.ADD,
				Y:  &ast.BasicLit{Value: "3", Kind: token.INT},
			},
		},
		{
			in: "1 * 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.MUL,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "1 * 2 + 3",
			want: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: "1", Kind: token.INT},
					Op: token.MUL,
					Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
				},
				Op: token.ADD,
				Y:  &ast.BasicLit{Value: "3", Kind: token.INT},
			},
		},
		{
			in: "1 * 2 * 3",
			want: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: "1", Kind: token.INT},
					Op: token.MUL,
					Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
				},
				Op: token.MUL,
				Y:  &ast.BasicLit{Value: "3", Kind: token.INT},
			},
		},
		{
			in: "1 + 2 * 3",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.ADD,
				Y: &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: "2", Kind: token.INT},
					Op: token.MUL,
					Y:  &ast.BasicLit{Value: "3", Kind: token.INT},
				},
			},
		},
		{
			in: "1 + 2 / 3",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.ADD,
				Y: &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: "2", Kind: token.INT},
					Op: token.QUO,
					Y:  &ast.BasicLit{Value: "3", Kind: token.INT},
				},
			},
		},
		{
			in: "1 == 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.EQL,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "1 == 2 == true",
			want: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: "1", Kind: token.INT},
					Op: token.EQL,
					Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
				},
				Y:  &ast.BasicLit{Value: "true", Kind: token.IDENT},
				Op: token.EQL,
			},
		},
		{
			in: "1 > 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.GTR,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "1 < 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.LSS,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "1 <= 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.LEQ,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "1 >= 2",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Op: token.GEQ,
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
			},
		},
		{
			in: "true || false",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "true", Kind: token.IDENT},
				Op: token.LOR,
				Y:  &ast.BasicLit{Value: "false", Kind: token.IDENT},
			},
		},
		{
			in: "true && false",
			want: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "true", Kind: token.IDENT},
				Op: token.LAND,
				Y:  &ast.BasicLit{Value: "false", Kind: token.IDENT},
			},
		},
	}

	var file token.File
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			file.Init(tt.in, 1, len(tt.in))
			exp := parser.ParseExpr([]byte(tt.in), 0)

			parser.Print(exp)

			if err := assertExpr(exp, tt.want); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func assertExpr(x ast.Expr, want ast.Expr) error {
	switch x := x.(type) {
	case *ast.BasicLit:
		return assertBasicLit(x, want.(*ast.BasicLit))
	case *ast.BinaryExpr:
		return assertBinaryExpr(x, want.(*ast.BinaryExpr))
	}

	return nil
}

func assertBinaryExpr(x ast.Expr, want *ast.BinaryExpr) error {
	bin, ok := x.(*ast.BinaryExpr)
	if !ok {
		return fmt.Errorf("expr type expected=%q. got=%T",
			reflect.TypeFor[*ast.BinaryExpr](), x)
	}

	if bin.Op != want.Op {
		return fmt.Errorf(
			"ast.BinaryExpr operator expected=%q. got=%q", want.Op, bin.Op,
		)
	}

	if err := assertExpr(bin.X, want.X); err != nil {
		return err
	}

	if err := assertExpr(bin.Y, want.Y); err != nil {
		return err
	}

	return nil
}

func assertBasicLit(x ast.Expr, want *ast.BasicLit) error {
	lit, ok := x.(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("ast.Expr type expected=%q. got=%T",
			reflect.TypeFor[*ast.BasicLit](), x)
	}

	if lit.Value != want.Value {
		return fmt.Errorf("literal value expected=%s. got=%s", want.Value, lit.Value)
	}

	if lit.Kind != want.Kind {
		return fmt.Errorf("literal kind expected=%s. got=%s", want.Kind, lit.Kind)
	}

	return nil
}
