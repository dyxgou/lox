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
		if s.ch == '\n' {
			s.linePos = s.pos
			s.file.AddLine(s.linePos)
		}

		s.ch = eof
	}
}

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
