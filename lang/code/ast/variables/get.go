package variables

import (
	"Falcon/code/ast"
	"Falcon/code/lex"
)

type Get struct {
	Where          *lex.Token
	Global         bool
	Name           string
	ValueSignature []ast.Signature
}

func (g *Get) String() string {
	if g.Global {
		return "this." + g.Name
	}
	return g.Name
}

func (g *Get) Blockly(flags ...bool) ast.Block {
	var name string
	if g.Global {
		name = "global " + g.Name
	} else {
		name = g.Name
	}
	return ast.Block{
		Type:   "lexical_variable_get",
		Fields: []ast.Field{{Name: "VAR", Value: name}},
	}
}

func (g *Get) Continuous() bool {
	return true
}

func (g *Get) Consumable(flags ...bool) bool {
	return true
}

func (g *Get) Signature() []ast.Signature {
	// TODO: later define a variable lookup table
	return []ast.Signature{ast.SignAny}
}
