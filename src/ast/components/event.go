package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
	"strings"
)

type Event struct {
	IsGeneric  bool
	Component  string
	Event      string
	Parameters []string
	Body       []blockly.Expr
}

func (e *Event) String() string {
	pFormat := "when %.%(%) {\n%}"
	if e.IsGeneric {
		pFormat = "when any %.%(%) {\n%}"
	}
	return sugar.Format(pFormat, e.Component, e.Event, strings.Join(e.Parameters, ", "), blockly.PadBody(e.Body))
}

func (e *Event) Blockly() blockly.Block {
	if e.IsGeneric {
		return blockly.Block{
			// TODO: add component_type to Mutation later
			Mutation: &blockly.Mutation{
				IsGeneric: true,
				EventName: e.Event,
			},
			Statements: []blockly.Statement{blockly.CreateStatement("DO", e.Body)},
		}
	}
	return blockly.Block{
		// TODO: add component_type to Mutation later
		Mutation: &blockly.Mutation{
			IsGeneric:    false,
			InstanceName: e.Component,
			EventName:    e.Event,
		},
		Fields:     []blockly.Field{{Name: "COMPONENT_SELECTOR", Value: e.Component}},
		Statements: []blockly.Statement{blockly.CreateStatement("DO", e.Body)},
	}
}

func (e *Event) Continuous() bool {
	return false
}
