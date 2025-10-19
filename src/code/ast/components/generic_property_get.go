package components

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type GenericPropertyGet struct {
	Component     blockly.Expr
	ComponentType string
	Property      string
}

func (g *GenericPropertyGet) String() string {
	return sugar.Format("get(%, %, %)", g.ComponentType, g.Component.String(), g.Property)
}

func (g *GenericPropertyGet) Blockly() blockly.Block {
	return blockly.Block{
		Type: "component_set_get",
		Mutation: &blockly.Mutation{
			SetOrGet:      "get",
			PropertyName:  g.Property,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Fields: []blockly.Field{{Name: "PROP", Value: g.Property}},
		Values: []blockly.Value{{Name: "COMPONENT", Block: g.Component.Blockly()}},
	}
}

func (g *GenericPropertyGet) Continuous() bool {
	return false
}

func (g *GenericPropertyGet) Consumable() bool {
	return true
}
