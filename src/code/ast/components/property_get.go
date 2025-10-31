package components

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type PropertyGet struct {
	ComponentName string
	ComponentType string
	Property      string
}

func (p *PropertyGet) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (p *PropertyGet) String() string {
	return sugar.Format("%.%", p.ComponentName, p.Property)
}

func (p *PropertyGet) Blockly() blockly.Block {
	return blockly.Block{
		Type: "component_set_get",
		Mutation: &blockly.Mutation{
			SetOrGet:      "get",
			PropertyName:  p.Property,
			IsGeneric:     false,
			InstanceName:  p.ComponentName,
			ComponentType: p.ComponentType,
		},
		Fields: []blockly.Field{
			{Name: "COMPONENT_SELECTOR", Value: p.ComponentName},
			{Name: "PROP", Value: p.Property},
		},
	}
}

func (p *PropertyGet) Continuous() bool {
	return false
}

func (p *PropertyGet) Consumable() bool {
	return true
}
