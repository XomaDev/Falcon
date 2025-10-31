package variables

import (
	blky "Falcon/code/ast"
)

type Global struct {
	Name  string
	Value blky.Expr
}

func (g *Global) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (g *Global) String() string {
	return "global " + g.Name + " = " + g.Value.String()
}

func (g *Global) Blockly() blky.Block {
	return blky.Block{
		Type:   "global_declaration",
		Fields: []blky.Field{{Name: "NAME", Value: g.Name}},
		Values: []blky.Value{{Name: "VALUE", Block: g.Value.Blockly()}},
	}
}

func (g *Global) Continuous() bool {
	return false
}

func (g *Global) Consumable() bool {
	return false
}
