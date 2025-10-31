package components

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
	"strings"
)

type Event struct {
	ComponentName string
	ComponentType string
	Event         string
	Parameters    []string
	Body          []ast.Expr
}

func (e *Event) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (e *Event) String() string {
	pFormat := "when %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentName, e.Event, strings.Join(e.Parameters, ", "), ast.PadBody(e.Body))
}

func (e *Event) Blockly() ast.Block {
	var statements []ast.Statement
	if len(e.Body) > 0 {
		statements = []ast.Statement{ast.CreateStatement("DO", e.Body)}
	}
	return ast.Block{
		Type: "component_event",
		Mutation: &ast.Mutation{
			IsGeneric:     false,
			InstanceName:  e.ComponentName,
			EventName:     e.Event,
			ComponentType: e.ComponentType,
		},
		Fields:     []ast.Field{{Name: "COMPONENT_SELECTOR", Value: e.ComponentName}},
		Statements: statements,
	}
}

func (e *Event) Continuous() bool {
	return false
}

func (e *Event) Consumable() bool {
	return false
}
