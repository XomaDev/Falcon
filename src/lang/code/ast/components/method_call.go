package components

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type MethodCall struct {
	ComponentName string
	ComponentType string
	Method        string
	Args          []ast2.Expr
}

func (m *MethodCall) Yail() string {
	yail := "(call-component-method "
	//TODO ??? there seems to be some special casing for Block
	yail += "'"
	yail += m.ComponentName
	yail += " '"
	yail += m.Method
	yail += " (*list-for-runtime*"
	yail += ast2.JoinYailExprs(" ", m.Args)
	yail += ") '("
	yail += strings.Repeat("any ", len(m.Args))
	yail += "))"
	yail += "))"
	return yail
}

func (m *MethodCall) String() string {
	return sugar.Format("%.%(%)", m.ComponentName, m.Method, ast2.JoinExprs(", ", m.Args))
}

func (m *MethodCall) Blockly() ast2.Block {
	return ast2.Block{
		Type: "component_method",
		Mutation: &ast2.Mutation{
			MethodName:    m.Method,
			IsGeneric:     false,
			InstanceName:  m.ComponentName,
			ComponentType: m.ComponentType,
		},
		Fields: []ast2.Field{{Name: "COMPONENT_SELECTOR", Value: m.ComponentName}},
		Values: ast2.ValuesByPrefix("ARG", m.Args),
	}
}

func (m *MethodCall) Continuous() bool {
	return false
}

func (m *MethodCall) Consumable() bool {
	return false // may be consumable too
}
