package analysis

import (
	blky "Falcon/ast/blockly"
	"Falcon/ast/common"
	"Falcon/ast/control"
	"Falcon/ast/list"
	"Falcon/ast/logic"
	"Falcon/ast/math"
	"Falcon/ast/text"
)
import l "Falcon/lex"

type Parser struct {
	Tokens    []l.Token
	currIndex int
	tokenSize int
}

func NewParser(tokens []l.Token) *Parser {
	return &Parser{
		Tokens:    tokens,
		tokenSize: len(tokens),
		currIndex: 0,
	}
}

func (p *Parser) ParseAll() []blky.Expr {
	var expressions []blky.Expr
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *Parser) parse() blky.Expr {
	switch p.peek().Type {
	case l.If:
		return p.ifExpr()
	default:
		return p.expr(0)
	}
}

func (p *Parser) ifExpr() *control.IfExpr {
	p.skip()
	var conditions []blky.Expr
	var bodies [][]blky.Expr

	conditions = append(conditions, p.expr(0))
	bodies = append(bodies, p.body())

	for p.notEOF() && p.consume(l.Elif) {
		conditions = append(conditions, p.expr(0))
		bodies = append(bodies, p.body())
	}
	var elseBody []blky.Expr
	if p.notEOF() && p.consume(l.Else) {
		elseBody = p.body()
	}
	return &control.IfExpr{
		Conditions: conditions,
		Bodies:     bodies,
		ElseBody:   elseBody,
	}
}

func (p *Parser) body() []blky.Expr {
	p.expect(l.OpenCurly)
	var expressions []blky.Expr
	if p.consume(l.CloseCurly) {
		return expressions
	}
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	p.expect(l.CloseCurly)
	return expressions
}

func (p *Parser) expr(minPrecedence int) blky.Expr {
	left := p.element()
	for p.notEOF() {
		opToken := p.peek()
		if !opToken.HasFlag(l.Operator) {
			break
		}
		precedence := precedenceOf(opToken.Flags[0])
		if precedence == -1 {
			break
		}
		p.skip()
		var right blky.Expr
		if opToken.HasFlag(l.PreserveOrder) {
			right = p.element()
		} else {
			right = p.expr(precedence)
		}
		if mExpr, ok := left.(*math.Expr); ok && mExpr.Operator.Type == opToken.Type {
			mExpr.Operands = append(mExpr.Operands, right)
		} else {
			left = &math.Expr{Operands: []blky.Expr{left, right}, Operator: opToken}
		}
	}
	return left
}

func precedenceOf(flag l.Flag) int {
	switch flag {
	case l.LLogicOr:
		return 1
	case l.LLogicAnd:
		return 2
	case l.BBitwiseOr:
		return 3
	case l.BBitwiseAnd:
		return 4
	case l.BBitwiseXor:
		return 5
	case l.Relational:
		return 6
	case l.Binary:
		return 7
	case l.BinaryL1:
		return 8
	default:
		return -1
	}
}

func (p *Parser) element() blky.Expr {
	left := p.term()
	for p.notEOF() {
		pe := p.peek()
		switch pe.Type {
		case l.RightArrow:
			left = &math.Prop{Where: p.next(), On: left, Name: p.name()}
			continue
		case l.Question:
			left = &common.Question{Where: p.next(), On: left, Question: p.name()}
			continue
		}
		break
	}
	return left
}

func (p *Parser) term() blky.Expr {
	token := p.next()
	switch token.Type {
	case l.OpenSquare:
		return p.list()
	case l.Not:
		return &logic.Not{Expr: p.expr(0)}
	default:
		value := p.value(token)
		if p.isEOF() {
			return value
		}
		nameExpr, ok := value.(*common.Name)
		if !ok {
			return value
		}
		pe := p.peek()
		switch pe.Type {
		case l.OpenCurve:
			return &common.FuncCall{Where: nameExpr.Where, Name: nameExpr.Name, Args: p.arguments()}
		default:
			pe.Error("Unexpected!")
			panic("") // unreachable
		}
	}
}

func (p *Parser) list() *list.Expr {
	var elements []blky.Expr
	if !p.consume(l.CloseSquare) {
		for p.notEOF() {
			elements = append(elements, p.expr(0))
			if !p.consume(l.Comma) {
				break
			}
		}
	}
	p.expect(l.CloseSquare)
	return &list.Expr{Elements: elements}
}

func (p *Parser) arguments() []blky.Expr {
	p.expect(l.OpenCurve)
	var args []blky.Expr
	if p.consume(l.CloseCurve) {
		return args
	}
	for p.notEOF() {
		args = append(args, p.expr(0))
		if !p.consume(l.Comma) {
			break
		}
	}
	p.expect(l.CloseCurve)
	return args
}

func (p *Parser) value(t l.Token) blky.Expr {
	switch t.Type {
	case l.True, l.False:
		return &logic.Bool{Value: t.Type == l.True}
	case l.Number:
		return &math.Num{Content: *t.Content}
	case l.Text:
		return &text.Expr{Content: *t.Content}
	case l.Name:
		return &common.Name{Where: t, Name: *t.Content}
	default:
		t.Error("Unknown value type '%'", t.Type.String())
		panic("") // unreachable
	}
}

func (p *Parser) name() string {
	return *p.expect(l.Name).Content
}

func (p *Parser) consume(t l.Type) bool {
	if p.peek().Type == t {
		p.currIndex++
		return true
	}
	return false
}

func (p *Parser) expect(t l.Type) l.Token {
	got := p.next()
	if got.Type != t {
		got.Error("Expected type % but got %", t.String(), got.Type.String())
	}
	return got
}

func (p *Parser) isNext(t l.Type) bool {
	return p.peek().Type == t
}

func (p *Parser) peek() l.Token {
	return p.Tokens[p.currIndex]
}

func (p *Parser) next() l.Token {
	token := p.Tokens[p.currIndex]
	p.currIndex++
	return token
}

func (p *Parser) skip() {
	p.currIndex++
}

func (p *Parser) notEOF() bool {
	return p.currIndex < p.tokenSize
}

func (p *Parser) isEOF() bool {
	return p.currIndex >= p.tokenSize
}
