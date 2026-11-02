package parser

import (
	"github.com/dyxgou/lox/pkg/ast"
	"github.com/dyxgou/lox/pkg/token"
)

func ParseFile(file *token.File, src []byte, mode Mode) (*ast.File, error) {
	if file == nil {
		panic("parser.ParseFile: no token.File provided (file == nil)")
	}

	var p parser
	p.init(file, src, mode)

	f, err := p.parseFile()
	if err != nil {
		return f, err
	}

	f.FileStart = token.Pos(file.Base())
	f.FileEnd = file.End()

	return f, nil
}

func ParseExpr(src []byte, mode Mode) ast.Expr {
	var file token.File
	file.Init("expr", 1, len(src))

	var p parser
	p.init(&file, src, mode)

	return p.parseExpr()
}
