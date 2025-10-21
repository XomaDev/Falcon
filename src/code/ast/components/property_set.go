package components

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type PropertySet struct {
	ComponentName string
	ComponentType string
	Property      string
	Value         blockly2.Expr
}

func (p *PropertySet) String() string {
	return sugar.Format("%.% = %", p.ComponentName, p.Property, p.Value.String())
}

func (p *PropertySet) Blockly() blockly2.Block {
	return blockly2.Block{
		Type: "component_set_get",
		Mutation: &blockly2.Mutation{
			SetOrGet:      "set",
			PropertyName:  p.Property,
			IsGeneric:     false,
			InstanceName:  p.ComponentName,
			ComponentType: p.ComponentType,
		},
		Fields: blockly2.FieldsFromMap(map[string]string{
			"COMPONENT_SELECTOR": p.ComponentName,
			"PROP":               p.Property,
		}),
		Values: []blockly2.Value{{Name: "VALUE", Block: p.Value.Blockly()}},
	}
}

func (p *PropertySet) Continuous() bool {
	return false
}

func (p *PropertySet) Consumable() bool {
	return false
}
