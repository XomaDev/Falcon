package components

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
	"strings"
)

type MethodCall struct {
	ComponentName string
	ComponentType string
	Method        string
	Args          []ast.Expr
}

func (m *MethodCall) Yail() string {
	yail := "(call-component-method "
	//TODO ??? there seems to be some special casing for Block
	yail += "'"
	yail += m.ComponentName
	yail += " '"
	yail += m.Method
	yail += " (*list-for-runtime*"
	yail += ast.JoinYailExprs(" ", m.Args)
	yail += ") '("
	yail += strings.Repeat("any ", len(m.Args))
	yail += "))"
	yail += "))"
	return yail
}

func (m *MethodCall) String() string {
	return sugar.Format("%.%(%)", m.ComponentName, m.Method, ast.JoinExprs(", ", m.Args))
}

func (m *MethodCall) Blockly() ast.Block {
	return ast.Block{
		Type: "component_method",
		Mutation: &ast.Mutation{
			MethodName:    m.Method,
			IsGeneric:     false,
			InstanceName:  m.ComponentName,
			ComponentType: m.ComponentType,
		},
		Fields: []ast.Field{{Name: "COMPONENT_SELECTOR", Value: m.ComponentName}},
		Values: ast.ValuesByPrefix("ARG", m.Args),
	}
}

func (m *MethodCall) Continuous() bool {
	return false
}

func (m *MethodCall) Consumable() bool {
	return false // may be consumable too
}
