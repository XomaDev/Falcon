package components

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type GenericPropertyGet struct {
	Component     ast.Expr
	ComponentType string
	Property      string
}

func (g *GenericPropertyGet) Yail() string {
	yail := "(get-property-and-check "
	yail += g.Component.Yail()
	yail += " '"
	yail += g.ComponentType
	yail += " '"
	yail += g.Property
	yail += ")"
	return yail
}

func (g *GenericPropertyGet) String() string {
	return sugar.Format("get(%, %, %)", g.ComponentType, g.Component.String(), g.Property)
}

func (g *GenericPropertyGet) Blockly() ast.Block {
	return ast.Block{
		Type: "component_set_get",
		Mutation: &ast.Mutation{
			SetOrGet:      "get",
			PropertyName:  g.Property,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Fields: []ast.Field{{Name: "PROP", Value: g.Property}},
		Values: []ast.Value{{Name: "COMPONENT", Block: g.Component.Blockly()}},
	}
}

func (g *GenericPropertyGet) Continuous() bool {
	return false
}

func (g *GenericPropertyGet) Consumable() bool {
	return true
}
