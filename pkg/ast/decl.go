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
)

func (d *BadDecl) Pos() token.Pos { return d.From }
func (d *BadDecl) End() token.Pos { return d.To }

func (d *BadDecl) declNode() {}

func (d *FuncDecl) Pos() token.Pos { return d.Name.Pos() }
func (d *FuncDecl) End() token.Pos { return d.Type.Returns.End() }

func (d *FuncDecl) declNode() {}
