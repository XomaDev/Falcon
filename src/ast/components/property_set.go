package components

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type PropertySet struct {
	ComponentName string
	ComponentType string
	Property      string
	Value         blockly.Expr
}

func (p *PropertySet) String() string {
	return sugar.Format("%.% = %", p.ComponentName, p.Property, p.Value.String())
}

func (p *PropertySet) Blockly() blockly.Block {
	return blockly.Block{
		Type: "component_set_get",
		Mutation: &blockly.Mutation{
			SetOrGet:      "set",
			PropertyName:  p.Property,
			IsGeneric:     false,
			InstanceName:  p.ComponentName,
			ComponentType: p.ComponentType,
		},
		Fields: blockly.FieldsFromMap(map[string]string{
			"COMPONENT_SELECTOR": p.ComponentName,
			"PROP":               p.Property,
		}),
		Values: []blockly.Value{{Name: "VALUE", Block: p.Value.Blockly()}},
	}
}

func (p *PropertySet) Continuous() bool {
	return false
}

func (p *PropertySet) Consumable() bool {
	return false
}
