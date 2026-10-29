package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dyxgou/lox/pkg/ast"
	"github.com/dyxgou/lox/pkg/scanner"
	"github.com/dyxgou/lox/pkg/token"
)

var ErrScanFailed = errors.New("failed scanning")

// parser holds the internal structure of the Parser
type parser struct {
	file    *token.File
	errors  scanner.ErrorList
	scanner scanner.Scanner

	// nextTok
	pos token.Pos
	tok token.Token
	lit string
}

func (p *parser) init(file *token.File, src []byte) {
	p.file = file
	eh := func(pos token.Position, msg string) {
		p.errors.Add(pos, msg)
	}

	p.scanner.Init(p.file, src, eh)
	p.next()
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *parser) parseFile() (*ast.File, error) {
	if p.errors.Len() != 0 {
		return nil, ErrScanFailed
	}

	p.expect(token.PACKAGE)
	ident := p.parseIdent()

	if ident.Name == "_" {
		p.error(p.pos, "invalid package name '_'")
	}
	p.expectSemi()

	if p.errors.Len() != 0 {
		slog.Error("error parsing", "e", p.errors[0].Error())
		return nil, errors.New("failed parsing package notation")
	}

	decls := make([]ast.Decl, 0, 10)
	if p.tok == token.IMPORT {
		panic("TODO: Implement import parsing")
	}

	prev := token.IMPORT
	for p.tok != token.EOF {
		if p.tok == token.SEMICOLON {
			// TODO: Implement advance til next meaningful token
			p.next()
			continue
		}

		if p.tok == token.IMPORT && prev != token.IMPORT {
			panic("TODO: Implement import parsing")
		}
		prev = p.tok

		decls = append(decls, p.parseDecl())
	}

	return &ast.File{Decls: decls}, nil
}

func (p *parser) parseDecl() ast.Decl {
	switch p.tok {
	case token.FUNC:
		return p.parseFuncDecl()
	default:
		p.errorExpected(p.pos, "declaration")
		p.next()
	}

	return nil
}

func (p *parser) parseFuncDecl() *ast.FuncDecl {
	pos := p.expect(token.FUNC)

	ident := p.parseIdent()
	params := p.parseParameters()
	rets := p.parseTypeList()
	body := p.parseBody()

	p.expectSemi()

	return &ast.FuncDecl{
		Name: ident,
		Type: &ast.FuncType{
			Func:    pos,
			Params:  params,
			Returns: rets,
		},
		Body: body,
	}
}

func (p *parser) parseParameters() *ast.FieldList {
	if p.tok == token.LPAREN {
		lparen := p.expect(token.LPAREN)

		fields := p.parseParametersList(token.RPAREN)

		rparen := p.expect(token.RPAREN)

		return &ast.FieldList{Opening: lparen, List: fields, Closing: rparen}
	}

	return nil
}

func (p *parser) parseParametersList(closing token.Token) []*ast.Field {
	params := make([]*ast.Field, 0, 5)

	for p.tok != closing && p.tok != token.EOF {
		if p.tok != token.IDENT {
			p.expect(token.IDENT)
		}

		ident := ast.NewIdent(p.lit)
		ident.NamePos = p.pos

		p.next()

		typ := p.parseType()
		f := &ast.Field{Name: ident, Type: typ}

		params = append(params, f)

		if !p.atComma("parameter list", closing) {
			break
		}

		// Consumes the Comma if encounter
		p.next()
	}

	return params
}

func (p *parser) parseTypeList() *ast.FieldList {
	typs := make([]*ast.Field, 0, 1)

	if p.tok == token.LPAREN {
		lparen := p.expect(token.LPAREN)

		for p.tok != token.RPAREN && p.tok != token.EOF {
			f := &ast.Field{Type: p.parseType()}
			typs = append(typs, f)

			if !p.atComma("func return type list", token.RPAREN) {
				break
			}

			p.next()
		}

		rparen := p.expect(token.RPAREN)

		return &ast.FieldList{Opening: lparen, List: typs, Closing: rparen}
	}

	if typ := p.tryIdentOrType(); typ != nil {
		f := &ast.Field{Type: typ}
		typs = append(typs, f)
		return &ast.FieldList{List: typs}
	}

	return nil
}

// atComma returns true if the current Token is a comma (',').
//
// context determines the execution context of this function in case that there's a
// syntax error
//
// follow is the token expected if the comma is not found
func (p *parser) atComma(context string, follow token.Token) bool {
	if p.tok == token.COMMA {
		return true
	}

	if p.tok != follow {
		var msg strings.Builder
		msg.WriteString("missing ','")

		if p.tok == token.SEMICOLON && p.lit == "\n" {
			msg.WriteString(" before new line")
		}

		msg.WriteString(" in ")
		msg.WriteString(context)
		p.error(p.pos, msg.String())

		// 'insert' a comma and continue to not stop parsing
		return true
	}

	return false
}

func (p *parser) parseType() ast.Expr {
	typ := p.tryIdentOrType()

	if typ == nil {
		pos := p.pos
		p.errorExpected(p.pos, "type")
		return &ast.BadExpr{From: pos, To: p.pos}
	}

	return typ
}

func (p *parser) tryIdentOrType() (typ ast.Expr) {
	switch p.tok {
	case token.IDENT:
		typ = p.parseIdent()
	}

	return
}

func (p *parser) parseBody() *ast.BlockStmt {
	lbrace := p.expect(token.LBRACE)
	list := p.parseStmtList()
	rbrace := p.expect(token.RBRACE)

	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}

func (p *parser) parseStmtList() []ast.Stmt {
	list := make([]ast.Stmt, 0, 10)

	switch p.tok {
	}

	return list
}

func (p *parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"

	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(token.IDENT)
	}

	return &ast.Ident{Name: name, NamePos: pos}
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos

	if p.tok != tok {
		var msg strings.Builder
		msg.WriteByte('\'')
		msg.WriteString(tok.String())
		msg.WriteByte('\'')

		p.errorExpected(pos, msg.String())
	}

	p.next()
	return pos
}

func (p *parser) expectSemi() {
	if p.tok != token.SEMICOLON {
		p.errorExpected(p.pos, "';'")
	}

	p.next()
}

func (p *parser) errorExpected(pos token.Pos, msg string) {
	var sb strings.Builder
	sb.WriteString("expected ")
	sb.WriteString(msg)

	if pos == p.pos {
		// the error happened at the current position
		// so we can be more specific
		switch {
		case p.tok == token.SEMICOLON && p.lit == "\n":
			sb.WriteString(", found new line")
		case p.tok.IsLiteral():
			sb.WriteString(", found ")
			sb.WriteString(p.lit)
		default:
			sb.WriteString(", found '")
			sb.WriteString(p.tok.String())
			sb.WriteByte('\'')
		}
	}

	p.error(pos, sb.String())
}

func (p *parser) error(pos token.Pos, msg string) {
	fpos := p.file.Position(pos)
	n := p.errors.Len()

	if n > 0 && p.errors[n-1].Pos.Line == fpos.Line {
		// Skip error if there's already one the same line
		return
	}

	if n > 10 {
		panic(fmt.Errorf("too many errors parsing: %d errors", n))
	}

	p.errors.Add(fpos, msg)
}
