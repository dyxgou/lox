package scanner

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dyxgou/lox/pkg/token"
)

type Error struct {
	Msg string
	Pos token.Position
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	var sb strings.Builder

	if e.Pos.Filename != "" || e.Pos.IsValid() {
		sb.WriteString(e.Pos.Filename)
		sb.WriteString(": ")
	}

	sb.WriteString(e.Msg)
	return sb.String()
}

type ErrorList []*Error

func (el ErrorList) Error() string {
	switch len(el) {
	case 0:
		return "no errors"
	case 1:
		return el[0].Error()
	}

	return fmt.Sprintf("%s (and %d more errors)", el[0], len(el)-1)
}

func (el *ErrorList) Add(pos token.Position, msg string) {
	*el = append(*el, &Error{msg, pos})
}

func (el *ErrorList) Reset() {
	*el = (*el)[0:0]
}

func (el ErrorList) Len() int      { return len(el) }
func (el ErrorList) Swap(i, j int) { el[i], el[j] = el[j], el[i] }

func (el ErrorList) Less(i, j int) bool {
	a := &el[i].Pos
	b := &el[j].Pos

	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}

	if a.Line != b.Line {
		return a.Line < b.Line
	}

	if a.Column != b.Column {
		return a.Column < b.Column
	}

	return el[i].Msg < el[j].Msg
}

func (el ErrorList) Sort() {
	sort.Sort(el)
}

func (el ErrorList) Err() error {
	if len(el) == 0 {
		return nil
	}

	return el[0]
}
