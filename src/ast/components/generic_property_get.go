package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type GenericPropertyGet struct {
	Component blockly.Expr
	Property  string
}

func (g *GenericPropertyGet) String() string {
	pFormat := "%->%"
	if !g.Component.Continuous() {
		pFormat = "(%)->%"
	}
	return sugar.Format(pFormat, g.Component.String(), g.Property)
}

func (g *GenericPropertyGet) Blockly() blockly.Block {
	return blockly.Block{
		// TODO: add component_type to Mutation
		Mutation: &blockly.Mutation{
			SetOrGet:     "get",
			PropertyName: g.Property,
			IsGeneric:    true,
		},
		Fields: []blockly.Field{{Name: "PROP", Value: g.Property}},
		Values: []blockly.Value{{Name: "COMPONENT", Block: g.Component.Blockly()}},
	}
}

func (g *GenericPropertyGet) Continuous() bool {
	return false
}
