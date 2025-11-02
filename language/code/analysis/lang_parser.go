package analysis

import (
	blky "Falcon/code/ast"
	common2 "Falcon/code/ast/common"
	components2 "Falcon/code/ast/components"
	control2 "Falcon/code/ast/control"
	fundamentals2 "Falcon/code/ast/fundamentals"
	list2 "Falcon/code/ast/list"
	"Falcon/code/ast/method"
	procedures2 "Falcon/code/ast/procedures"
	variables2 "Falcon/code/ast/variables"
	lex2 "Falcon/code/lex"
	"Falcon/code/sugar"
	"strings"
)

type NameResolver struct {
	Procedures        map[string]Procedure
	ComponentTypesMap map[string]string // Button1 -> Button
	ComponentNameMap  map[string][]string
}

type Procedure struct {
	Name       string
	Parameters []string
	Returning  bool
}

type LangParser struct {
	Tokens    []*lex2.Token
	currIndex int
	tokenSize int
	Resolver  *NameResolver
}

func NewLangParser(tokens []*lex2.Token) *LangParser {
	return &LangParser{
		Tokens:    tokens,
		tokenSize: len(tokens),
		currIndex: 0,
		Resolver: &NameResolver{
			Procedures:        map[string]Procedure{},
			ComponentTypesMap: map[string]string{},
			ComponentNameMap:  map[string][]string{},
		},
	}
}

func (p *LangParser) SetComponentDefinitions(definitions map[string][]string, reverseDefinitions map[string]string) {
	p.Resolver.ComponentNameMap = definitions
	p.Resolver.ComponentTypesMap = reverseDefinitions
}

func (p *LangParser) GetComponentDefinitionsCode() string {
	// convert the AST back to syntax
	var definitions strings.Builder
	for key, value := range p.Resolver.ComponentNameMap {
		definitions.WriteString(sugar.Format("@% { % }\n", key, strings.Join(value, ", ")))
	}
	return definitions.String()
}

func (p *LangParser) ParseAll() []blky.Expr {
	var expressions []blky.Expr
	if p.notEOF() {
		p.defineStatements()
	}
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *LangParser) defineStatements() {
	for p.notEOF() && p.consume(lex2.At) {
		compType := p.name()
		p.expect(lex2.OpenCurly)
		if !p.consume(lex2.CloseCurly) {
			var componentNames []string
			for {
				name := p.name()
				componentNames = append(componentNames, name)
				p.Resolver.ComponentTypesMap[name] = compType
				if !p.consume(lex2.Comma) {
					break
				}
			}
			p.Resolver.ComponentNameMap[compType] = componentNames
			p.expect(lex2.CloseCurly)
		}
	}
}

func (p *LangParser) parse() blky.Expr {
	switch p.peek().Type {
	case lex2.If:
		return p.ifExpr()
	case lex2.For:
		return p.forExpr()
	case lex2.Each:
		return p.eachExpr()
	case lex2.While:
		return p.whileExpr()
	case lex2.Break:
		p.skip()
		return &control2.Break{}
	case lex2.WalkAll:
		p.skip()
		return &fundamentals2.WalkAll{}
	case lex2.Local:
		return p.varExpr()
	case lex2.Global:
		return p.globVar()
	case lex2.Func:
		return p.funcSmt()
	case lex2.When:
		p.skip()
		if p.consume(lex2.Any) {
			return p.genericEvent()
		}
		return p.event()
	default:
		return p.expr(0)
	}
}

func (p *LangParser) genericEvent() blky.Expr {
	componentType := p.componentType()
	p.expect(lex2.Dot)
	eventName := p.name()
	var parameters []string
	if p.isNext(lex2.OpenCurve) {
		parameters = p.parameters()
	}
	body := p.body()
	return &components2.GenericEvent{
		ComponentType: componentType,
		Event:         eventName,
		Parameters:    parameters,
		Body:          body,
	}
}

