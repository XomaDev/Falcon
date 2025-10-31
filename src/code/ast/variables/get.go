package variables

import (
	blky "Falcon/code/ast/blockly"
	"Falcon/code/lex"
)

type Get struct {
	Where  *lex.Token
	Global bool
	Name   string
}

func (g *Get) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (g *Get) String() string {
	if g.Global {
		return "this." + g.Name
	}
	return g.Name
}

func (g *Get) Blockly() blky.Block {
	var name string
	if g.Global {
		name = "global " + g.Name
	} else {
		name = g.Name
	}
	return blky.Block{
		Type:   "lexical_variable_get",
		Fields: []blky.Field{{Name: "VAR", Value: name}},
	}
}

func (g *Get) Continuous() bool {
	return true
}

func (g *Get) Consumable() bool {
	return true
}
