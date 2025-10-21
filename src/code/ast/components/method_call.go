package components

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type MethodCall struct {
	ComponentName string
	ComponentType string
	Method        string
	Args          []blockly.Expr
}

func (m *MethodCall) String() string {
	return sugar.Format("%.%(%)", m.ComponentName, m.Method, blockly.JoinExprs(", ", m.Args))
}

func (m *MethodCall) Blockly() blockly.Block {
	return blockly.Block{
		Type: "component_method",
		Mutation: &blockly.Mutation{
			MethodName:    m.Method,
			IsGeneric:     false,
			InstanceName:  m.ComponentName,
			ComponentType: m.ComponentType,
		},
		Fields: []blockly.Field{{Name: "COMPONENT_SELECTOR", Value: m.ComponentName}},
		Values: blockly.ValuesByPrefix("ARG", m.Args),
	}
}

func (m *MethodCall) Continuous() bool {
	return false
}

func (m *MethodCall) Consumable() bool {
	return false // may be consumable too
}
