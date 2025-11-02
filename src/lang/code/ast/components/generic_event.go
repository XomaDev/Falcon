package components

import (
	"Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type GenericEvent struct {
	ComponentType string
	Event         string
	Parameters    []string
	Body          []ast.Expr
}

func (e *GenericEvent) Yail() string {
	yail := "(define-generic-event "
	yail += e.ComponentType
	yail += " "
	yail += e.Event
	yail += " ("
	for _, p := range e.Parameters {
		yail += "$" + p + " "
	}
	yail += ") (set-this-form) "
	yail += ast.PadBodyYail(e.Body)
	yail += ")"
	return yail
}

func (e *GenericEvent) String() string {
	pFormat := "when any %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentType, e.Event, strings.Join(e.Parameters, ", "), ast.PadBody(e.Body))
}

func (e *GenericEvent) Blockly() ast.Block {
	var statements []ast.Statement
	if len(e.Body) > 0 {
		statements = []ast.Statement{ast.CreateStatement("DO", e.Body)}
	}
	return ast.Block{
		Type: "component_event",
		Mutation: &ast.Mutation{
			IsGeneric:     true,
			EventName:     e.Event,
			ComponentType: e.ComponentType,
		},
		Statements: statements,
	}
}

func (e *GenericEvent) Continuous() bool {
	return false
}

func (e *GenericEvent) Consumable() bool {
	return false
}
