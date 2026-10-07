package scanner

import (
	"math"
	"testing"

	"github.com/dyxgou/lox/pkg/token"
)

type ExpTok struct {
	name string
	lit  string
	want string
	tok  token.Token
}

var tokens = [...]ExpTok{
	{tok: token.IDENT, lit: "hello"},
	{tok: token.INT, lit: "0"},
	{tok: token.INT, lit: "1"},
	{tok: token.INT, lit: "2"},
	{tok: token.INT, lit: "3"},
	{tok: token.INT, lit: "4"},
	{tok: token.INT, lit: "5"},
	{tok: token.INT, lit: "6"},
	{tok: token.INT, lit: "7"},
	{tok: token.INT, lit: "8"},
	{tok: token.INT, lit: "9"},
	{tok: token.INT, lit: "12345678910"},
	{tok: token.FLOAT, lit: ".0"},
	{tok: token.FLOAT, lit: "0."},
	{tok: token.FLOAT, lit: "3.14159265"},
	{tok: token.IMAG, lit: "0i"},
	{tok: token.IMAG, lit: "1i"},
	{tok: token.IMAG, lit: "2i"},
	{tok: token.IMAG, lit: "3i"},
	{tok: token.IMAG, lit: "4i"},
	{tok: token.IMAG, lit: "5i"},
	{tok: token.IMAG, lit: "6i"},
	{tok: token.IMAG, lit: "7i"},
	{tok: token.IMAG, lit: "8i"},
	{tok: token.IMAG, lit: "9i"},
	{tok: token.IMAG, lit: "12345678910i"},
	{tok: token.IMAG, lit: ".0i"},
	{tok: token.IMAG, lit: "0.i"},
	{tok: token.IMAG, lit: "3.14159265i"},
	{tok: token.CHAR, lit: "'a'"},
	{tok: token.STRING, lit: "`Hello world`"},

	// Operators and delimiters
	{tok: token.ADD, lit: "+"},
	{tok: token.SUB, lit: "-"},
	{tok: token.MUL, lit: "*"},
	{tok: token.QUO, lit: "/"},
	{tok: token.REM, lit: "%"},

	{tok: token.AND, lit: "&"},
	{tok: token.OR, lit: "|"},
	{tok: token.XOR, lit: "^"},
	{tok: token.SHL, lit: "<<"},
	{tok: token.SHR, lit: ">>"},
	{tok: token.AND_NOT, lit: "&^"},

	{tok: token.ADD_ASSIGN, lit: "+="},
	{tok: token.SUB_ASSIGN, lit: "-="},
	{tok: token.MUL_ASSIGN, lit: "*="},
	{tok: token.QUO_ASSIGN, lit: "/="},
	{tok: token.REM_ASSIGN, lit: "%="},

	{tok: token.AND_ASSIGN, lit: "&="},
	{tok: token.OR_ASSIGN, lit: "|="},
	{tok: token.XOR_ASSIGN, lit: "^="},
	{tok: token.SHL_ASSIGN, lit: "<<="},
	{tok: token.SHR_ASSIGN, lit: ">>="},
	{tok: token.AND_NOT_ASSIGN, lit: "&^="},

	{tok: token.LAND, lit: "&&"},
	{tok: token.LOR, lit: "||"},
	{tok: token.INC, lit: "++"},
	{tok: token.DEC, lit: "--"},

	{tok: token.EQL, lit: "=="},
	{tok: token.LSS, lit: "<"},
	{tok: token.GTR, lit: ">"},
	{tok: token.ASSIGN, lit: "="},
	{tok: token.NOT, lit: "!"},

	{tok: token.NEQ, lit: "!="},
	{tok: token.LEQ, lit: "<="},
	{tok: token.GEQ, lit: ">="},
	{tok: token.DEFINE, lit: ":="},

	{tok: token.LPAREN, lit: "("},
	{tok: token.LBRACK, lit: "["},
	{tok: token.LBRACE, lit: "{"},
	{tok: token.COMMA, lit: ","},
	{tok: token.PERIOD, lit: "."},

	{tok: token.RPAREN, lit: ")"},
	{tok: token.RBRACK, lit: "]"},
	{tok: token.RBRACE, lit: "}"},
	{tok: token.SEMICOLON, lit: ";"},
	{tok: token.COLON, lit: ":"},
}

const whiteSpace = "	\t\n\n\n"

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

type keywordInfo struct {
	lit string
	tok token.Token
}

