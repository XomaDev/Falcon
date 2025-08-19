package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type PropertyGet struct {
	Component string
	Property  string
}

func (p *PropertyGet) String() string {
	return sugar.Format("%.%", p.Component, p.Property)
}

func (p *PropertyGet) Blockly() blockly.Block {
	return blockly.Block{
		// TODO: add component_type to Mutation
		Mutation: &blockly.Mutation{
			SetOrGet:     "get",
			PropertyName: p.Property,
			IsGeneric:    false,
			InstanceName: p.Component,
		},
		Fields: []blockly.Field{
			{Name: "COMPONENT_SELECTOR", Value: p.Component},
			{Name: "PROP", Value: p.Property},
		},
	}
}

func (p *PropertyGet) Continuous() bool {
	return false
}
