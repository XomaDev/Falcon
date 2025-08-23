package components

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type MethodCall struct {
	ComponentName string
	ComponentType string
	Method        string
	Args          []blky.Expr
}

func (m *MethodCall) String() string {
	return sugar.Format("%.%(%)", m.ComponentName, m.Method, blky.JoinExprs(", ", m.Args))
}

func (m *MethodCall) Blockly() blky.Block {
	return blky.Block{
		Type: "component_method",
		Mutation: &blky.Mutation{
			MethodName:    m.Method,
			IsGeneric:     false,
			InstanceName:  m.ComponentName,
			ComponentType: m.ComponentType,
		},
		Fields: []blky.Field{{Name: "COMPONENT_SELECTOR", Value: m.ComponentName}},
		Values: blky.ValuesByPrefix("ARG", m.Args),
	}
}

func (m *MethodCall) Continuous() bool {
	return false
}

func (m *MethodCall) Consumable() bool {
	return false // may be consumable too
}
