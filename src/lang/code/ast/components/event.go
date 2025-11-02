package components

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type Event struct {
	ComponentName string
	ComponentType string
	Event         string
	Parameters    []string
	Body          []ast2.Expr
}

func (e *Event) Yail() string {
	yail := "(define-event "
	yail += e.ComponentName
	yail += " "
	yail += e.Event
	yail += " ("
	for _, p := range e.Parameters {
		yail += "$" + p + " "
	}
	yail += ") (set-this-form) "
	yail += ast2.PadBodyYail(e.Body)
	yail += ")"
	return yail
}

func (e *Event) String() string {
	pFormat := "when %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentName, e.Event, strings.Join(e.Parameters, ", "), ast2.PadBody(e.Body))
}

func (e *Event) Blockly() ast2.Block {
	var statements []ast2.Statement
	if len(e.Body) > 0 {
		statements = []ast2.Statement{ast2.CreateStatement("DO", e.Body)}
	}
	return ast2.Block{
		Type: "component_event",
		Mutation: &ast2.Mutation{
			IsGeneric:     false,
			InstanceName:  e.ComponentName,
			EventName:     e.Event,
			ComponentType: e.ComponentType,
		},
		Fields:     []ast2.Field{{Name: "COMPONENT_SELECTOR", Value: e.ComponentName}},
		Statements: statements,
	}
}

func (e *Event) Continuous() bool {
	return false
}

func (e *Event) Consumable() bool {
	return false
}
