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

	// tracing and modes
	trace  bool
	indent int
	mode   Mode

	// nextTok
	pos token.Pos
	tok token.Token
	lit string

	// nestLev is used to track the depth level in the recursion parsing
	nestLev int

	// sync
	syncPos token.Pos // reached position when advance()
	syncCnt int       // C'nt call to advance without making progress
}

type Mode byte

const (
	TraceMode Mode = 1 << iota
	SkipPackage
)

func (p *parser) init(file *token.File, src []byte, mode Mode) {
	p.file = file
	eh := func(pos token.Position, msg string) {
		p.errors.Add(pos, msg)
	}

	p.mode = mode
	p.trace = mode&TraceMode == 1

	p.scanner.Init(p.file, src, eh)
	p.next()
}

func (p *parser) printTrace(a ...any) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	pos := p.file.Position(p.pos)
	fmt.Printf("%5d:%3d: ", pos.Line, pos.Column)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Print(a...)
	if len(a) > 1 {
		fmt.Printf("tok=%q : lit=%q", p.tok.String(), p.lit)
	}
	fmt.Println()
}

func trace(p *parser, msg string) *parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *parser) {
	p.indent--
	p.printTrace(")")
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

var declStart = map[token.Token]bool{
	token.IMPORT: true,
	token.CONST:  true,
	token.VAR:    true,
}

var stmtStart = map[token.Token]bool{
	token.BREAK:    true,
	token.CONST:    true,
	token.CONTINUE: true,
	token.FOR:      true,
	token.IF:       true,
	token.RETURN:   true,
	token.SWITCH:   true,
	token.VAR:      true,
}

func (p *parser) advance(tm map[token.Token]bool) {
	for ; p.tok != token.EOF; p.next() {
		if tm[p.tok] {
			if p.pos == p.syncPos {
				p.syncCnt++
				return
			}

			if p.pos > p.syncPos {
				p.syncPos = p.pos
				p.syncCnt = 0
				return
			}
		}
	}
}

func (p *parser) parseFile() (*ast.File, error) {
	if p.trace {
		defer un(trace(p, "file"))
	}

	if p.errors.Len() != 0 {
		return nil, ErrScanFailed
	}

	if p.mode&SkipPackage == 0 {
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

func (p *parser) end() token.Pos {
	return p.scanner.End()
}

func (p *parser) parseDecl() ast.Decl {
	if p.trace {
		defer un(trace(p, "Declaration"))
	}

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
	if p.trace {
		defer un(trace(p, "Func Decl"))
	}

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
	if p.trace {
		defer un(trace(p, "Parameters"))
	}

	if p.tok == token.LPAREN {
		lparen := p.expect(token.LPAREN)

		fields := p.parseParametersList(token.RPAREN)

		rparen := p.expect(token.RPAREN)

		return &ast.FieldList{Opening: lparen, List: fields, Closing: rparen}
	}

	return nil
}

func (p *parser) parseParametersList(closing token.Token) []*ast.Field {
	if p.trace {
		defer un(trace(p, "Param List"))
	}

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
	if p.trace {
		defer un(trace(p, "Type List"))
	}

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
	if p.trace {
		defer un(trace(p, "Type"))
	}

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
	if p.trace {
		defer un(trace(p, "Body"))
	}

	lbrace := p.expect(token.LBRACE)
	list := p.parseStmtList()
	rbrace := p.expect(token.RBRACE)

	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}

func (p *parser) parseStmtList() []ast.Stmt {
	if p.trace {
		defer un(trace(p, "Stmt List"))
	}

	list := make([]ast.Stmt, 0, 10)

	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseStmt())
	}

	return list
}

func (p *parser) parseStmt() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case token.IDENT:
		s = p.parseSimpleStmt()
	case token.IF:
		s = p.parseIfStmt()
	case token.SEMICOLON:
		s = &ast.EmptyStmt{Semicolon: p.pos, Implicit: p.lit == "\n"}
		p.next()
	default:
		pos := p.pos
		p.errorExpected(pos, "statement")
		p.advance(stmtStart)
		s = &ast.BadStmt{From: pos, To: p.pos}
	}

	return
}

