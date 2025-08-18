package analysis

import (
	blky "Falcon/ast/blockly"
	"Falcon/ast/common"
	"Falcon/ast/control"
	"Falcon/ast/fundamentals"
	"Falcon/ast/list"
	"Falcon/ast/method"
	"Falcon/ast/procedures"
	"Falcon/ast/variables"
)
import l "Falcon/lex"

type NameResolver struct {
	Procedures map[string]Procedure
}

type Procedure struct {
	Name       string
	Parameters []string
	Returning  bool
}

type LangParser struct {
	Tokens    []*l.Token
	currIndex int
	tokenSize int
	resolver  *NameResolver
}

func NewLangParser(tokens []*l.Token) *LangParser {
	return &LangParser{
		Tokens:    tokens,
		tokenSize: len(tokens),
		currIndex: 0,
		resolver:  &NameResolver{Procedures: map[string]Procedure{}},
	}
}

func (p *LangParser) ParseAll() []blky.Expr {
	var expressions []blky.Expr
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *LangParser) parse() blky.Expr {
	switch p.peek().Type {
	case l.If:
		return p.ifExpr()
	case l.For:
		return p.forExpr()
	case l.Each:
		return p.eachExpr()
	case l.While:
		return p.whileExpr()
	case l.Break:
		p.skip()
		return &control.Break{}
	case l.WalkAll:
		p.skip()
		return &fundamentals.WalkAll{}
	case l.Local:
		return p.varExpr()
	case l.Global:
		return p.globVar()
	case l.Func:
		return p.funcSmt()
	default:
		return p.expr(0)
	}
}

func (p *LangParser) funcSmt() blky.Expr {
	p.skip()
	name := p.name()
	p.expect(l.OpenCurve)
	var parameters []string
	if !p.consume(l.CloseCurve) {
		for p.notEOF() && !p.isNext(l.CloseCurve) {
			parameters = append(parameters, p.name())
			if !p.consume(l.Comma) {
				break
			}
		}
		p.expect(l.CloseCurve)
	}
	returning := p.consume(l.Assign)
	p.resolver.Procedures[name] = Procedure{Name: name, Parameters: parameters, Returning: returning}
	if returning {
		return &procedures.RetProcedure{Name: name, Parameters: parameters, Result: p.parse()}
	} else {
		return &procedures.VoidProcedure{Name: name, Parameters: parameters, Body: p.body()}
	}
}

func (p *LangParser) globVar() blky.Expr {
	p.skip()
	name := p.name()
	p.expect(l.Assign)
	return &variables.Global{Name: name, Value: p.parse()}
}

func (p *LangParser) varExpr() blky.Expr {
	p.skip()

	var varNames []string
	var varValues []blky.Expr
	if p.consume(l.OpenCurve) {
		// a result local var
		for p.notEOF() && !p.isNext(l.CloseCurve) {
			name := p.name()
			p.expect(l.Assign)
			value := p.parse()

			varNames = append(varNames, name)
			varValues = append(varValues, value)

			if !p.consume(l.Comma) {
				break
			}
		}
		p.expect(l.CloseCurve)
		if p.consume(l.RightArrow) {
			return &variables.VarResult{Names: varNames, Values: varValues, Result: p.parse()}
		} else {
			return &variables.Var{Names: varNames, Values: varValues, Body: p.body()}
		}
	} else {
		// a clean full scope variable
		name := p.name()
		p.expect(l.Assign)
		value := p.parse()
		// we gotta parse rest of the body here
		return &variables.SimpleVar{Name: name, Value: value, Body: p.bodyUntilCurly()}
	}
}

func (p *LangParser) whileExpr() *control.While {
	p.skip()
	condition := p.expr(0)
	body := p.body()
	return &control.While{Condition: condition, Body: body}
}

func (p *LangParser) eachExpr() blky.Expr {
	p.skip()
	keyName := p.name()
	if p.consume(l.DoubleColon) {
		// a dictionary pair iteration
		valueName := p.name()
		p.expect(l.RightArrow)
		return &control.EachPair{KeyName: keyName, ValueName: valueName, Iterable: p.element(), Body: p.body()}
	} else {
		// a simple list iteration
		p.expect(l.RightArrow)
		return &control.Each{IName: keyName, Iterable: p.element(), Body: p.body()}
	}
}

func (p *LangParser) forExpr() *control.For {
	p.skip()
	iName := p.name()
	p.expect(l.Colon)
	from := p.expr(0)
	p.expect(l.To)
	to := p.expr(0)
	p.expect(l.By)
	by := p.expr(0)
	body := p.body()
	return &control.For{
		IName: iName,
		From:  from,
		To:    to,
		By:    by,
		Body:  body,
	}
}

