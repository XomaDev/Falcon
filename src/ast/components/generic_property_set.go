package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type GenericPropertySet struct {
	Component     blockly.Expr
	ComponentType string
	Property      string
	Value         blockly.Expr
}

func (g *GenericPropertySet) String() string {
	return sugar.Format("set(%, %, %, %)", g.ComponentType, g.Component.String(), g.Property, g.Value.String())

}

func (g *GenericPropertySet) Blockly() blockly.Block {
	return blockly.Block{
		Type: "component_set_get",
		Mutation: &blockly.Mutation{
			SetOrGet:      "set",
			PropertyName:  g.Property,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Fields: []blockly.Field{{Name: "PROP", Value: g.Property}},
		Values: blockly.MakeValues([]blockly.Expr{g.Component, g.Value}, "COMPONENT", "VALUE"),
	}
}

func (g *GenericPropertySet) Continuous() bool {
	return false
}
