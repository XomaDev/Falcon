package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
	"strings"
)

type Event struct {
	IsGeneric     bool
	ComponentName string
	ComponentType string
	Event         string
	Parameters    []string
	Body          []blockly.Expr
}

func (e *Event) String() string {
	pFormat := "when %.%(%) {\n%}"
	if e.IsGeneric {
		pFormat = "when any %.%(%) {\n%}"
	}
	return sugar.Format(pFormat, e.ComponentName, e.Event, strings.Join(e.Parameters, ", "), blockly.PadBody(e.Body))
}

func (e *Event) Blockly() blockly.Block {
	if e.IsGeneric {
		return blockly.Block{
			Type: "component_event",
			Mutation: &blockly.Mutation{
				IsGeneric:     true,
				EventName:     e.Event,
				ComponentType: e.ComponentType,
			},
			Statements: []blockly.Statement{blockly.CreateStatement("DO", e.Body)},
		}
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
		Statements: []blockly.Statement{blockly.CreateStatement("DO", e.Body)},
	}
}

func (e *Event) Continuous() bool {
	return false
}
