package main

import (
	"Falcon/ast"
	"Falcon/sugar"
	"Falcon/types"
)

type Parser struct {
	Tokens    []types.Token
	curIndex  int
	tokenSize int
}

func NewParser(tokens []types.Token) *Parser {
	return &Parser{
		Tokens:    tokens,
		curIndex:  0,
		tokenSize: len(tokens),
	}
}

func (p *Parser) ParseAll() []ast.Expr {
	var expressions []ast.Expr
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *Parser) parse() ast.Expr {
	token := p.peek()
	switch token.Type {
	default:
		return p.expression()
	}
}

func (p *Parser) expression() ast.Expr {
	left := p.element()
	for p.notEOF() {
		op := p.peek()
		if op.Type != types.Operator {
			break
		}
		p.skip()
		newOperand := p.element()
		if mExpr, ok := left.(*ast.MathExpr); ok && mExpr.Operator.Type == op.Type {
			// New MathExpr if last operator mismatch
			mExpr.Operands = append(mExpr.Operands, newOperand)
		} else {
			left = &ast.MathExpr{Operands: []ast.Expr{left, newOperand}, Operator: op}
		}
	}
	return left
}

func (p *Parser) element() ast.Expr {
	left := p.term()
	for p.notEOF() {
		peek := p.peek()
		if peek.Type != types.RightArrow {
			break
		}
		p.skip()
		left = &ast.PropExpr{Where: peek, On: left, Name: p.readName()}
	}
	return left
}

func (p *Parser) term() ast.Expr {
	token := p.next()

	switch token.Type {
	case types.OpenSquare:
		return p.list()
	}

	value := p.value(token)
	if p.isEOF() {
		return value
	}
	neExpr, ok := value.(*ast.NameExpr)
	if !ok {
		return value
	}
	peek := p.peek()
	switch peek.Type {
	case types.OpenCurve:
		return &ast.FuncCall{Where: neExpr.Where, Name: neExpr.Name, Args: p.arguments()}
	default:
		peek.Error("Unexpected token!")
	}
	panic("") // unreachable
}

func (p *Parser) list() ast.Expr {
	var elements []ast.Expr
	for p.notEOF() {
		elements = append(elements, p.parse())
		if !p.consume(types.Comma) {
			break
		}
	}
	p.expect(types.CloseSquare)
	return &ast.ListExpr{Elements: elements}
}

func (p *Parser) arguments() []ast.Expr {
	p.expect(types.OpenCurve)
	var arguments []ast.Expr
	if p.consume(types.CloseCurve) {
		return arguments
	}
	for p.notEOF() {
		arguments = append(arguments, p.parse())
		if !p.consume(types.Comma) {
			break
		}
	}
	p.expect(types.CloseCurve)
	return arguments
}

func (p *Parser) value(token types.Token) ast.Expr {
	switch token.Type {
	case types.Number:
		return &ast.NumExpr{Content: token.Content}
	case types.Bool:
		return &ast.BoolExpr{Value: token.Content}
	case types.Text:
		return &ast.TextExpr{Content: token.Content}
	case types.Alpha:
		return &ast.NameExpr{Where: token, Name: token.Content, Global: false}
	}
	panic(sugar.Format("Unknown value type '%'", token.Type.String()))
}

func (p *Parser) readName() *string {
	next := p.next()
	if next.Type != types.Alpha {
		next.Error("Expected name but got %", next.Type.String())
	}
	return next.Content
}

func (p *Parser) consume(t types.Type) bool {
	if p.peek().Type == t {
		p.skip()
		return true
	}
	return false
}

func (p *Parser) expect(t types.Type) {
	next := p.next()
	if next.Type != t {
		next.Error("Expected type '%' but got '%'", t.String(), next.Type.String())
	}
}

func (p *Parser) isNext(t types.Type) bool {
	return p.peek().Type == t
}

func (p *Parser) peek() types.Token {
	return p.Tokens[p.curIndex]
}

func (p *Parser) next() types.Token {
	token := p.Tokens[p.curIndex]
	p.curIndex++
	return token
}

func (p *Parser) skip() {
	p.curIndex++
}

func (p *Parser) notEOF() bool {
	return p.curIndex < p.tokenSize
}

func (p *Parser) isEOF() bool {
	return p.curIndex >= p.tokenSize
}
