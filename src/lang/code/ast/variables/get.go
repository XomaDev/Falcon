package variables

import (
	blky "Falcon/lang/code/ast"
	"Falcon/lang/code/lex"
)

type Get struct {
	Where  *lex.Token
	Global bool
	Name   string
}

func (g *Get) Yail() string {
	var yail string
	if g.Global {
		yail = "(get-var g$"
	} else {
		yail += "(lexical-value $"
	}
	yail += g.Name
	yail += ")"
	return yail
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
