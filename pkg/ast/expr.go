package ast

import (
	"bytes"
	"go/ast"
	tok "go/token"

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
)

func (e *BadExpr) Pos() token.Pos { return e.From }
func (e *BadExpr) End() token.Pos { return e.To }

func (e *Ident) Pos() token.Pos { return e.NamePos }
func (e *Ident) End() token.Pos { return e.NamePos + token.Pos(len(e.Name)) }

func NewIdent(name string) *Ident { return &Ident{Name: name, NamePos: token.NoPos} }

func (e *File) Pos() token.Pos { return e.FileStart }
func (e *File) End() token.Pos { return e.FileEnd }

func (e *File) exprNode()    {}
func (e *Ident) exprNode()   {}
func (e *BadExpr) exprNode() {}

func (e *File) String() string {
	var buf bytes.Buffer

	ast.Fprint(&buf, tok.NewFileSet(), e, ast.NotNilFilter)

	return buf.String()
}
