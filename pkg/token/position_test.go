package token

import (
	"testing"
)

const defaultBase = 1

func TestSetLinesForContent(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		lines   []int
	}{
		{name: "Empty", content: []byte(""), lines: []int{0}},
		{name: "ABC", content: []byte("ab\nc"), lines: []int{0, 3}},
		{name: "NewLineEnd", content: []byte("ab\nc\n"), lines: []int{0, 3}},
		{name: "ABC2", content: []byte("a\nb\nc"), lines: []int{0, 2, 4}},
		{name: "ABC3", content: []byte("a\nb\nc\nd\ne"), lines: []int{0, 2, 4, 6, 8}},
		{name: "Mult New Lines", content: []byte("a\n\nb"), lines: []int{0, 3}},
		{name: "Package1", content: []byte("print 'hello world'\n"), lines: []int{0}},
		{
			name:    "Package2",
			content: []byte("print 'hello world';\nprint 'world hello'\n"),
			lines:   []int{0, 22},
		},
	}

	var f File
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f.Init(tt.name, defaultBase, len(tt.content))
			f.SetLinesForContent(tt.content)

			for _, nl := range f.Lines() {
				if nl == 0 {
					continue
				}

				prev := tt.content[nl-1]
				if prev != '\n' {
					t.Fatalf("SetLinesForContent: have not got line at pos=%d with char=%q", nl-1, prev)
				}
			}
		})
	}
}
