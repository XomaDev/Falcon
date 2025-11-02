package components

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type GenericMethodCall struct {
	Component     ast2.Expr
	ComponentType string
	Method        string
	Args          []ast2.Expr
}

func (g *GenericMethodCall) Yail() string {
	yail := "(call-component-type-method "
	yail += g.Component.Yail()
	yail += " '"
	yail += g.ComponentType
	yail += " '"
	yail += g.Method
	yail += " (*list-for-runtime* "
	yail += ast2.JoinYailExprs(" ", g.Args)
	yail += ") '("
	yail += strings.Repeat("any ", len(g.Args))
	yail += "))"
	return yail
}

func (g *GenericMethodCall) String() string {
	return sugar.Format("call(%, %, %, %)", g.ComponentType, g.Component.String(), g.Method, ast2.JoinExprs(", ", g.Args))
}

func (g *GenericMethodCall) Blockly() ast2.Block {
	return ast2.Block{
		Type: "component_method",
		Mutation: &ast2.Mutation{
			MethodName:    g.Method,
			IsGeneric:     true,
			ComponentType: g.ComponentType,
		},
		Values: ast2.ValueArgsByPrefix(g.Component, "COMPONENT", "ARG", g.Args),
	}
}

func (g *GenericMethodCall) Continuous() bool {
	return false
}

func (g *GenericMethodCall) Consumable() bool {
	return false // play safe, may be consumable too
}
