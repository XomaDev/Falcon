package components

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type GenericMethodCall struct {
	Component     blockly.Expr
	ComponentType string
	Method        string
	Args          []blockly.Expr
}

func (g *GenericMethodCall) String() string {
	return sugar.Format("call(%, %, %, %)", g.ComponentType, g.Component.String(), g.Method, blockly.JoinExprs(", ", g.Args))
}

func (g *GenericMethodCall) Blockly() blockly.Block {
	return blockly.Block{
		Type: "component_method",
		Mutation: &blockly.Mutation{
			MethodName:    g.Method,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Values: blockly.ValueArgsByPrefix(g.Component, "COMPONENT", "ARG", g.Args),
	}
}

func (g *GenericMethodCall) Continuous() bool {
	return false
}

func (g *GenericMethodCall) Consumable() bool {
	return false // play safe, may be consumable too
}
