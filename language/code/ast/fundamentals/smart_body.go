package fundamentals

import (
	"Falcon/code/ast"
	"Falcon/code/ast/variables"
)

type SmartBody struct {
	Body []ast.Expr
}

func (s *SmartBody) Yail() string {
	panic("implement me")
}

func (s *SmartBody) String() string {
	return ast.PadBody(s.Body)
}

func (s *SmartBody) Blockly() ast.Block {
	// a single expression, just inline it
	if len(s.Body) == 1 {
		return s.Body[0].Blockly()
	}
	// prepare a do expression out of the Then
	doResult := s.Body[len(s.Body)-1]
	doBody := s.Body[:len(s.Body)-1]
	doExpr := ast.Block{
		Type:       "controls_do_then_return",
		Statements: []ast.Statement{ast.CreateStatement("STM", doBody)},
		Values:     []ast.Value{{Name: "VALUE", Block: doResult.Blockly()}},
	}

	var namesLocal = s.mutateVars()
	if len(namesLocal) == 0 {
		// no variables declared in the Then, a do expression is enough
		return doExpr
	}
	// We'd need to use a local result expression
	var defaultLocalVals []ast.Expr
	for k := range defaultLocalVals {
		defaultLocalVals[k] = &Boolean{Value: false}
	}
	return ast.Block{
		Type:     "local_declaration_expression",
		Mutation: &ast.Mutation{LocalNames: ast.MakeLocalNames(namesLocal...)},
		Fields:   ast.ToFields("VAR", namesLocal),
		Values: append(ast.ValuesByPrefix("DECL", defaultLocalVals),
			ast.Value{Name: "RETURN", Block: doExpr}),
	}
}

// mutateVars returns a name list of declared variables, and the declarations are mutated to a set call.
// The variables will later be defined at the top.
func (s *SmartBody) mutateVars() []string {
	var names []string
	for k, expr := range s.Body {
		// We only have simple variables
		if e, ok := expr.(*variables.SimpleVar); ok {
			names = append(names, e.Name)
			// Mutate it to a set function
			s.Body[k] = &variables.Set{Global: false, Name: e.Name, Expr: e.Value}
		}
	}
	return names
}

func (s *SmartBody) Continuous() bool {
	return false
}

func (s *SmartBody) Consumable() bool {
	return true
}

func (s *SmartBody) Signature() []ast.Signature {
	return s.Body[len(s.Body)-1].Signature()
}