var keywords = func() []keywordInfo {
	kwrds := make([]keywordInfo, 0, math.MaxInt8)

	for i := range math.MaxInt8 {
		if tok := token.Token(i); tok.IsKeyword() {
			kwrds = append(kwrds, keywordInfo{tok.String(), tok})
		}
	}

	return kwrds
}()

func TestScanIdentAndKeywords(t *testing.T) {
	tests := []ExpTok{
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
		tests = append(tests, ExpTok{
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

func TestNumberInt(t *testing.T) {
	tests := []ExpTok{
		{
			name: "Int",
			lit:  "1234",
			want: "1234",
			tok:  token.INT,
		},
		{
			name: "Int with spaces",
			lit:  "1234    ",
			want: "1234",
			tok:  token.INT,
		},
		{
			name: "Float",
			lit:  "1234.1234",
			want: "1234.1234",
			tok:  token.FLOAT,
		},
		{
			name: "Float With Spaces",
			lit:  "    1234.1234     ",
			want: "1234.1234",
			tok:  token.FLOAT,
		},
		{
			name: "Float Repeated Dot",
			lit:  "1234..1234",
			tok:  token.ILLEGAL,
		},
		{
			name: "Float dot at end",
			lit:  "1234.",
			want: "1234.",
			tok:  token.FLOAT,
		},
		{
			name: "Float dot at start",
			lit:  ".1234",
			want: ".1234",
			tok:  token.FLOAT,
		},
		{
			name: "Imag",
			lit:  "1234i",
			want: "1234i",
			tok:  token.IMAG,
		},
		{
			name: "Imag with spaces",
			lit:  "	  1234i     ",
			want: "1234i",
			tok:  token.IMAG,
		},
		{
			name: "Imag float",
			lit:  "1234.1234i",
			want: "1234.1234i",
			tok:  token.IMAG,
		},
		{
			name: "Imag float with spaces",
			lit:  "    	1234.1234i    ",
			want: "1234.1234i",
			tok:  token.IMAG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file.Init(tt.lit, defaultBase, len(tt.lit))
			scanner.Init(&file, []byte(tt.lit), nil)

			_, tok, lit := scanner.Scan()
			if tok != tt.tok {
				t.Errorf("scanNumber token expected=%q. got=%q", tt.tok.String(), tok.String())
			}

			if tok != token.ILLEGAL && lit != tt.want {
				t.Errorf("scanNumber literal expected=%q. got=%q", tt.want, lit)
			}
		})
	}
}

func TestScanComments(t *testing.T) {
	tt := struct {
		src []byte
		tok token.Token
		lit string
	}{
		src: []byte("// Lox comment"),
		tok: token.COMMENT,
		lit: " Lox comment",
	}

	file.Init("comment", 1, len(tt.src))
	scanner.Init(&file, tt.src, nil)

	_, tok, lit := scanner.Scan()

	if tok != tt.tok {
		t.Errorf("TestScanComment tok expected=%q. got=%q", tt.tok, tok)
	}

	if lit != tt.lit {
		t.Errorf("TestScanComment literal expected=%q. got=%q", tt.lit, lit)
	}
}

func TestScanSource(t *testing.T) {
	file.Init("tokens_source", defaultBase, len(source))
	scanner.Init(&file, []byte(source), nil)

	lastPos := 0
	for _, tt := range tokens {
		name := "Scanning: " + tt.lit
		t.Run(name, func(t *testing.T) {
			pos, tok, lit := scanner.Scan()

			for tok == token.SEMICOLON && lit == "\n" {
				pos, tok, lit = scanner.Scan()
			}

			expPos := indexFrom(source, tt.lit[0], lastPos) + 1

			if tok != tt.tok {
				t.Errorf("ScanNextTok token expected=%q. got=%q", tt.tok.String(), tok.String())
			}

			if tok.IsLiteral() && lit != tt.lit {
				t.Errorf("ScanNextTok literal expected=%q. got=%q", tt.lit, lit)
			}

			if expPos == -1 {
				t.Fatalf("ScanNextTok lit=%q not found.", tt.lit)
			}

			if int(pos) != expPos {
				t.Fatalf("ScanNextTok pos expected=%d. got=%d", expPos, pos)
			}

			lastPos = int(pos) + len(whiteSpace)
		})
	}
}

func indexFrom(src []byte, x byte, start int) int {
	for i := start; i < len(src); i++ {
		if src[i] == x {
			return i
		}
	}

	return -1
}
