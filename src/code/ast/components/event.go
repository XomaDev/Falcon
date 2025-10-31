package components

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
	"strings"
)

type Event struct {
	ComponentName string
	ComponentType string
	Event         string
	Parameters    []string
	Body          []blockly2.Expr
}

func (e *Event) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (e *Event) String() string {
	pFormat := "when %.%(%) {\n%}"
	return sugar.Format(pFormat, e.ComponentName, e.Event, strings.Join(e.Parameters, ", "), blockly2.PadBody(e.Body))
}

func (e *Event) Blockly() blockly2.Block {
	var statements []blockly2.Statement
	if len(e.Body) > 0 {
		statements = []blockly2.Statement{blockly2.CreateStatement("DO", e.Body)}
	}
	return blockly2.Block{
		Type: "component_event",
		Mutation: &blockly2.Mutation{
			IsGeneric:     false,
			InstanceName:  e.ComponentName,
			EventName:     e.Event,
			ComponentType: e.ComponentType,
		},
		Fields:     []blockly2.Field{{Name: "COMPONENT_SELECTOR", Value: e.ComponentName}},
		Statements: statements,
	}
}

func (e *Event) Continuous() bool {
	return false
}

func (e *Event) Consumable() bool {
	return false
}
