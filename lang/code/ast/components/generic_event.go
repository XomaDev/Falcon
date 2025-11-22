package components

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
	"strings"
)

type GenericEvent struct {
	ComponentType string
	Event         string
	Parameters    []string
	Body          []ast.Expr
}

func (g *GenericEvent) String() string {
	pFormat := "when any %.%(%) {\n%}"
	return sugar.Format(pFormat, g.ComponentType, g.Event, strings.Join(g.Parameters, ", "), ast.PadBody(g.Body))
}

func (g *GenericEvent) Blockly(flags ...bool) ast.Block {
	var statements []ast.Statement
	if len(g.Body) > 0 {
		statements = []ast.Statement{ast.CreateStatement("DO", g.Body)}
	}
	return ast.Block{
		Type: "component_event",
		Mutation: &ast.Mutation{
			IsGeneric:     true,
			EventName:     g.Event,
			ComponentType: g.ComponentType,
		},
		Statements: statements,
	}
}

func (g *GenericEvent) Continuous() bool {
	return false
}

func (g *GenericEvent) Consumable(flags ...bool) bool {
	return false
}

func (g *GenericEvent) Signature() []ast.Signature {
	return []ast.Signature{ast.SignVoid}
}
