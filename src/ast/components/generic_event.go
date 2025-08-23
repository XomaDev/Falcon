package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
	"strings"
)

type GenericEvent struct {
	ComponentType string
	Event         string
	Parameters    []string
	Body          []blockly.Expr
}

func (e *GenericEvent) String() string {
	pFormat := "when any %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentType, e.Event, strings.Join(e.Parameters, ", "), blockly.PadBody(e.Body))
}

func (e *GenericEvent) Blockly() blockly.Block {
	var statements []blockly.Statement
	if len(e.Body) > 0 {
		statements = []blockly.Statement{blockly.CreateStatement("DO", e.Body)}
	}
	return blockly.Block{
		Type: "component_event",
		Mutation: &blockly.Mutation{
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
