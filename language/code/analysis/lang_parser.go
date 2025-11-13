package analysis

import (
	"Falcon/code/ast"
	"Falcon/code/ast/common"
	"Falcon/code/ast/components"
	"Falcon/code/ast/control"
	"Falcon/code/ast/fundamentals"
	"Falcon/code/ast/list"
	"Falcon/code/ast/method"
	"Falcon/code/ast/procedures"
	"Falcon/code/ast/variables"
	"Falcon/code/sugar"
	"strings"
)
import l "Falcon/code/lex"

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
	Tokens      []*l.Token
	currIndex   int
	tokenSize   int
	Resolver    *NameResolver
	ScopeCursor *ScopeCursor
}

func NewLangParser(tokens []*l.Token) *LangParser {
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

func (p *LangParser) ParseAll() []ast.Expr {
	var expressions []ast.Expr
	if p.notEOF() {
		p.defineStatements()
	}
	for p.notEOF() {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *LangParser) defineStatements() {
	for p.notEOF() && p.consume(l.At) {
		compType := p.name()
		p.expect(l.OpenCurly)
		if !p.consume(l.CloseCurly) {
			var componentNames []string
			for {
				name := p.name()
				componentNames = append(componentNames, name)
				p.Resolver.ComponentTypesMap[name] = compType
				if !p.consume(l.Comma) {
					break
				}
			}
			p.Resolver.ComponentNameMap[compType] = componentNames
			p.expect(l.CloseCurly)
		}
	}
}

func (p *LangParser) parse() ast.Expr {
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
	case l.When:
		p.skip()
		if p.consume(l.Any) {
			return p.genericEvent()
		}
		return p.event()
	default:
		return p.expr(0)
	}
}

func (p *LangParser) genericEvent() ast.Expr {
	componentType := p.componentType()
	p.expect(l.Dot)
	eventName := p.name()
	var parameters []string
	if p.isNext(l.OpenCurve) {
		parameters = p.parameters()
	}
	body := p.body(ScopeGenericEvent)
	return &components.GenericEvent{
		ComponentType: componentType,
		Event:         eventName,
		Parameters:    parameters,
		Body:          body,
	}
}

func (p *LangParser) event() ast.Expr {
	component := p.component()
	p.expect(l.Dot)
	eventName := p.name()
	var parameters []string
	if p.isNext(l.OpenCurve) {
		parameters = p.parameters()
	}
	body := p.body(ScopeEvent)
	return &components.Event{
		ComponentName: component.Name,
		ComponentType: component.Type,
		Event:         eventName,
		Parameters:    parameters,
		Body:          body,
	}
}

func (p *LangParser) funcSmt() ast.Expr {
	p.skip()
	name := p.name()
	var parameters = p.parameters()
	returning := p.consume(l.Assign)
	p.Resolver.Procedures[name] = Procedure{Name: name, Parameters: parameters, Returning: returning}
	if returning {
		return &procedures.RetProcedure{Name: name, Parameters: parameters, Result: p.parse()}
	} else {
		return &procedures.VoidProcedure{Name: name, Parameters: parameters, Body: p.body(ScopeProc)}
	}
}

func (p *LangParser) globVar() ast.Expr {
	where := p.next()
	if !p.ScopeCursor.AtRoot() {
		where.Error("Global variables can only be defined at the root.")
	}
	name := p.name()
	p.expect(l.Assign)
	return &variables.Global{Name: name, Value: p.parse()}
}

func (p *LangParser) varExpr() ast.Expr {
	where := p.next()
	if p.ScopeCursor.AtRoot() {
		where.Error("Cannot define a local variable at the root level.")
	}
	// a clean full scope variable
	name := p.name()
	p.expect(l.Assign)
	value := p.parse()
	// we have to parse rest of the body here
	return &variables.SimpleVar{Name: name, Value: value, Body: p.bodyUntilCurly()}
}

func (p *LangParser) whileExpr() *control.While {
	p.skip()
	condition := p.expr(0)
	body := p.body(Loop)
	return &control.While{Condition: condition, Body: body}
}

func (p *LangParser) eachExpr() ast.Expr {
	p.skip()
	if p.consume(l.OpenCurve) {
		// a dictionary pair iteration
		keyName := p.name()
		p.expect(l.Comma)
		valueName := p.name()
		p.expect(l.CloseCurve)
		p.expect(l.In)
		return &control.EachPair{KeyName: keyName, ValueName: valueName, Iterable: p.element(), Body: p.body(Loop)}
	} else {
		keyName := p.name()
		p.expect(l.In)
		return &control.Each{IName: keyName, Iterable: p.element(), Body: p.body(Loop)}
	}
}

func (p *LangParser) forExpr() *control.For {
	p.skip()
	iName := p.name()
	p.expect(l.Colon)
	from := p.expr(0)
	p.expect(l.To)
	to := p.expr(0)
	var by ast.Expr
	if p.consume(l.By) {
		by = p.expr(0)
	} else {
		by = &fundamentals.Number{Content: "1"}
	}
	body := p.body(Loop)
	return &control.For{
		IName: iName,
		From:  from,
		To:    to,
		By:    by,
		Body:  body,
	}
}

func (p *LangParser) ifExpr() ast.Expr {
	p.skip()
	if p.isNext(l.OpenCurve) {
		return p.simpleIf()
	}
	var conditions []ast.Expr
	var bodies [][]ast.Expr

	conditions = append(conditions, p.expr(0))
	bodies = append(bodies, p.body(IfBody))

	for p.notEOF() && p.consume(l.Elif) {
		conditions = append(conditions, p.expr(0))
		bodies = append(bodies, p.body(IfBody))
	}
	var elseBody []ast.Expr
	if p.notEOF() && p.consume(l.Else) {
		elseBody = p.body(IfBody)
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

func (p *LangParser) body(scope Scope) []ast.Expr {
	where := p.expect(l.OpenCurly)
	p.ScopeCursor.Enter(where, scope)
	expressions := p.bodyUntilCurly()
	p.ScopeCursor.Exit(scope)
	p.expect(l.CloseCurly)
	return expressions
}

func (p *LangParser) bodyUntilCurly() []ast.Expr {
	var expressions []ast.Expr
	if p.isNext(l.CloseCurly) {
		return expressions
	}
	for p.notEOF() && !p.isNext(l.CloseCurly) {
		expressions = append(expressions, p.parse())
	}
	return expressions
}

func (p *LangParser) expr(minPrecedence int) ast.Expr {
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
		var right ast.Expr
		if opToken.HasFlag(l.PreserveOrder) {
			right = p.element()
		} else {
			right = p.expr(precedence)
		}
		if rBinExpr, ok := right.(*common.BinaryExpr); ok && rBinExpr.CanRepeat(opToken.Type) {
			// for NoPreserveOrder: merge binary expr with same operator (towards right)
			rBinExpr.Operands = append([]ast.Expr{left}, rBinExpr.Operands...)
			left = rBinExpr
		} else if lBinExpr, ok := left.(*common.BinaryExpr); ok && lBinExpr.CanRepeat(opToken.Type) {
			// for PreserveOder: merge binary expr with same operator (towards left)
			lBinExpr.Operands = append(lBinExpr.Operands, right)
		} else {
			// a new binary node
			left = &common.BinaryExpr{Where: opToken, Operands: []ast.Expr{left, right}, Operator: opToken.Type}
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

func (p *LangParser) element() ast.Expr {
	left := p.term()
	for p.notEOF() {
		pe := p.peek()
		// check if it's a variable Get, if so, check if it refers to a component
		if getExpr, ok := left.(*fundamentals.Component); ok && pe.Type == l.Dot {
			if compType, exists := p.Resolver.ComponentTypesMap[getExpr.Name]; exists {
				// a specific component call (MethodCall, PropertyGet, PropertySet)
				left = p.componentCall(getExpr.Name, compType)
				continue
			}
		}

		switch pe.Type {
		case l.At:
			left = p.helperDropdown(left)
		case l.Dot:
			left = p.objectCall(left)
			continue
		//case l.RightArrow:
		//left = &common.Convert{Where: p.next(), On: left, Name: p.name()}
		//continue
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

func (p *LangParser) componentCall(compName string, compType string) ast.Expr {
	p.expect(l.Dot)
	resource := p.name()
	if p.isNext(l.OpenCurve) {
		return &components.MethodCall{
			ComponentName: compName,
			ComponentType: compType,
			Method:        resource,
			Args:          p.arguments(),
		}
	} else if p.consume(l.Assign) {
		assignment := p.expr(0)
		return &components.PropertySet{
			ComponentName: compName,
			ComponentType: compType,
			Property:      resource,
			Value:         assignment,
		}
	}
	return &components.PropertyGet{ComponentName: compName, ComponentType: compType, Property: resource}
}

func (p *LangParser) helperDropdown(keyExpr ast.Expr) ast.Expr {
	where := p.next()
	if key, ok := keyExpr.(*variables.Get); ok {
		return &fundamentals.HelperDropdown{Key: key.Name, Option: p.name()}
	}
	where.Error("Invalid Helper Access operation ")
	panic("")
}

func (p *LangParser) objectCall(object ast.Expr) ast.Expr {
	p.skip()
	where := p.next()
	name := *where.Content

	var args []ast.Expr
	if p.isNext(l.OpenCurve) {
		args = p.arguments()
		if !p.isNext(l.OpenCurly) {
			// he's a simple call!
			return &method.Call{Where: where, On: object, Name: name, Args: args}
		}
	}
	p.expect(l.OpenCurly)
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

func (p *LangParser) term() ast.Expr {
	token := p.next()
	switch token.Type {
	case l.Undefined:
		return &common.EmptySocket{}
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
	//case l.Do:
	//	return p.doExpr()
	case l.If:
		return p.simpleIf()
	case l.Compute:
		return p.computeExpr()
	default:
		if token.HasFlag(l.Value) {
			value := p.value(token)
			if nameExpr, ok := value.(*variables.Get); ok && p.notEOF() && p.isNext(l.OpenCurve) {
				signature, ok := p.Resolver.Procedures[nameExpr.Name]
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
	// TODO: a returning local statement might be possible here
}

func (p *LangParser) computeExpr() *variables.VarResult {
	var varNames []string
	var varValues []ast.Expr
	p.expect(l.OpenCurve)

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
	p.expect(l.RightArrow)
	return &variables.VarResult{Names: varNames, Values: varValues, Result: p.parse()}
}

//func (p *LangParser) doExpr() *control.Do {
//	body := p.body()
//	p.expect(l.RightArrow)
//	result := p.expr(0)
//	return &control.Do{Body: body, Result: result}
//}

func (p *LangParser) dictionary() *fundamentals.Dictionary {
	var elements []ast.Expr
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
	var elements []ast.Expr
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

func (p *LangParser) parameters() []string {
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
	return parameters
}

func (p *LangParser) arguments() []ast.Expr {
	p.expect(l.OpenCurve)
	var args []ast.Expr
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

func (p *LangParser) value(t *l.Token) ast.Expr {
	switch t.Type {
	case l.True, l.False:
		return &fundamentals.Boolean{Value: t.Type == l.True}
	case l.Number:
		return &fundamentals.Number{Content: *t.Content}
	case l.Text:
		return &fundamentals.Text{Content: *t.Content}
	case l.Name:
		if compType, exists := p.Resolver.ComponentTypesMap[*t.Content]; exists {
			return &fundamentals.Component{Name: *t.Content, Type: compType}
		}
		return &variables.Get{Where: t, Global: false, Name: *t.Content}
	case l.This:
		p.expect(l.Dot)
		return &variables.Get{Where: t, Global: true, Name: p.name()}
	case l.ColorCode:
		return &fundamentals.Color{Where: t, Hex: *t.Content}
	default:
		t.Error("Unknown value type '%'", t.String())
		panic("") // unreachable
	}
}

func (p *LangParser) componentType() string {
	token := p.expect(l.Name)
	name := *token.Content
	if _, exists := p.Resolver.ComponentNameMap[name]; exists {
		return name
	}
	token.Error("Undefined component group %", name)
	panic("")
}

func (p *LangParser) component() fundamentals.Component {
	token := p.expect(l.Name)
	name := *token.Content
	if compType, exists := p.Resolver.ComponentTypesMap[name]; exists {
		return fundamentals.Component{Name: name, Type: compType}
	}
	token.Error("Undefined component %", name)
	panic("")
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
		got.Error("Expected type % but got %", t.String(), got.String())
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
