// Package scanner provides the logic for scanning the Lox programming language tokens
package scanner

import (
	"fmt"
	"log"

	"github.com/dyxgou/lox/pkg/token"
)

type ErrorHandler func(pos token.Pos, msg string)

type Scanner struct {
	file *token.File
	src  []byte
	err  ErrorHandler

	pos     int
	readPos int
	linePos int
	ch      byte
}

const eof = byte(token.EOF)

func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler) {
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}

	*s = Scanner{
		file: file,
		src:  src,
		err:  err,

		ch: ' ',
	}

	s.next()
}

func (s *Scanner) next() {
	if s.readPos < len(s.src) {
		s.pos = s.readPos
		if s.ch == '\n' {
			s.linePos = s.pos
			s.file.AddLine(s.linePos)
		}

		s.ch = s.src[s.readPos]
		s.readPos++
	} else {
		s.pos = len(s.src)
		s.readPos = len(s.src)

		if s.ch == '\n' {
			s.linePos = s.pos
			s.file.AddLine(s.linePos)
		}

		s.ch = eof
	}
}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	s.skipWhitespace()

	pos = s.file.Pos(s.pos)
	tok = token.ILLEGAL

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = token.IDENT
		if len(lit) > 1 {
			tok = token.Lookup(lit)
		}
	case isDecimal(ch) || s.ch == '.' && isDecimal(s.peek()):
		tok, lit = s.scanNumber()
	default:
		s.next()

		switch ch {
		case eof:
			tok = token.EOF
		case '\n':
			return pos, token.SEMICOLON, "\n"
		case '"':
			tok = token.STRING
			// TODO(scanner): implement scan string
		case '`':
			tok = token.STRING
			lit = s.scanRawString()
		case '\'':
			tok = token.CHAR
			lit = s.scanChar()
		case ':':
			tok = s.switch2(token.COLON, token.DEFINE)
		case ',':
			tok = token.COMMA
		case ';':
			tok = token.SEMICOLON
			lit = ";"
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case '[':
			tok = token.LBRACK
		case ']':
			tok = token.RBRACK
		case '{':
			tok = token.LBRACE
		case '}':
			tok = token.RBRACE
		case '.':
			tok = token.PERIOD
		case '+':
			tok = s.switch3(token.ADD, token.ADD_ASSIGN, '+', token.INC)
		case '-':
			tok = s.switch3(token.SUB, token.SUB_ASSIGN, '-', token.DEC)
		case '*':
			tok = s.switch2(token.MUL, token.MUL_ASSIGN)
		case '/':
			tok = s.switch2(token.QUO, token.QUO_ASSIGN)
		case '%':
			tok = s.switch2(token.REM, token.REM_ASSIGN)
		case '&':
			if s.ch == '^' {
				s.next()
				tok = s.switch2(token.AND_NOT, token.AND_NOT_ASSIGN)
			} else {
				tok = s.switch3(token.AND, token.AND_ASSIGN, '&', token.LAND)
			}
		case '^':
			tok = s.switch2(token.XOR, token.XOR_ASSIGN)
		case '|':
			tok = s.switch3(token.OR, token.OR_ASSIGN, '|', token.LOR)
		case '>':
			tok = s.switch4(token.GTR, token.GEQ, '>', token.SHR, token.SHR_ASSIGN)
		case '<':
			tok = s.switch4(token.LSS, token.LEQ, '<', token.SHL, token.SHL_ASSIGN)
		case '=':
			tok = s.switch2(token.ASSIGN, token.EQL)
		case '!':
			tok = s.switch2(token.NOT, token.NEQ)
		}
	}

	return
}

/*
switch2 let us choose between two tokens(tok1, tok2) depending on s.ch.

if s.ch == '=' -> tok2.

else -> tok1.
*/
func (s *Scanner) switch2(tok1, tok2 token.Token) token.Token {
	if s.ch == '=' {
		s.next()
		return tok2
	}

	return tok1
}

/*
switch3 let us choose between three tokens(tok1, tok2, tok3) depending on s.ch and ch2.

else(all conds are false) -> tok1.

if s.ch == '=' -> tok2.

if s.ch == ch2  -> tok3.
*/
func (s *Scanner) switch3(tok1, tok2 token.Token, ch2 byte, tok3 token.Token) token.Token {
	if s.ch == '=' {
		s.next()
		return tok2
	}

	if s.ch == ch2 {
		s.next()
		return tok3
	}

	return tok1
}

/*
switch4 let us choose between four tokens(tok1, tok2, tok3, tok4) depending on s.ch and ch2.

else(all conds are false) -> tok1.

if s.ch == '=' -> tok2.

if s.ch == ch2  -> tok3.

if s.ch == ch2 && s.peek() == '=' -> tok4
*/
func (s *Scanner) switch4(tok1, tok2 token.Token, ch2 byte, tok3, tok4 token.Token) token.Token {
	if s.ch == '=' {
		s.next()
		return tok2
	}

	if s.ch == ch2 {
		s.next()
		if s.ch == '=' {
			s.next()
			return tok4
		}
		return tok3
	}

	return tok1
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\r' || s.ch == '\t' {
		s.next()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDecimal(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) peek() byte {
	if s.readPos < len(s.src) {
		return s.src[s.readPos]
	}

	return eof
}

func (s *Scanner) scanIdentifier() string {
	pos := s.pos

	for rdOffset, b := range s.src[s.readPos:] {
		// Optimization: As we're ringing over the source, we can skip the bound
		// checks of s.next if the following chars are letters or numbers.
		if isLetter(b) || isDecimal(b) {
			continue
		}

		s.readPos += rdOffset
		s.next()
		goto exit
	}

	s.pos = len(s.src)
	s.readPos = len(s.src)
	s.ch = eof

exit:
	return string(s.src[pos:s.pos])
}

func (s *Scanner) scanNumber() (token.Token, string) {
	pos := s.pos
	tok := token.INT

	if s.ch == '.' {
		tok = token.FLOAT
	}

	for rdOffset, b := range s.src[s.readPos:] {
		if isDecimal(b) {
			continue
		}

		if b == '.' {
			if tok == token.FLOAT {
				tok = token.ILLEGAL
				goto exit
			}

			tok = token.FLOAT
			continue
		}

		if b == 'i' {
			tok = token.IMAG
			s.next()
		}

		s.readPos += rdOffset
		s.next()
		goto exit
	}
	s.pos = len(s.src)
	s.readPos = len(s.src)
	s.ch = eof

exit:
	return tok, string(s.src[pos:s.pos])
}

func (s *Scanner) scanChar() string {
	// Char "'" has already been consumed
	pos := s.pos - 1

	len := 0
	for {
		ch := s.ch
		if ch == '\n' || ch == eof {
			log.Fatalf("char quote never closed at pos=%d", pos+len)
		}

		s.next()
		len++

		if ch == '\'' {
			break
		}
	}

	return string(s.src[pos:s.pos])
}

func (s *Scanner) scanRawString() string {
	// char '`' already consumed
	pos := s.pos - 1

	for {
		ch := s.ch

		if ch == eof {
			log.Fatalf("raw string reached eof at=%d", s.pos)
		}

		s.next()
		if ch == '`' {
			break
		}
	}

	return string(s.src[pos:s.pos])
}
