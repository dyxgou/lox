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
)

func (s *BadStmt) Pos() token.Pos { return s.From }
func (s *BadStmt) End() token.Pos { return s.To }

func (s *EmptyStmt) Pos() token.Pos { return s.Semicolon }
func (s *EmptyStmt) End() token.Pos { return s.Semicolon }

func (s *BlockStmt) Pos() token.Pos { return s.Lbrace }
func (s *BlockStmt) End() token.Pos { return s.Rbrace }

func (s *AssignStmt) Pos() token.Pos { return s.Lhs.Pos() }
func (s *AssignStmt) End() token.Pos { return s.Rhs.End() }

func (s *BadStmt) stmtNode()    {}
func (s *EmptyStmt) stmtNode()  {}
func (s *BlockStmt) stmtNode()  {}
func (s *AssignStmt) stmtNode() {}
