package token

import (
	"strconv"
	"strings"
)

type Position struct {
	Filename string

	// Offset of the respective [Token] relative to its absolute position in the file starting at 0.
	Offset int
	// Line of the current [Token] counted from (\n). [Position] is valid if Line > 0.
	Line   int
	Column int // Column of the respective [Token] on its line.
}

func (p *Position) IsValid() bool {
	return p.Line > 0
}

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
func (p *Position) String() string {
	var sb strings.Builder
	sb.WriteString(p.Filename)

	if p.IsValid() {
		if sb.Len() != 0 {
			sb.WriteByte(':')
		}

		sb.WriteString(strconv.Itoa(p.Line))

		if p.Column != 0 {
			sb.WriteByte(':')
			sb.WriteString(strconv.Itoa(p.Column))
		}
	}

	if sb.Len() == 0 {
		sb.WriteByte('-')
	}

	return sb.String()
}

// Pos is a compact representation of the Token position within a source file.
// It can be converted into a [Position] for a better but larger
// representation.
//
// The Pos value for a given file is a number that ranges between [base,
// base+size], where base and size are specified when a [File] is created,
// being base the starting position of a file within a [FileSet]
type Pos int

// NoPos is the zero value for [Pos], NoPos is the smallest possible position
// representing the beginning of any file. The corresponding [Position] value
// for NoPos is the zero value of [Position].
const NoPos Pos = 0

func (pos Pos) IsValid() bool {
	return pos != NoPos
}
