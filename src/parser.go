package main

import (
	"Falcon/ast/blockly"
	"Falcon/ast/common"
	"Falcon/ast/control"
	"Falcon/ast/list"
	"Falcon/ast/logic"
	"Falcon/ast/math"
	"Falcon/ast/text"
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

func (p *Parser) ParseAll() []blockly.Expr {
	var expressions []blockly.Expr
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *Parser) parse() blockly.Expr {
	token := p.peek()
	switch token.Type {
	case types.If:
		return p.ifExpr()
	default:
		return p.expression()
	}
}

func (p *Parser) ifExpr() blockly.Expr {
	p.skip()
	var conditions []blockly.Expr
	var bodies [][]blockly.Expr

	conditions = append(conditions, p.parse())
	bodies = append(bodies, p.body())

	for p.notEOF() && p.consume(types.Elif) {
		conditions = append(conditions, p.parse())
		bodies = append(bodies, p.body())
	}
	var elseBody []blockly.Expr
	if p.notEOF() && p.consume(types.Else) {
		elseBody = p.body()
	}
	return &control.IfExpr{
		Conditions: conditions,
		Bodies:     bodies,
		ElseBody:   elseBody,
	}
}

func (p *Parser) body() []blockly.Expr {
	p.expect(types.OpenCurly)
	var expressions []blockly.Expr
	for p.notEOF() && !p.isNext(types.CloseCurly) {
		expressions = append(expressions, p.parse())
	}
	p.expect(types.CloseCurly)
	return expressions
}

func (p *Parser) expression() blockly.Expr {
	left := p.element()
	for p.notEOF() {
		op := p.peek()
		if op.Type != types.Operator {
			break
		}
		p.skip()
		newOperand := p.element()
		if mExpr, ok := left.(*math.MathExpr); ok && mExpr.Operator.Type == op.Type {
			// New MathExpr if last operator mismatch
			mExpr.Operands = append(mExpr.Operands, newOperand)
		} else {
			left = &math.MathExpr{Operands: []blockly.Expr{left, newOperand}, Operator: op}
		}
	}
	return left
}

func (p *Parser) element() blockly.Expr {
	left := p.term()
	for p.notEOF() {
		peek := p.peek()
		switch {
		case peek.Type == types.RightArrow:
			left = &math.PropExpr{Where: p.next(), On: left, Name: p.readName()}
			continue
		case peek.Type == types.Question:
			left = &common.QuestionExp{Where: p.next(), On: left, Question: p.readName()}
			continue
		}
		break
	}
	return left
}

func (p *Parser) term() blockly.Expr {
	token := p.next()

	switch token.Type {
	case types.OpenSquare:
		return p.list()
	}

	value := p.value(token)
	if p.isEOF() {
		return value
	}
	neExpr, ok := value.(*common.NameExpr)
	if !ok {
		return value
	}
	peek := p.peek()
	switch peek.Type {
	case types.OpenCurve:
		return &common.FuncCall{Where: neExpr.Where, Name: neExpr.Name, Args: p.arguments()}
	default:
		peek.Error("Unexpected token!")
	}
	panic("") // unreachable
}

func (p *Parser) list() blockly.Expr {
	var elements []blockly.Expr
	for p.notEOF() {
		elements = append(elements, p.parse())
		if !p.consume(types.Comma) {
			break
		}
	}
	p.expect(types.CloseSquare)
	return &list.ListExpr{Elements: elements}
}

func (p *Parser) arguments() []blockly.Expr {
	p.expect(types.OpenCurve)
	var arguments []blockly.Expr
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

func (p *Parser) value(token types.Token) blockly.Expr {
	switch token.Type {
	case types.Number:
		return &math.NumExpr{Content: token.Content}
	case types.Bool:
		return &logic.BoolExpr{Value: token.Content}
	case types.Text:
		return &text.TextExpr{Content: token.Content}
	case types.Alpha:
		return &common.NameExpr{Where: token, Name: token.Content, Global: false}
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

func (p *Parser) isNext(types ...types.Type) bool {
	got := p.peek().Type
	for _, t := range types {
		if got == t {
			return true
		}
	}
	return false
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
