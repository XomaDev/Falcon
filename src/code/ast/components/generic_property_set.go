package components

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type GenericPropertySet struct {
	Component     ast.Expr
	ComponentType string
	Property      string
	Value         ast.Expr
}

func (g *GenericPropertySet) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (g *GenericPropertySet) String() string {
	return sugar.Format("set(%, %, %, %)", g.ComponentType, g.Component.String(), g.Property, g.Value.String())

}

func (g *GenericPropertySet) Blockly() ast.Block {
	return ast.Block{
		Type: "component_set_get",
		Mutation: &ast.Mutation{
			SetOrGet:      "set",
			PropertyName:  g.Property,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Fields: []ast.Field{{Name: "PROP", Value: g.Property}},
		Values: ast.MakeValues([]ast.Expr{g.Component, g.Value}, "COMPONENT", "VALUE"),
	}
}

func (g *GenericPropertySet) Continuous() bool {
	return false
}

func (g *GenericPropertySet) Consumable() bool {
	return false
}
