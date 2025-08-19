package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type PropertySet struct {
	Component string
	Property  string
	Value     blockly.Expr
}

func (p *PropertySet) String() string {
	return sugar.Format("%.% = %", p.Component, p.Property, p.Value.String())
}

func (p *PropertySet) Blockly() blockly.Block {
	return blockly.Block{
		// TODO: add component_type to Mutation
		Mutation: &blockly.Mutation{
			SetOrGet:     "set",
			PropertyName: p.Property,
			IsGeneric:    false,
			InstanceName: p.Component,
		},
		Fields: blockly.FieldsFromMap(map[string]string{
			"COMPONENT_SELECTOR": p.Component,
			"PROP":               p.Property,
		}),
		Values: []blockly.Value{{Name: "VALUE", Block: p.Value.Blockly()}},
	}
}

func (p *PropertySet) Continuous() bool {
	return false
}
