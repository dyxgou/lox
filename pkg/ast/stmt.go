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

	DeclStmt struct {
		Decl Decl
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

	ForStmt struct {
		For  token.Pos // position of 'for' keyword
		Init Stmt      // initialization statement or nil
		Cond Expr      // condition expression
		Post Stmt      // post iteration statement or nil
		Body *BlockStmt
	}

	ReturnStmt struct {
		Results []Expr    // Returned Value
		Return  token.Pos // position of 'return' keyword
	}

	SwitchStmt struct {
		Switch token.Pos  // position of the 'switch' keyword
		Cond   Expr       // condition
		Body   *BlockStmt // Just CauseClauses
	}

	CaseClause struct {
		Case  token.Pos // position of the 'case' keyword
		Cond  Expr      // cond of the case clause, nil means default
		Colon token.Pos // position of the ':'
		Body  []Stmt
	}

	IncDecStmt struct {
		X      Expr        // expression
		TokPos token.Pos   // position of the tok
		Tok    token.Token // INC or DEC
	}

	// BranchStmt represents a 'continue' and 'break' keyword
	BranchStmt struct {
		Tok    token.Token // keyword token (continue, break)
		TokPos token.Pos   // position of keyword
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

func (s *ForStmt) Pos() token.Pos { return s.For }
func (s *ForStmt) End() token.Pos {
	return s.Body.End()
}

func (s *IncDecStmt) Pos() token.Pos { return s.X.Pos() }
func (s *IncDecStmt) End() token.Pos { return s.TokPos + 2 /* len("++") */ }

func (s *ReturnStmt) Pos() token.Pos { return s.Return }
func (s *ReturnStmt) End() token.Pos {
	if n := len(s.Results); n > 0 {
		return s.Results[n-1].End()
	}

	return s.Return + 6 /* len("return") */
}

func (s *SwitchStmt) Pos() token.Pos { return s.Switch }
func (s *SwitchStmt) End() token.Pos {
	return s.Body.End()
}

func (s *CaseClause) Pos() token.Pos { return s.Case }
func (s *CaseClause) End() token.Pos {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}

	return s.Colon + 1
}

func (s *BranchStmt) Pos() token.Pos { return s.TokPos }
func (s *BranchStmt) End() token.Pos {
	return s.TokPos + token.Pos(len(s.Tok.String()))
}

func (s *DeclStmt) Pos() token.Pos { return s.Decl.Pos() }
func (s *DeclStmt) End() token.Pos { return s.Decl.End() }

func (s *BadStmt) stmtNode()    {}
func (s *EmptyStmt) stmtNode()  {}
func (s *ExprStmt) stmtNode()   {}
func (s *BlockStmt) stmtNode()  {}
func (s *AssignStmt) stmtNode() {}
func (s *IfStmt) stmtNode()     {}
func (s *ForStmt) stmtNode()    {}
func (s *IncDecStmt) stmtNode() {}
func (s *ReturnStmt) stmtNode() {}
func (s *SwitchStmt) stmtNode() {}
func (s *CaseClause) stmtNode() {}
func (s *BranchStmt) stmtNode() {}
func (s *DeclStmt) stmtNode()   {}
