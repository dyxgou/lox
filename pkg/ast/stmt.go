package ast

import "github.com/dyxgou/lox/pkg/token"

type (
	BadStmt struct {
		From, To token.Pos
	}

	EmptyStmt struct {
		Semicolon token.Pos // position of ';'
		Implicit  bool      // if set, ';' was omitted in the source code
	}

	// ExprStmt represents a stand-alone expression on a statement list
	ExprStmt struct {
		X Expr // expression
	}

	BlockStmt struct {
		Lbrace token.Pos // position of '{'
		List   []Stmt
		Rbrace token.Pos // position of '}'
	}

	AssignStmt struct {
		Lhs    Expr
		TokPos token.Pos
		Tok    token.Token
		Rhs    Expr
	}

	IfStmt struct {
		If   token.Pos // position of 'if' keyword
		Cond Expr      // condition
		Body *BlockStmt
		Else Stmt // else branch or nil
	}
)

func (s *BadStmt) Pos() token.Pos { return s.From }
func (s *BadStmt) End() token.Pos { return s.To }

func (s *EmptyStmt) Pos() token.Pos { return s.Semicolon }
func (s *EmptyStmt) End() token.Pos { return s.Semicolon }

func (s *ExprStmt) Pos() token.Pos { return s.X.Pos() }
func (s *ExprStmt) End() token.Pos { return s.X.End() }

func (s *BlockStmt) Pos() token.Pos { return s.Lbrace }
func (s *BlockStmt) End() token.Pos { return s.Rbrace }

func (s *AssignStmt) Pos() token.Pos { return s.Lhs.Pos() }
func (s *AssignStmt) End() token.Pos { return s.Rhs.End() }

func (s *IfStmt) Pos() token.Pos { return s.If }
func (s *IfStmt) End() token.Pos {
	if s.Else != nil {
		return s.End()
	}

	return s.Body.End()
}

func (s *BadStmt) stmtNode()    {}
func (s *EmptyStmt) stmtNode()  {}
func (s *ExprStmt) stmtNode()   {}
func (s *BlockStmt) stmtNode()  {}
func (s *AssignStmt) stmtNode() {}
func (s *IfStmt) stmtNode()     {}