func (p *LangParser) event() blky.Expr {
	component := p.component()
	p.expect(lex2.Dot)
	eventName := p.name()
	var parameters []string
	if p.isNext(lex2.OpenCurve) {
		parameters = p.parameters()
	}
	body := p.body()
	return &components2.Event{
		ComponentName: component.Name,
		ComponentType: component.Type,
		Event:         eventName,
		Parameters:    parameters,
		Body:          body,
	}
}

func (p *LangParser) funcSmt() blky.Expr {
	p.skip()
	name := p.name()
	var parameters = p.parameters()
	returning := p.consume(lex2.Assign)
	p.Resolver.Procedures[name] = Procedure{Name: name, Parameters: parameters, Returning: returning}
	if returning {
		return &procedures2.RetProcedure{Name: name, Parameters: parameters, Result: p.parse()}
	} else {
		return &procedures2.VoidProcedure{Name: name, Parameters: parameters, Body: p.body()}
	}
}

func (p *LangParser) globVar() blky.Expr {
	p.skip()
	name := p.name()
	p.expect(lex2.Assign)
	return &variables2.Global{Name: name, Value: p.parse()}
}

func (p *LangParser) varExpr() blky.Expr {
	p.skip()

	var varNames []string
	var varValues []blky.Expr
	if p.consume(lex2.OpenCurve) {
		// a result local var
		for p.notEOF() && !p.isNext(lex2.CloseCurve) {
			name := p.name()
			p.expect(lex2.Assign)
			value := p.parse()

			varNames = append(varNames, name)
			varValues = append(varValues, value)

			if !p.consume(lex2.Comma) {
				break
			}
		}
		p.expect(lex2.CloseCurve)
		return &variables2.Var{Names: varNames, Values: varValues, Body: p.body()}
	} else {
		// a clean full scope variable
		name := p.name()
		p.expect(lex2.Assign)
		value := p.parse()
		// we have to parse rest of the body here
		return &variables2.SimpleVar{Name: name, Value: value, Body: p.bodyUntilCurly()}
	}
}

func (p *LangParser) whileExpr() *control2.While {
	p.skip()
	condition := p.expr(0)
	body := p.body()
	return &control2.While{Condition: condition, Body: body}
}

func (p *LangParser) eachExpr() blky.Expr {
	p.skip()
	keyName := p.name()
	if p.consume(lex2.DoubleColon) {
		// a dictionary pair iteration
		valueName := p.name()
		p.expect(lex2.RightArrow)
		return &control2.EachPair{KeyName: keyName, ValueName: valueName, Iterable: p.element(), Body: p.body()}
	} else {
		// a simple list iteration
		p.expect(lex2.RightArrow)
		return &control2.Each{IName: keyName, Iterable: p.element(), Body: p.body()}
	}
}

func (p *LangParser) forExpr() *control2.For {
	p.skip()
	iName := p.name()
	p.expect(lex2.Colon)
	from := p.expr(0)
	p.expect(lex2.To)
	to := p.expr(0)
	p.expect(lex2.By)
	by := p.expr(0)
	body := p.body()
	return &control2.For{
		IName: iName,
		From:  from,
		To:    to,
		By:    by,
		Body:  body,
	}
}

func (p *LangParser) ifExpr() blky.Expr {
	p.skip()
	if p.isNext(lex2.OpenCurve) {
		return p.simpleIf()
	}
	var conditions []blky.Expr
	var bodies [][]blky.Expr

	conditions = append(conditions, p.expr(0))
	bodies = append(bodies, p.body())

	for p.notEOF() && p.consume(lex2.Elif) {
		conditions = append(conditions, p.expr(0))
		bodies = append(bodies, p.body())
	}
	var elseBody []blky.Expr
	if p.notEOF() && p.consume(lex2.Else) {
		elseBody = p.body()
	}
	return &control2.If{Conditions: conditions, Bodies: bodies, ElseBody: elseBody}
}

func (p *LangParser) simpleIf() *control2.SimpleIf {
	p.expect(lex2.OpenCurve)
	condition := p.parse()
	p.expect(lex2.CloseCurve)
	then := p.parse()
	p.expect(lex2.Else)
	elze := p.parse()
	return &control2.SimpleIf{Condition: condition, Then: then, Else: elze}
}

