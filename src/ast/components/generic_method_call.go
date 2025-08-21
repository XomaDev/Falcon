package components

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type GenericMethodCall struct {
	Component     blky.Expr
	ComponentType string
	Method        string
	Args          []blky.Expr
}

func (g *GenericMethodCall) String() string {
	return sugar.Format("call(%, %, %, %)", g.ComponentType, g.Component.String(), g.Method, blky.JoinExprs(", ", g.Args))
}

func (g *GenericMethodCall) Blockly() blky.Block {
	return blky.Block{
		Type: "component_method",
		Mutation: &blky.Mutation{
			MethodName:    g.Method,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Values: blky.ValueArgsByPrefix(g.Component, "COMPONENT", "ARG", g.Args),
	}
}

func (g *GenericMethodCall) Continuous() bool {
	return false
}