func (p *LangParser) ifExpr() blky.Expr {
	p.skip()
	if p.isNext(l.OpenCurve) {
		return p.simpleIf()
	}
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
	return &control.If{Conditions: conditions, Bodies: bodies, ElseBody: elseBody}
}

func (p *LangParser) simpleIf() *control.SimpleIf {
	p.expect(l.OpenCurve)
	condition := p.parse()
	p.expect(l.CloseCurve)
	then := p.parse()
	p.expect(l.Else)
	elze := p.parse()
	return &control.SimpleIf{Condition: condition, Then: then, Else: elze}
}

func (p *LangParser) body() []blky.Expr {
	p.expect(l.OpenCurly)
	expressions := p.bodyUntilCurly()
	return expressions
}

func (p *LangParser) bodyUntilCurly() []blky.Expr {
	var expressions []blky.Expr
	if p.consume(l.CloseCurly) {
		return expressions
	}
	for p.notEOF() && !p.isNext(l.CloseCurly) {
		expressions = append(expressions, p.parse())
	}
	p.expect(l.CloseCurly)
	return expressions
}

func (p *LangParser) expr(minPrecedence int) blky.Expr {
	left := p.element()
	for p.notEOF() {
		opToken := p.peek()
		if !opToken.HasFlag(l.Operator) {
			break
		}
		precedence := precedenceOf(opToken.Flags[0])
		if precedence == -1 || precedence < minPrecedence {
			break
		}
		p.skip()
		var right blky.Expr
		if opToken.HasFlag(l.PreserveOrder) {
			right = p.element()
		} else {
			right = p.expr(precedence)
		}
		if rBinExpr, ok := right.(*common.BinaryExpr); ok && rBinExpr.CanRepeat(opToken.Type) {
			// for NoPreserveOrder: merge binary expr with same operator (towards right)
			rBinExpr.Operands = append([]blky.Expr{left}, rBinExpr.Operands...)
			left = rBinExpr
		} else if lBinExpr, ok := left.(*common.BinaryExpr); ok && lBinExpr.CanRepeat(opToken.Type) {
			// for PreserveOder: merge binary expr with same operator (towards left)
			lBinExpr.Operands = append(lBinExpr.Operands, right)
		} else {
			// a new binary node
			left = &common.BinaryExpr{Where: opToken, Operands: []blky.Expr{left, right}, Operator: opToken.Type}
		}
	}
	return left
}

func precedenceOf(flag l.Flag) int {
	switch flag {
	case l.AssignmentType:
		return 0
	case l.Pair:
		return 1
	case l.TextJoin:
		return 2
	case l.LLogicOr:
		return 3
	case l.LLogicAnd:
		return 4
	case l.BBitwiseOr:
		return 5
	case l.BBitwiseAnd:
		return 6
	case l.BBitwiseXor:
		return 7
	case l.Equality:
		return 8
	case l.Relational:
		return 9
	case l.Binary:
		return 10
	case l.BinaryL1:
		return 11
	case l.BinaryL2:
		return 12
	default:
		return -1
	}
}

func (p *LangParser) element() blky.Expr {
	left := p.term()
	for p.notEOF() {
		pe := p.peek()
		switch pe.Type {
		case l.At:
			left = p.helperDropdown(left)
		case l.Dot:
			left = p.objectCall(left)
			continue
		case l.RightArrow:
			left = &common.Convert{Where: p.next(), On: left, Name: p.name()}
			continue
		case l.Question:
			left = &common.Question{Where: p.next(), On: left, Question: p.name()}
			continue
		case l.DoubleColon:
			// constant value transformer
			left = &common.Transform{Where: p.next(), On: left, Name: p.name()}
		case l.OpenSquare:
			p.skip()
			// an index element access
			left = &list.Get{List: left, Index: p.parse()}
			p.expect(l.CloseSquare)
			continue
		}
		break
	}
	return left
}

func (p *LangParser) helperDropdown(keyExpr blky.Expr) blky.Expr {
	where := p.next()
	if key, ok := keyExpr.(*variables.Get); ok {
		return &fundamentals.HelperDropdown{Key: key.Name, Option: p.name()}
	}
	where.Error("Invalid Helper Access operation ")
	panic("")
}

func (p *LangParser) objectCall(object blky.Expr) blky.Expr {
	p.skip()
	where := p.next()
	name := *where.Content

	var args []blky.Expr
	if p.isNext(l.OpenCurve) {
		args = p.arguments()
	}
	if p.isEOF() || !p.consume(l.OpenCurly) {
		// he's a simple call!
		return &method.Call{Where: where, On: object, Name: name, Args: args}
	}
	// oh, no! he's a transformer >_>
	var namesUsed []string
	if !p.consume(l.RightArrow) {
		for {
			namesUsed = append(namesUsed, p.name())
			if !p.consume(l.Comma) {
				break
			}
		}
		p.consume(l.RightArrow)
	}
	transformer := p.parse()
	p.consume(l.CloseCurly)
	return &list.Transformer{
		Where:       where,
		List:        object,
		Name:        name,
		Args:        args,
		Names:       namesUsed,
		Transformer: transformer}
}

