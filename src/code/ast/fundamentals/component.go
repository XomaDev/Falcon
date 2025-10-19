package fundamentals

import (
	"Falcon/code/ast/blockly"
)

type Component struct {
	Name string
	Type string
}

func (c *Component) String() string {
	return c.Name
}

func (c *Component) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "component_component_block",
		Mutation: &blockly.Mutation{InstanceName: c.Name, ComponentType: c.Type},
		Fields:   []blockly.Field{{Name: "COMPONENT_SELECTOR", Value: c.Name}},
	}
}

func (c *Component) Continuous() bool {
	return true
}

func (c *Component) Consumable() bool {
	return true
}