func (p *parser) parseSimpleStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "Simple Stmt"))
	}

	x := p.parseExpr()

	switch p.tok {
	case token.DEFINE:
		pos, tok := p.pos, p.tok
		p.next()

		y := p.parseBinaryExpr(nil, token.LowestPrec+1)
		return &ast.AssignStmt{Lhs: x, TokPos: pos, Tok: tok, Rhs: y}
	}

	return &ast.ExprStmt{X: x}
}

func (p *parser) parseExpr() ast.Expr {
	if p.trace {
		defer un(trace(p, "Expression"))
	}

	return p.parseBinaryExpr(nil, token.LowestPrec+1)
}

func (p *parser) parseBinaryExpr(x ast.Expr, prec int) ast.Expr {
	if p.trace {
		defer un(trace(p, "Binary Expr"))
	}

	if x == nil {
		x = p.parseUnaryExpr()
	}

	var n int
	defer func() { p.nestLev -= n }()

	for n = 1; ; n++ {
		incNestLev(p) // prevents inf-looping
		op, oprec := p.tok, p.tokPrec()

		if oprec < prec {
			return x
		}

		pos := p.expect(op)
		y := p.parseBinaryExpr(nil, oprec+1)
		x = &ast.BinaryExpr{X: x, OpPos: pos, Op: op, Y: y}
	}
}

func (p *parser) parseUnaryExpr() (e ast.Expr) {
	if p.trace {
		defer un(trace(p, "Unary Expr"))
	}

	switch p.tok {
	case token.IDENT:
		e = p.parseIdent()
	case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
		e = &ast.BasicLit{ValuePos: p.pos, ValueEnd: p.end(), Kind: p.tok, Value: p.lit}
		p.next()
	case token.MUL:
		pos := p.pos
		p.next()
		x := p.parseUnaryExpr()
		e = &ast.PtrExpr{Star: pos, X: x}
	case token.ADD, token.SUB, token.NOT, token.XOR, token.AND:
		pos, op := p.pos, p.tok
		p.next()

		x := p.parseUnaryExpr()
		e = &ast.UnaryExpr{Op: op, OpPos: pos, X: x}
	}

	return
}

func (p *parser) parseIfStmt() *ast.IfStmt {
	if p.trace {
		defer un(trace(p, "IfStmt"))
	}

	pos := p.expect(token.IF)
	cond := p.parseIfCond()
	p.expect(token.LBRACE)
	body := p.parseBody()

	var elseStmt ast.Stmt
	if p.tok == token.ELSE {
		p.next()

		switch p.tok {
		case token.IF:
			elseStmt = p.parseIfStmt()
		case token.LBRACE:
			elseStmt = p.parseBody()
		default:
			p.errorExpected(p.pos, "if statement or block")
			elseStmt = &ast.BadStmt{From: p.pos, To: p.pos}
		}

	} else {
		p.expectSemi()
	}

	return &ast.IfStmt{
		If:   pos,
		Cond: cond,
		Body: body,
		Else: elseStmt,
	}
}

func (p *parser) parseIfCond() ast.Expr {
	if p.trace {
		defer un(trace(p, "IfCond"))
	}

	if p.tok == token.LBRACE {
		p.error(p.pos, "expected cond in if statement")
		return &ast.BadExpr{From: p.pos, To: p.pos}
	}

	return p.makeExpr(p.parseSimpleStmt(), "boolean expression")
}

func (p *parser) makeExpr(s ast.Stmt, want string) ast.Expr {
	if p.trace {
		defer un(trace(p, "Make Expr"))
	}

	if s == nil {
		return nil
	}

	if es, isExpr := s.(*ast.ExprStmt); isExpr {
		return es.X
	}

	found := "simple statement"
	if _, isAss := s.(*ast.AssignStmt); isAss {
		found = "assignment"
	}

	p.error(s.Pos(), fmt.Sprintf("expected=%q, got=%q", want, found))
	return &ast.BadExpr{From: s.Pos(), To: s.End()}
}

func (p *parser) parseIdent() *ast.Ident {
	if p.trace {
		defer un(trace(p, "Ident"))
	}

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

func (p *parser) tokPrec() int {
	return p.tok.Precedence()
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

const maxNestLev = 1e4

func incNestLev(p *parser) *parser {
	p.nestLev++

	if p.nestLev > maxNestLev {
		p.error(p.pos, "exceeded max nesting level")
		panic("exceeded max nesting level")
	}

	return p
}

func decNestLev(p *parser) {
	p.nestLev--
}
