package ast

import "github.com/dyxgou/lox/pkg/token"

type (
	Field struct {
		Name *Ident
		Type Expr
	}

	FieldList struct {
		Opening token.Pos
		List    []*Field
		Closing token.Pos
	}

	FuncType struct {
		Func    token.Pos  // positions of the "func" keyword
		Params  *FieldList // func parameters or nil
		Returns *FieldList // func return types or nil
	}

	Comment struct {
		Slash token.Pos // position of '//'
		Text  string
	}
)

func (f *Field) Pos() token.Pos {
	if f.Name != nil {
		return f.Name.Pos()
	}

	if f.Type != nil {
		return f.Type.Pos()
	}

	return token.NoPos
}

func (f *Field) End() token.Pos {
	if f.Type != nil {
		return f.Type.End()
	}
	if f.Name != nil {
		return f.Name.End()
	}

	return token.NoPos
}

func (f *FieldList) Pos() token.Pos {
	if f.Opening.IsValid() {
		return f.Opening
	}
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}

	return token.NoPos
}

func (f *FieldList) End() token.Pos {
	if f.Closing.IsValid() {
		return f.Closing
	}

	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}

	return token.NoPos
}
