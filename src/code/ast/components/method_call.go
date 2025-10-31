package components

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type MethodCall struct {
	ComponentName string
	ComponentType string
	Method        string
	Args          []ast.Expr
}

func (m *MethodCall) Yail() string {
	//TODO implement me
	panic("implement me")
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
