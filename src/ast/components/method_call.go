package components

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type MethodCall struct {
	Component string
	Method    string
	Args      []blky.Expr
}

func (m *MethodCall) String() string {
	return sugar.Format("%.%(%)", m.Component, m.Method, blky.JoinExprs(", ", m.Args))
}

func (m *MethodCall) Blockly() blky.Block {
	return blky.Block{
		// TODO: add component_type to Mutation
		Mutation: &blky.Mutation{
			MethodName:   m.Method,
			IsGeneric:    false,
			InstanceName: m.Component,
		},
		Fields: []blky.Field{{Name: "COMPONENT_SELECTOR", Value: m.Component}},
		Values: blky.ValuesByPrefix("ARG", m.Args),
	}
}

func (m *MethodCall) Continuous() bool {
	return false
}
