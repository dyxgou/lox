package ast

import "github.com/dyxgou/lox/pkg/token"

type (
	BlockStmt struct {
		Lbrace token.Pos // position of '{'
		List   []Stmt
		Rbrace token.Pos // position of '}'
	}
)

func (s *BlockStmt) stmtNode() {}
