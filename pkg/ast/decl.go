package ast

import "github.com/dyxgou/lox/pkg/token"

type (
	BadDecl struct {
		From, To token.Pos
	}

	FuncDecl struct {
		Name *Ident
		Type *FuncType
		Body *BlockStmt
	}

	ValueDecl struct {
		TokPos token.Pos   // position of the keyword
		Tok    token.Token // VAR or CONST
		Names  []*Ident    // value names; len(Names) > 0
		Type   Expr        // value type
		Values []Expr      // initial values or nil
	}
)

func (d *BadDecl) Pos() token.Pos { return d.From }
func (d *BadDecl) End() token.Pos { return d.To }

func (d *FuncDecl) Pos() token.Pos { return d.Name.Pos() }
func (d *FuncDecl) End() token.Pos { return d.Type.Returns.End() }

func (d *ValueDecl) Pos() token.Pos { return d.Names[0].Pos() }
func (d *ValueDecl) End() token.Pos {
	if n := len(d.Values); n > 0 {
		return d.Values[n-1].End()
	}

	if d.Type != nil {
		return d.Type.End()
	}

	return d.Names[len(d.Names)-1].End()
}

func (d *FuncDecl) declNode()  {}
func (d *ValueDecl) declNode() {}
func (d *BadDecl) declNode()   {}
