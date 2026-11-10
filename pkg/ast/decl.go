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

	GenDecl struct {
		TokPos token.Pos
		Tok    token.Token // IMPORT, TYPE, VAR, CONST
		Lparen token.Pos   // position of '(' if any
		Specs  []Spec
		Rparen token.Pos // position of ')' if any
	}
)

func (d *BadDecl) Pos() token.Pos { return d.From }
func (d *BadDecl) End() token.Pos { return d.To }

func (d *FuncDecl) Pos() token.Pos { return d.Name.Pos() }
func (d *FuncDecl) End() token.Pos { return d.Type.Returns.End() }

func (d *GenDecl) Pos() token.Pos { return d.TokPos }
func (d *GenDecl) End() token.Pos { return d.Rparen }

func (d *FuncDecl) declNode() {}
func (d *BadDecl) declNode()  {}
func (d *GenDecl) declNode()  {}

type (
	Spec interface {
		Node
		specNode()
	}

	// TypeSpec represents a 'type' decl
	TypeSpec struct {
		Name *Ident // name of the type
		Type Expr   // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	}

	ValueSpec struct {
		Names  []*Ident // value names; len(Names) > 0
		Type   Expr     // value type
		Values []Expr   // initial values or nil
	}
)

func (s *TypeSpec) Pos() token.Pos { return s.Name.Pos() }
func (s *TypeSpec) End() token.Pos { return s.Type.End() }

func (s *ValueSpec) Pos() token.Pos { return s.Names[0].Pos() }
func (s *ValueSpec) End() token.Pos {
	if n := len(s.Values); n > 0 {
		return s.Values[n-1].End()
	}

	if s.Type != nil {
		return s.Type.End()
	}

	return s.Names[len(s.Names)-1].End()
}

func (s *TypeSpec) specNode()  {}
func (s *ValueSpec) specNode() {}