func (p *LangParser) term() blky.Expr {
	token := p.next()
	switch token.Type {
	case l.OpenSquare:
		return p.list()
	case l.OpenCurly:
		return p.dictionary()
	case l.OpenCurve:
		e := p.parse()
		p.expect(l.CloseCurve)
		return e
	case l.Not:
		return &fundamentals.Not{Expr: p.element()}
	case l.Do:
		return p.doExpr()
	case l.If:
		return p.simpleIf()
	default:
		if token.HasFlag(l.Value) {
			value := p.value(token)
			if nameExpr, ok := value.(*variables.Get); ok && p.notEOF() && p.isNext(l.OpenCurve) {
				signature, ok := p.resolver.Procedures[nameExpr.Name]
				if ok {
					return &procedures.Call{
						Name:       nameExpr.Name,
						Parameters: signature.Parameters,
						Arguments:  p.arguments(),
						Returning:  signature.Returning}
				} else {
					return &common.FuncCall{Where: nameExpr.Where, Name: nameExpr.Name, Args: p.arguments()}
				}
			}
			return value
		}
		token.Error("Unexpected! %", token.String())
		panic("") // unreachable
	}
}

func (p *LangParser) doExpr() *control.Do {
	body := p.body()
	p.expect(l.RightArrow)
	result := p.expr(0)
	return &control.Do{Body: body, Result: result}
}

func (p *LangParser) dictionary() *fundamentals.Dictionary {
	var elements []blky.Expr
	if !p.consume(l.CloseCurly) {
		for p.notEOF() {
			elements = append(elements, p.expr(0))
			if !p.consume(l.Comma) {
				break
			}
		}
		p.expect(l.CloseCurly)
	}
	return &fundamentals.Dictionary{Elements: elements}
}

func (p *LangParser) list() *fundamentals.List {
	var elements []blky.Expr
	if !p.consume(l.CloseSquare) {
		for p.notEOF() {
			elements = append(elements, p.expr(0))
			if !p.consume(l.Comma) {
				break
			}
		}
		p.expect(l.CloseSquare)
	}
	return &fundamentals.List{Elements: elements}
}

func (p *LangParser) arguments() []blky.Expr {
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

func (p *LangParser) value(t *l.Token) blky.Expr {
	switch t.Type {
	case l.True, l.False:
		return &fundamentals.Boolean{Value: t.Type == l.True}
	case l.Number:
		return &fundamentals.Number{Content: *t.Content}
	case l.Text:
		return &fundamentals.Text{Content: *t.Content}
	case l.Name:
		return &variables.Get{Where: t, Global: false, Name: *t.Content}
	case l.This:
		p.expect(l.Dot)
		return &variables.Get{Where: t, Global: true, Name: p.name()}
	case l.Color:
		p.expect(l.Colon)
		return &fundamentals.Color{Where: t, Name: p.name()}
	default:
		t.Error("Unknown value type '%'", t.Type.String())
		panic("") // unreachable
	}
}

func (p *LangParser) name() string {
	return *p.expect(l.Name).Content
}

func (p *LangParser) consume(t l.Type) bool {
	if p.notEOF() && p.peek().Type == t {
		p.currIndex++
		return true
	}
	return false
}

func (p *LangParser) expect(t l.Type) *l.Token {
	if p.isEOF() {
		panic("Early EOF! Was expecting type " + t.String())
	}
	got := p.next()
	if got.Type != t {
		got.Error("Expected type % but got %", t.String(), got.Type.String())
	}
	return got
}

func (p *LangParser) isNext(checkTypes ...l.Type) bool {
	if p.isEOF() {
		return false
	}
	pType := p.peek().Type
	for _, checkType := range checkTypes {
		if checkType == pType {
			return true
		}
	}
	return false
}

func (p *LangParser) peek() *l.Token {
	if p.isEOF() {
		panic("Early EOF! Expected more content.")
	}
	return p.Tokens[p.currIndex]
}

func (p *LangParser) next() *l.Token {
	token := p.Tokens[p.currIndex]
	p.currIndex++
	return token
}

func (p *LangParser) skip() {
	p.currIndex++
}

func (p *LangParser) notEOF() bool {
	return p.currIndex < p.tokenSize
}

func (p *LangParser) isEOF() bool {
	return p.currIndex >= p.tokenSize
}
