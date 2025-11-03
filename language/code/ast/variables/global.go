package variables

import (
	ast2 "Falcon/code/ast"
)

type Global struct {
	Name  string
	Value ast2.Expr
}

func (g *Global) Yail() string {
	yail := "(def g$"
	yail += g.Name
	yail += " "
	yail += g.Value.Yail()
	yail += ")"
	return yail
}

func (g *Global) String() string {
	return "global " + g.Name + " = " + g.Value.String()
}

func (g *Global) Blockly() ast2.Block {
	return ast2.Block{
		Type:   "global_declaration",
		Fields: []ast2.Field{{Name: "NAME", Value: g.Name}},
		Values: []ast2.Value{{Name: "VALUE", Block: g.Value.Blockly()}},
	}
}

func (g *Global) Continuous() bool {
	return false
}

func (g *Global) Consumable() bool {
	return false
}
