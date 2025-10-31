package components

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type GenericPropertySet struct {
	Component     blockly2.Expr
	ComponentType string
	Property      string
	Value         blockly2.Expr
}

func (g *GenericPropertySet) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (g *GenericPropertySet) String() string {
	return sugar.Format("set(%, %, %, %)", g.ComponentType, g.Component.String(), g.Property, g.Value.String())

}

func (g *GenericPropertySet) Blockly() blockly2.Block {
	return blockly2.Block{
		Type: "component_set_get",
		Mutation: &blockly2.Mutation{
			SetOrGet:      "set",
			PropertyName:  g.Property,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Fields: []blockly2.Field{{Name: "PROP", Value: g.Property}},
		Values: blockly2.MakeValues([]blockly2.Expr{g.Component, g.Value}, "COMPONENT", "VALUE"),
	}
}

func (g *GenericPropertySet) Continuous() bool {
	return false
}

func (g *GenericPropertySet) Consumable() bool {
	return false
}
