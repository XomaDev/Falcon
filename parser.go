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
		newOperand := p.parse()

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
	return left
}

func (p *Parser) term() ast.Expr {
	token := p.next()
	switch token.Type {
	case types.Number:
		return &ast.NumExpr{Content: token.Content}
	case types.Bool:
		return &ast.BoolExpr{Value: token.Content}
	}
	panic(sugar.Format("Unknown value type '%'", token.Type.String()))
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
