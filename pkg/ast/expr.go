package ast

import (
	"github.com/dyxgou/lox/pkg/token"
)

type (
	// BadExpr represents an expression that cannot produce a correct expression
	BadExpr struct {
		From, To token.Pos
	}

	File struct {
		Decls []Decl // Top level declarations

		FileStart, FileEnd token.Pos // Start and End of the respective file
	}

	// Ident node represents an indentifier
	Ident struct {
		Name    string
		NamePos token.Pos
	}

	UnaryExpr struct {
		OpPos token.Pos
		Op    token.Token
		X     Expr
	}

	BinaryExpr struct {
		X     Expr
		OpPos token.Pos
		Op    token.Token
		Y     Expr
	}

	PtrExpr struct {
		Star token.Pos
		X    Expr
	}

	BasicLit struct {
		ValuePos token.Pos
		Kind     token.Token
		Value    string
		ValueEnd token.Pos
	}

	ParenExpr struct {
		Lparen token.Pos // position of '('
		X      Expr      // Expression
		Rparen token.Pos // position of ')'
	}
)

func (e *BadExpr) Pos() token.Pos { return e.From }
func (e *BadExpr) End() token.Pos { return e.To }

func (e *Ident) Pos() token.Pos { return e.NamePos }
func (e *Ident) End() token.Pos { return e.NamePos + token.Pos(len(e.Name)) }

func NewIdent(name string) *Ident { return &Ident{Name: name, NamePos: token.NoPos} }

func (e *File) Pos() token.Pos { return e.FileStart }
func (e *File) End() token.Pos { return e.FileEnd }

func (e *UnaryExpr) Pos() token.Pos { return e.OpPos }
func (e *UnaryExpr) End() token.Pos { return e.X.End() }

func (e *BinaryExpr) Pos() token.Pos { return e.X.Pos() }
func (e *BinaryExpr) End() token.Pos { return e.Y.End() }

func (e *PtrExpr) Pos() token.Pos { return e.Star }
func (e *PtrExpr) End() token.Pos { return e.X.End() }

func (e *BasicLit) Pos() token.Pos { return e.ValuePos }
func (e *BasicLit) End() token.Pos { return e.ValueEnd }

func (e *ParenExpr) Pos() token.Pos { return e.Lparen }
func (e *ParenExpr) End() token.Pos { return e.Rparen }

func (*File) exprNode()       {}
func (*Ident) exprNode()      {}
func (*BadExpr) exprNode()    {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*PtrExpr) exprNode()    {}
func (*BasicLit) exprNode()   {}
func (*ParenExpr) exprNode()  {}
