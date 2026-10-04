// Package scanner provides the logic for scanning the Lox programming language tokens
package scanner

import (
	"fmt"

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
	}

	return
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
