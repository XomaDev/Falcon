package main

import (
	"Falcon/ast/blockly"
	"Falcon/ast/common"
	"Falcon/ast/control"
	"Falcon/ast/list"
	"Falcon/ast/logic"
	"Falcon/ast/math"
	"Falcon/ast/text"
	"Falcon/label"
	"Falcon/sugar"
)

type Parser struct {
	Tokens    []label.Token
	curIndex  int
	tokenSize int
}

func NewParser(tokens []label.Token) *Parser {
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
	peek := p.peek()
	switch peek.Quality {
	case label.If:
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

	for p.notEOF() && p.consume(label.Elif) {
		conditions = append(conditions, p.parse())
		bodies = append(bodies, p.body())
	}
	var elseBody []blockly.Expr
	if p.notEOF() && p.consume(label.Else) {
		elseBody = p.body()
	}
	return &control.IfExpr{
		Conditions: conditions,
		Bodies:     bodies,
		ElseBody:   elseBody,
	}
}

func (p *Parser) body() []blockly.Expr {
	p.expect(label.OpenCurly)
	var expressions []blockly.Expr
	for p.notEOF() && !p.isNext(label.CloseCurly) {
		expressions = append(expressions, p.parse())
	}
	p.expect(label.CloseCurly)
	return expressions
}

func (p *Parser) expression() blockly.Expr {
	left := p.element()
	for p.notEOF() {
		op := p.peek()
		if p.consume(label.Operator) {
			newOperand := p.element()
			if mExpr, ok := left.(*math.Expr); ok && mExpr.Operator.Quality == op.Quality {
				// New MathExpr if last operator mismatch
				mExpr.Operands = append(mExpr.Operands, newOperand)
			} else {
				left = &math.Expr{Operands: []blockly.Expr{left, newOperand}, Operator: op}
			}
		} else {
			break
		}
	}
	return left
}

func operatorPrecedence(quality label.Quality) {
	switch quality {

	}
}

func (p *Parser) element() blockly.Expr {
	left := p.term()
	for p.notEOF() {
		peek := p.peek()
		switch {
		case peek.Quality == label.RightArrow:
			left = &math.PropExpr{Where: p.next(), On: left, Name: p.readName()}
			continue
		case peek.Quality == label.Question:
			left = &common.QuestionExp{Where: p.next(), On: left, Question: p.readName()}
			continue
		}
		break
	}
	return left
}

func (p *Parser) term() blockly.Expr {
	token := p.next()

	switch token.Quality {
	case label.OpenSquare:
		return p.list()
	case label.Not:
		return &logic.NotExpr{Expr: p.expression()}
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
	switch peek.Quality {
	case label.OpenCurve:
		return &common.FuncCall{Where: neExpr.Where, Name: neExpr.Name, Args: p.arguments()}
	default:
		peek.Error("Unexpected label!")
	}
	panic("") // unreachable
}

func (p *Parser) list() blockly.Expr {
	var elements []blockly.Expr
	for p.notEOF() {
		elements = append(elements, p.parse())
		if !p.consume(label.Comma) {
			break
		}
	}
	p.expect(label.CloseSquare)
	return &list.Expr{Elements: elements}
}

func (p *Parser) arguments() []blockly.Expr {
	p.expect(label.OpenCurve)
	var arguments []blockly.Expr
	if p.consume(label.CloseCurve) {
		return arguments
	}
	for p.notEOF() {
		arguments = append(arguments, p.parse())
		if !p.consume(label.Comma) {
			break
		}
	}
	p.expect(label.CloseCurve)
	return arguments
}

func (p *Parser) value(token label.Token) blockly.Expr {
	switch token.Quality {
	case label.Number:
		return &math.NumExpr{Content: token.Content}
	case label.Bool:
		return &logic.BoolExpr{Value: token.Content}
	case label.Text:
		return &text.Expr{Content: token.Content}
	case label.Alpha:
		return &common.NameExpr{Where: token, Name: token.Content, Global: false}
	default:
		panic(sugar.Format("Unknown value type '%'", token.Quality.String()))
	}
}

func (p *Parser) readName() *string {
	next := p.next()
	if next.Quality != label.Alpha {
		next.Error("Expected name but got %", next.Quality.String())
	}
	return next.Content
}

func (p *Parser) consume(t label.Quality) bool {
	if p.peek().Quality == t {
		p.skip()
		return true
	}
	return false
}

func (p *Parser) expect(t label.Quality) {
	next := p.next()
	if next.Quality != t {
		next.Error("Expected type '%' but got '%'", t.String(), next.Quality.String())
	}
}

func (p *Parser) isNext(types ...label.Quality) bool {
	got := p.peek().Quality
	for _, t := range types {
		if got == t {
			return true
		}
	}
	return false
}

func (p *Parser) peek() label.Token {
	return p.Tokens[p.curIndex]
}

func (p *Parser) next() label.Token {
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
