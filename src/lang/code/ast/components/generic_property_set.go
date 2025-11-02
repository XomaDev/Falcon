package components

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
)

type GenericPropertySet struct {
	Component     ast2.Expr
	ComponentType string
	Property      string
	Value         ast2.Expr
}

func (g *GenericPropertySet) Yail() string {
	yail := "(set-and-coerce-property-and-check! "
	yail += g.Component.Yail()
	yail += " '"
	yail += g.ComponentType
	yail += " '"
	yail += g.Property
	yail += " "
	yail += ast2.PadDirect(g.Value.Yail())
	yail += " '"
	yail += "any)"
	return yail
}

func (g *GenericPropertySet) String() string {
	return sugar.Format("set(%, %, %, %)", g.ComponentType, g.Component.String(), g.Property, g.Value.String())

}

func (g *GenericPropertySet) Blockly() ast2.Block {
	return ast2.Block{
		Type: "component_set_get",
		Mutation: &ast2.Mutation{
			SetOrGet:      "set",
			PropertyName:  g.Property,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Fields: []ast2.Field{{Name: "PROP", Value: g.Property}},
		Values: ast2.MakeValues([]ast2.Expr{g.Component, g.Value}, "COMPONENT", "VALUE"),
	}
}

func (g *GenericPropertySet) Continuous() bool {
	return false
}

func (g *GenericPropertySet) Consumable() bool {
	return false
}