func (p *LangParser) body() []blky.Expr {
	p.expect(lex2.OpenCurly)
	expressions := p.bodyUntilCurly()
	p.expect(lex2.CloseCurly)
	return expressions
}

func (p *LangParser) bodyUntilCurly() []blky.Expr {
	var expressions []blky.Expr
	if p.isNext(lex2.CloseCurly) {
		return expressions
	}
	for p.notEOF() && !p.isNext(lex2.CloseCurly) {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *LangParser) expr(minPrecedence int) blky.Expr {
	left := p.element()
	for p.notEOF() {
		opToken := p.peek()
		if !opToken.HasFlag(lex2.Operator) {
			break
		}
		precedence := precedenceOf(opToken.Flags[0])
		if precedence == -1 || precedence < minPrecedence {
			break
		}
		p.skip()
		var right blky.Expr
		if opToken.HasFlag(lex2.PreserveOrder) {
			right = p.element()
		} else {
			right = p.expr(precedence)
		}
		if rBinExpr, ok := right.(*common2.BinaryExpr); ok && rBinExpr.CanRepeat(opToken.Type) {
			// for NoPreserveOrder: merge binary expr with same operator (towards right)
			rBinExpr.Operands = append([]blky.Expr{left}, rBinExpr.Operands...)
			left = rBinExpr
		} else if lBinExpr, ok := left.(*common2.BinaryExpr); ok && lBinExpr.CanRepeat(opToken.Type) {
			// for PreserveOder: merge binary expr with same operator (towards left)
			lBinExpr.Operands = append(lBinExpr.Operands, right)
		} else {
			// a new binary node
			left = &common2.BinaryExpr{Where: opToken, Operands: []blky.Expr{left, right}, Operator: opToken.Type}
		}
	}
	return left
}

func precedenceOf(flag lex2.Flag) int {
	switch flag {
	case lex2.AssignmentType:
		return 0
	case lex2.Pair:
		return 1
	case lex2.TextJoin:
		return 2
	case lex2.LLogicOr:
		return 3
	case lex2.LLogicAnd:
		return 4
	case lex2.BBitwiseOr:
		return 5
	case lex2.BBitwiseAnd:
		return 6
	case lex2.BBitwiseXor:
		return 7
	case lex2.Equality:
		return 8
	case lex2.Relational:
		return 9
	case lex2.Binary:
		return 10
	case lex2.BinaryL1:
		return 11
	case lex2.BinaryL2:
		return 12
	default:
		return -1
	}
}

func (p *LangParser) element() blky.Expr {
	left := p.term()
	for p.notEOF() {
		pe := p.peek()
		// check if it's a variable Get, if so, check if it refers to a component
		if getExpr, ok := left.(*fundamentals2.Component); ok && pe.Type == lex2.Dot {
			if compType, exists := p.Resolver.ComponentTypesMap[getExpr.Name]; exists {
				// a specific component call (MethodCall, PropertyGet, PropertySet)
				left = p.componentCall(getExpr.Name, compType)
				continue
			}
		}

		switch pe.Type {
		case lex2.At:
			left = p.helperDropdown(left)
		case lex2.Dot:
			left = p.objectCall(left)
			continue
		//case l.RightArrow:
		//left = &common.Convert{Where: p.next(), On: left, BlocklyName: p.name()}
		//continue
		case lex2.Question:
			left = &common2.Question{Where: p.next(), On: left, Question: p.name()}
			continue
		case lex2.DoubleColon:
			// constant value transformer
			left = &common2.Transform{Where: p.next(), On: left, Name: p.name()}
		case lex2.OpenSquare:
			p.skip()
			// an index element access
			left = &list2.Get{List: left, Index: p.parse()}
			p.expect(lex2.CloseSquare)
			continue
		}
		break
	}
	return left
}

func (p *LangParser) componentCall(compName string, compType string) blky.Expr {
	p.expect(lex2.Dot)
	resource := p.name()
	if p.isNext(lex2.OpenCurve) {
		return &components2.MethodCall{
			ComponentName: compName,
			ComponentType: compType,
			Method:        resource,
			Args:          p.arguments(),
		}
	} else if p.consume(lex2.Assign) {
		assignment := p.expr(0)
		return &components2.PropertySet{
			ComponentName: compName,
			ComponentType: compType,
			Property:      resource,
			Value:         assignment,
		}
	}
	return &components2.PropertyGet{ComponentName: compName, ComponentType: compType, Property: resource}
}

func (p *LangParser) helperDropdown(keyExpr blky.Expr) blky.Expr {
	where := p.next()
	if key, ok := keyExpr.(*variables2.Get); ok {
		return &fundamentals2.HelperDropdown{Key: key.Name, Option: p.name()}
	}
	where.Error("Invalid Helper Access operation ")
	panic("")
}

func (p *LangParser) objectCall(object blky.Expr) blky.Expr {
	p.skip()
	where := p.next()
	name := *where.Content

	var args []blky.Expr
	if p.isNext(lex2.OpenCurve) {
		args = p.arguments()
		if !p.isNext(lex2.OpenCurly) {
			// he's a simple call!
			return &method.Call{Where: where, On: object, Name: name, Args: args}
		}
	}
	p.expect(lex2.OpenCurly)
	// oh, no! he's a transformer >_>
	var namesUsed []string
	if !p.consume(lex2.RightArrow) {
		for {
			namesUsed = append(namesUsed, p.name())
			if !p.consume(lex2.Comma) {
				break
			}
		}
		p.consume(lex2.RightArrow)
	}
	transformer := p.parse()
	p.consume(lex2.CloseCurly)
	return &list2.Transformer{
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
	case lex2.Undefined:
		return &common2.EmptySocket{}
	case lex2.OpenSquare:
		return p.list()
	case lex2.OpenCurly:
		return p.dictionary()
	case lex2.OpenCurve:
		e := p.parse()
		p.expect(lex2.CloseCurve)
		return e
	case lex2.Not:
		return &fundamentals2.Not{Expr: p.element()}
	case lex2.Do:
		return p.doExpr()
	case lex2.If:
		return p.simpleIf()
	case lex2.Compute:
		return p.computeExpr()
	default:
		if token.HasFlag(lex2.Value) {
			value := p.value(token)
			if nameExpr, ok := value.(*variables2.Get); ok && p.notEOF() && p.isNext(lex2.OpenCurve) {
				signature, ok := p.Resolver.Procedures[nameExpr.Name]
				if ok {
					return &procedures2.Call{
						Name:       nameExpr.Name,
						Parameters: signature.Parameters,
						Arguments:  p.arguments(),
						Returning:  signature.Returning}
				} else {
					return &common2.FuncCall{Where: nameExpr.Where, Name: nameExpr.Name, Args: p.arguments()}
				}
			}
			return value
		}
		token.Error("Unexpected! %", token.String())
		panic("") // unreachable
	}
	// TODO: a returning local statement might be possible here
}

func (p *LangParser) computeExpr() *variables2.VarResult {
	var varNames []string
	var varValues []blky.Expr
	p.expect(lex2.OpenCurve)

	// a result local var
	for p.notEOF() && !p.isNext(lex2.CloseCurve) {
		name := p.name()
		p.expect(lex2.Assign)
		value := p.parse()

		varNames = append(varNames, name)
		varValues = append(varValues, value)

		if !p.consume(lex2.Comma) {
			break
		}
	}
	p.expect(lex2.CloseCurve)
	p.expect(lex2.RightArrow)
	return &variables2.VarResult{Names: varNames, Values: varValues, Result: p.parse()}
}

func (p *LangParser) doExpr() *control2.Do {
	body := p.body()
	p.expect(lex2.RightArrow)
	result := p.expr(0)
	return &control2.Do{Body: body, Result: result}
}

func (p *LangParser) dictionary() *fundamentals2.Dictionary {
	var elements []blky.Expr
	if !p.consume(lex2.CloseCurly) {
		for p.notEOF() {
			elements = append(elements, p.expr(0))
			if !p.consume(lex2.Comma) {
				break
			}
		}
		p.expect(lex2.CloseCurly)
	}
	return &fundamentals2.Dictionary{Elements: elements}
}

func (p *LangParser) list() *fundamentals2.List {
	var elements []blky.Expr
	if !p.consume(lex2.CloseSquare) {
		for p.notEOF() {
			elements = append(elements, p.expr(0))
			if !p.consume(lex2.Comma) {
				break
			}
		}
		p.expect(lex2.CloseSquare)
	}
	return &fundamentals2.List{Elements: elements}
}

func (p *LangParser) parameters() []string {
	p.expect(lex2.OpenCurve)
	var parameters []string
	if !p.consume(lex2.CloseCurve) {
		for p.notEOF() && !p.isNext(lex2.CloseCurve) {
			parameters = append(parameters, p.name())
			if !p.consume(lex2.Comma) {
				break
			}
		}
		p.expect(lex2.CloseCurve)
	}
	return parameters
}

func (p *LangParser) arguments() []blky.Expr {
	p.expect(lex2.OpenCurve)
	var args []blky.Expr
	if p.consume(lex2.CloseCurve) {
		return args
	}
	for p.notEOF() {
		args = append(args, p.expr(0))
		if !p.consume(lex2.Comma) {
			break
		}
	}
	p.expect(lex2.CloseCurve)
	return args
}

func (p *LangParser) value(t *lex2.Token) blky.Expr {
	switch t.Type {
	case lex2.True, lex2.False:
		return &fundamentals2.Boolean{Value: t.Type == lex2.True}
	case lex2.Number:
		return &fundamentals2.Number{Content: *t.Content}
	case lex2.Text:
		return &fundamentals2.Text{Content: *t.Content}
	case lex2.Name:
		if compType, exists := p.Resolver.ComponentTypesMap[*t.Content]; exists {
			return &fundamentals2.Component{Name: *t.Content, Type: compType}
		}
		return &variables2.Get{Where: t, Global: false, Name: *t.Content}
	case lex2.This:
		p.expect(lex2.Dot)
		return &variables2.Get{Where: t, Global: true, Name: p.name()}
	case lex2.Color:
		p.expect(lex2.Colon)
		return &fundamentals2.Color{Where: t, Name: p.name()}
	default:
		t.Error("Unknown value type '%'", t.String())
		panic("") // unreachable
	}
}

func (p *LangParser) componentType() string {
	token := p.expect(lex2.Name)
	name := *token.Content
	if _, exists := p.Resolver.ComponentNameMap[name]; exists {
		return name
	}
	token.Error("Undefined component group %", name)
	panic("")
}

func (p *LangParser) component() fundamentals2.Component {
	token := p.expect(lex2.Name)
	name := *token.Content
	if compType, exists := p.Resolver.ComponentTypesMap[name]; exists {
		return fundamentals2.Component{Name: name, Type: compType}
	}
	token.Error("Undefined component %", name)
	panic("")
}

func (p *LangParser) name() string {
	return *p.expect(lex2.Name).Content
}

func (p *LangParser) consume(t lex2.Type) bool {
	if p.notEOF() && p.peek().Type == t {
		p.currIndex++
		return true
	}
	return false
}

func (p *LangParser) expect(t lex2.Type) *lex2.Token {
	if p.isEOF() {
		panic("Early EOF! Was expecting type " + t.String())
	}
	got := p.next()
	if got.Type != t {
		got.Error("Expected type % but got %", t.String(), got.String())
	}
	return got
}

func (p *LangParser) isNext(checkTypes ...lex2.Type) bool {
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

func (p *LangParser) peek() *lex2.Token {
	if p.isEOF() {
		panic("Early EOF! Expected more content.")
	}
	return p.Tokens[p.currIndex]
}

func (p *LangParser) next() *lex2.Token {
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
