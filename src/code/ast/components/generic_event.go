package components

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
	"strings"
)

type GenericEvent struct {
	ComponentType string
	Event         string
	Parameters    []string
	Body          []blockly2.Expr
}

func (e *GenericEvent) String() string {
	pFormat := "when any %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentType, e.Event, strings.Join(e.Parameters, ", "), blockly2.PadBody(e.Body))
}

func (e *GenericEvent) Blockly() blockly2.Block {
	var statements []blockly2.Statement
	if len(e.Body) > 0 {
		statements = []blockly2.Statement{blockly2.CreateStatement("DO", e.Body)}
	}
	return blockly2.Block{
		Type: "component_event",
		Mutation: &blockly2.Mutation{
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
