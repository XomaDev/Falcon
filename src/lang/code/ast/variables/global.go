package variables

import (
	"Falcon/lang/code/ast"
)

type Global struct {
	Name  string
	Value ast.Expr
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

func (g *Global) Blockly() ast.Block {
	return ast.Block{
		Type:   "global_declaration",
		Fields: []ast.Field{{Name: "NAME", Value: g.Name}},
		Values: []ast.Value{{Name: "VALUE", Block: g.Value.Blockly()}},
	}
}

func (g *Global) Continuous() bool {
	return false
}

func (g *Global) Consumable() bool {
	return false
}
