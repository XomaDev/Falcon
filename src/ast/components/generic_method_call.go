package components

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type GenericMethodCall struct {
	Component blky.Expr
	Method    string
	Args      []blky.Expr
}

func (g *GenericMethodCall) String() string {
	pFormat := "%->%(%)"
	if !g.Component.Continuous() {
		pFormat = "(%)->%(%)"
	}
	return sugar.Format(pFormat, g.Component.String(), g.Method, blky.JoinExprs(", ", g.Args))
}

func (g *GenericMethodCall) Blockly() blky.Block {
	return blky.Block{
		// TODO: add component_type to Mutation later
		Mutation: &blky.Mutation{
			MethodName: g.Method,
			IsGeneric:  true,
		},
		Values: blky.ValueArgsByPrefix(g.Component, "COMPONENT", "ARG", g.Args),
	}
}

func (g *GenericMethodCall) Continuous() bool {
	return false
}
