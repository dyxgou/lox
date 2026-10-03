package token

import "fmt"

// File is a representation of a Lox File within a [FileSet]. A File has a
// name, base, size and line offset table which is composed of the offsets of
// the \n symbols within the file
type File struct {
	name string
	base int // base is the starting position of the file within the [FileSet].
	size int // size is the corresponding byte size of the File.

	lines []int
}

func (f *File) Init(name string, base, size int) {
	*f = File{
		name:  name,
		base:  base,
		size:  size,
		lines: make([]int, 1, 50),
	}
}

// String returns a brief description of the File.
func (f *File) String() string {
	return fmt.Sprintf("%s(%d-%d)", f.Name(), f.Base(), f.End())
}

// Name returns the file name of file.
func (f *File) Name() string {
	return f.name
}

// Base returns the base offset of file.
func (f *File) Base() int {
	return f.base
}

// Size returns the size of file f.
func (f *File) Size() int {
	return f.size
}

// End returns the end position of file.
func (f *File) End() Pos {
	return Pos(f.base + f.size)
}

func (f *File) LineCount() int {
	return len(f.lines)
}

// AddLine add the NewLine (\n) offset for a new line.
// The line offset must be larger than the offset for the previous line
// and smaller than the file size; otherwise the line offset is ignored.
func (f *File) AddLine(offset int) {
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

// Lines returns the effective line offset table of the form described by [File.SetLines].
// Callers must not mutate the result.
func (f *File) Lines() []int {
	return f.lines
}

func (f *File) SetLinesForContent(content []byte) {
	f.lines = f.lines[0:0]

	nlPos := 0
	for offset, ch := range content {
		if nlPos >= 0 {
			f.lines = append(f.lines, nlPos)
		}
		// Setting nlPos = -1 allow us to stop registering new lines if we've
		// reached the end of the content
		nlPos = -1
		if ch == '\n' {
			nlPos = offset + 1
		}
	}
}
