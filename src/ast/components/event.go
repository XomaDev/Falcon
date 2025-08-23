package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
	"strings"
)

type Event struct {
	ComponentName string
	ComponentType string
	Event         string
	Parameters    []string
	Body          []blockly.Expr
}

func (e *Event) String() string {
	pFormat := "when %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentName, e.Event, strings.Join(e.Parameters, ", "), blockly.PadBody(e.Body))
}

func (e *Event) Blockly() blockly.Block {
	var statements []blockly.Statement
	if len(e.Body) > 0 {
		statements = []blockly.Statement{blockly.CreateStatement("DO", e.Body)}
	}
	return blockly.Block{
		Type: "component_event",
		Mutation: &blockly.Mutation{
			IsGeneric:     false,
			InstanceName:  e.ComponentName,
			EventName:     e.Event,
			ComponentType: e.ComponentType,
		},
		Fields:     []blockly.Field{{Name: "COMPONENT_SELECTOR", Value: e.ComponentName}},
		Statements: statements,
	}
}

func (e *Event) Continuous() bool {
	return false
}

func (e *Event) Consumable() bool {
	return false
}
