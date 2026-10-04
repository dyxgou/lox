package scanner

import (
	"math"
	"testing"

	"github.com/dyxgou/lox/pkg/token"
)

type ExpTok struct {
	tok token.Token
	lit string
}

var tokens = [...]ExpTok{
	{token.COMMENT, "// Code Comment"},
	{token.IDENT, "hello"},
	{token.INT, "0"},
	{token.INT, "1"},
	{token.INT, "2"},
	{token.INT, "3"},
	{token.INT, "4"},
	{token.INT, "5"},
	{token.INT, "6"},
	{token.INT, "7"},
	{token.INT, "8"},
	{token.INT, "9"},
	{token.INT, "12345678910"},
	{token.FLOAT, ".0"},
	{token.FLOAT, "0."},
	{token.FLOAT, "3.14159265"},
	{token.IMAG, "0i"},
	{token.IMAG, "1i"},
	{token.IMAG, "2i"},
	{token.IMAG, "3i"},
	{token.IMAG, "4i"},
	{token.IMAG, "5i"},
	{token.IMAG, "6i"},
	{token.IMAG, "7i"},
	{token.IMAG, "8i"},
	{token.IMAG, "9i"},
	{token.IMAG, "12345678910i"},
	{token.IMAG, ".0i"},
	{token.IMAG, "0.i"},
	{token.IMAG, "3.14159265i"},
	{token.CHAR, "'a'"},
	{token.STRING, "`Hello world`"},

	// Operators and delimiters
	{token.ADD, "+"},
	{token.SUB, "-"},
	{token.MUL, "*"},
	{token.QUO, "/"},
	{token.REM, "%"},

	{token.AND, "&"},
	{token.OR, "|"},
	{token.XOR, "^"},
	{token.SHL, "<<"},
	{token.SHR, ">>"},
	{token.AND_NOT, "&^"},

	{token.ADD_ASSIGN, "+="},
	{token.SUB_ASSIGN, "-="},
	{token.MUL_ASSIGN, "*="},
	{token.QUO_ASSIGN, "/="},
	{token.REM_ASSIGN, "%="},

	{token.AND_ASSIGN, "&="},
	{token.OR_ASSIGN, "|="},
	{token.XOR_ASSIGN, "^="},
	{token.SHL_ASSIGN, "<<="},
	{token.SHR_ASSIGN, ">>="},
	{token.AND_NOT_ASSIGN, "&^="},

	{token.LAND, "&&"},
	{token.LOR, "||"},
	{token.ARROW, "<-"},
	{token.INC, "++"},
	{token.DEC, "--"},

	{token.EQL, "=="},
	{token.LSS, "<"},
	{token.GTR, ">"},
	{token.ASSIGN, "="},
	{token.NOT, "!"},

	{token.NEQ, "!="},
	{token.LEQ, "<="},
	{token.GEQ, ">="},
	{token.DEFINE, ":="},

	{token.LPAREN, "("},
	{token.LBRACK, "["},
	{token.LBRACE, "{"},
	{token.COMMA, ","},
	{token.PERIOD, "."},

	{token.RPAREN, ")"},
	{token.RBRACK, "]"},
	{token.RBRACE, "}"},
	{token.SEMICOLON, ";"},
	{token.COLON, ":"},
}

const whiteSpace = "	\t	\n\n\n"

var source = func() []byte {
	src := make([]byte, 0, 700)

	for _, tok := range tokens {
		src = append(src, tok.lit...)
		src = append(src, whiteSpace...)
	}

	return src
}()

var file token.File
var scanner Scanner

const defaultBase = 1

func TestNext(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
	}{
		{
			name: "EOF",
			src:  []byte(""),
		},
		{
			name: "Hello world",
			src:  []byte("hello world"),
		},
		{
			name: "Lox tokens",
			src:  source,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file.Init(tt.name, defaultBase, len(tt.src))
			scanner.Init(&file, tt.src, nil)

			for _, ch := range tt.src {

				if ch != scanner.ch {
					t.Fatalf("scanner next char expected=%q. got=%q", ch, scanner.ch)
				}

				scanner.next()
			}
		})
	}
}

type keyword struct {
	lit string
	tok token.Token
}

var keywords = func() []keyword {
	kwrds := make([]keyword, 0, math.MaxInt8)

	for i := range math.MaxInt8 {
		if tok := token.Token(i); tok.IsKeyword() {
			kwrds = append(kwrds, keyword{tok.String(), tok})
		}
	}

	return kwrds
}()

func TestScanIdentAndKeywords(t *testing.T) {
	tests := []struct {
		name string
		lit  string
		want string
		tok  token.Token
	}{
		{
			name: "Ident",
			lit:  "foo",
			want: "foo",
			tok:  token.IDENT,
		},
		{
			name: "Single letter ident",
			lit:  "a",
			want: "a",
			tok:  token.IDENT,
		},
		{
			name: "Single letter ident with spaces",
			lit:  "a    ",
			want: "a",
			tok:  token.IDENT,
		},
		{
			name: "Ident with Spaces",
			lit:  "foo     ",
			want: "foo",
			tok:  token.IDENT,
		},
		{
			name: "Ident with Nums",
			lit:  "foo123",
			want: "foo123",
			tok:  token.IDENT,
		},
		{
			name: "Ident with Nums and letters",
			lit:  "foo123bar",
			want: "foo123bar",
			tok:  token.IDENT,
		},
	}

	for _, kwrd := range keywords {
		tests = append(tests, struct {
			name string
			lit  string
			want string
			tok  token.Token
		}{
			name: "Keyword: " + kwrd.lit,
			lit:  kwrd.lit,
			want: kwrd.lit,
			tok:  kwrd.tok,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file.Init(tt.lit, defaultBase, len(tt.lit))
			scanner.Init(&file, []byte(tt.lit), nil)

			_, tok, lit := scanner.Scan()

			if lit != tt.want {
				t.Errorf("scanIdentifier literal expected=%q. got=%q", tt.want, lit)
			}

			if tok != tt.tok {
				t.Errorf("scanIdentifier token expected=%q. got=%q", tt.tok.String(), tok.String())
			}
		})
	}
}
