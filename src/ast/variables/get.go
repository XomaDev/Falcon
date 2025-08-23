package variables

import (
	blky "Falcon/ast/blockly"
	"Falcon/lex"
)

type Get struct {
	Where  *lex.Token
	Global bool
	Name   string
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
