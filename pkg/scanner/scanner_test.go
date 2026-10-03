package scanner

import (
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
